package repository

import (
	"cli-service/internal/model"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type CLIRepository struct {
	db *gorm.DB
}

func NewCLIRepository(db *gorm.DB) *CLIRepository {
	return &CLIRepository{db: db}
}

func (repo *CLIRepository) Create(events model.Events) (model.Events, error) {
	result := repo.db.Create(events.Events)
	if result.Error != nil {
		log.Println("error with create gorm model:", result.Error)
		return model.Events{}, result.Error
	}

	return events, nil
}

func (repo *CLIRepository) Get() (model.Events, error) {
	var events model.Events

	result := repo.db.Find(&events.Events)
	if result.Error != nil {
		log.Println("fail with find events gorm:", result.Error)
		return model.Events{}, result.Error
	}

	return events, nil
}

func (repo *CLIRepository) Drop() error {
	var events model.Events

	result := repo.db.Where("language = ?", "golang").Delete(&events.Events)
	if result.Error != nil {
		log.Println("fail delete events gorm:", result.Error)
		return result.Error
	}

	fmt.Println("deleting successful")
	return nil
}
