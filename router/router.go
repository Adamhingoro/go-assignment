package router

import (
	"assigment/core"
	"assigment/handler"
	"assigment/middleware"
	"log"

	"github.com/gorilla/mux"
)

type Router struct {
	core           *core.CoreConfig
	mux            *mux.Router
	UserHandler    *handler.UserHandler
	CompanyHandler *handler.CompanyHandler
	AuthHandler    *handler.AuthHandler
}

func NewRouter(userHandler *handler.UserHandler, companyHandler *handler.CompanyHandler, core *core.CoreConfig, authHandler *handler.AuthHandler) *Router {
	log.Println("Initialiing the router")
	r := mux.NewRouter()
	r.Use(middleware.JSONMiddleware)
	return &Router{mux: r, UserHandler: userHandler, core: core, CompanyHandler: companyHandler, AuthHandler: authHandler}
}

func (r *Router) RegisterRoutes() {
	r.RegisterUserRoutes()
	r.RegisterAuthRoutes()
	r.RegisterCompanyRoutes()
	log.Println("All routes registered successfully")
}

func (r *Router) RegisterUserRoutes() {
	r.mux.HandleFunc("/users", middleware.ValidateJWTMiddleware(r.UserHandler.GetUsers, r.core)).Methods("GET")
	r.mux.HandleFunc("/users", middleware.ValidateJWTMiddleware(r.UserHandler.CreateUser, r.core)).Methods("POST")
	r.mux.HandleFunc("/users/{id}", middleware.ValidateJWTMiddleware(r.UserHandler.UpdateUser, r.core)).Methods("PUT")
	r.mux.HandleFunc("/users/{id}", middleware.ValidateJWTMiddleware(r.UserHandler.DeleteUser, r.core)).Methods("DELETE")
}

func (r *Router) RegisterCompanyRoutes() {
	r.mux.HandleFunc("/companies", middleware.ValidateJWTMiddleware(r.CompanyHandler.GetCompanies, r.core)).Methods("GET")
	r.mux.HandleFunc("/companies", middleware.ValidateJWTMiddleware(r.CompanyHandler.CreateCompany, r.core)).Methods("POST")
	r.mux.HandleFunc("/companies/{id}", middleware.ValidateJWTMiddleware(r.CompanyHandler.UpdateCompany, r.core)).Methods("PUT")
	r.mux.HandleFunc("/companies/{id}", middleware.ValidateJWTMiddleware(r.CompanyHandler.DeleteCompany, r.core)).Methods("DELETE")
}

func (r *Router) RegisterAuthRoutes() {
	r.mux.HandleFunc("/login", r.AuthHandler.Login).Methods("POST")
}

func (r *Router) GetMux() *mux.Router {
	return r.mux
}
