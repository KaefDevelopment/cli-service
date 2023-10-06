package repository

import (
	"cli-service/internal/model"
	"cli-service/internal/service/dto"
	"gorm.io/gorm"
	"log"
)

type CLIRepository struct {
	db *gorm.DB
}

func NewCLIRepository(db *gorm.DB) *CLIRepository {
	return &CLIRepository{db: db}
}

func (repo *CLIRepository) Create(events model.Events) error {
	result := repo.db.Create(&events.Events)
	if result.Error != nil {
		log.Println("error with create gorm model:", result.Error)
		return result.Error
	}

	return nil
}

func (repo *CLIRepository) Update() error {
	result := repo.db.Table("events").Where("send = ?", 0).Updates(map[string]interface{}{"send": 1})
	if result.Error != nil {
		log.Println("fail with update gorm column:", result.Error)
	}

	return nil
}

func (repo *CLIRepository) Get(authKey []string) (model.EventsByAuthKey, error) {
	var res model.EventsByAuthKey

	for _, key := range authKey {
		var events model.Events

		result := repo.db.Where("send = ?", 1).Where("authKey = ?", key).Find(&events.Events)
		if result.Error != nil {
			log.Println("fail with find events gorm:", result.Error)
			return model.EventsByAuthKey{}, result.Error
		}

		res.Events = append(res.Events, events)
	}

	return res, nil
}

func (repo *CLIRepository) GetAuthKeys() ([]string, error) {
	result := repo.db.Table("events")

	result = result.Raw("select distinct authKey from events")

	rows, err := result.Rows()
	if err != nil {
		log.Println("fail with rows gorm:", err)
		return nil, err
	}

	defer rows.Close()

	var keys []string

	for rows.Next() {
		var authKey model.Event

		if err := rows.Scan(&authKey.AuthKey); err != nil {
			log.Println("fail rows scan gorm:", err)
			return nil, err
		}

		keys = append(keys, authKey.AuthKey)
	}

	return keys, nil
}

func (repo *CLIRepository) Drop() error {
	var events dto.Events

	result := repo.db.Where("send = ?", "1").Delete(&events.Events)
	if result.Error != nil {
		log.Println("fail delete events gorm:", result.Error)
		return result.Error
	}

	log.Println("deleting successful")

	return nil
}
