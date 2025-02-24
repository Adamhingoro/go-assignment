package handler

import (
	"assigment/database"
	"assigment/model"
	"encoding/json"
	"net/http"
)

type CompanyHandler struct {
	DB *database.Database
}

func NewCompanyHandler(db *database.Database) *CompanyHandler {
	return &CompanyHandler{DB: db}
}

func (h *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	var companies []model.Company
	h.DB.DB.Find(&companies)
	json.NewEncoder(w).Encode(companies)
}

func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company model.Company
	json.NewDecoder(r.Body).Decode(&company)

	h.DB.DB.Create(&company)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}

func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var company model.Company
	json.NewDecoder(r.Body).Decode(&company)

	id := r.URL.Query().Get("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Company ID is required")
		return
	}

	result := h.DB.DB.Model(&company).Where("id = ?", id).Updates(company)
	if result.Error != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to update company")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}

func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Company ID is required")
		return
	}

	result := h.DB.DB.Where("id = ?", id).Delete(&model.Company{})
	if result.Error != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to delete company")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
