package cronjob

import (
	"assigment/database"
	"assigment/model"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func CheckAndUpdateUsers(db *database.Database) {
	log.Println("Starting JOB")
	// Calculate the time threshold (5 hours ago)
	fiveHoursAgo := time.Now().Add(-5 * time.Hour)

	// Update users where last_seen is older than 5 hours
	result := db.DB.Model(&model.User{}).
		Where("last_seen < ?", fiveHoursAgo).
		Update("is_available", false)

	if result.Error != nil {
		log.Printf("Error updating users: %v", result.Error)
	} else {
		log.Printf("Updated %d users", result.RowsAffected)
	}
}

func NewCheckAndUpdateUsersCronJob(db *database.Database) *cron.Cron {
	log.Println("Scheduling the Cronjob for updating users")
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		CheckAndUpdateUsers(db)
	})
	if err != nil {
		log.Fatalln("Error initilizing the cronjob", err)
	}

	c.Start()

	return c
}
