package interactionservice

import "gorm.io/gorm"

type InteractionRepository struct {
	DB           *gorm.DB
	Interactions *[]UserInteraction
}

func _NewRepository(db *gorm.DB, interatcions *[]UserInteraction) *InteractionRepository {
	return &InteractionRepository{DB: db, Interactions: interatcions}
}

func (repo *InteractionRepository) CreateUserInteraction() error {

	for _, interaction := range *repo.Interactions {
		err := repo.DB.Create(&interaction).Error
		if err != nil {
			return err
		}
	}
	return nil
}
