package interactionservice

import "gorm.io/gorm"

type InteractionService struct {
	DB           *gorm.DB
	Interactions *[]UserInteraction
}

func _NewService(db *gorm.DB, interatcions *[]UserInteraction) *InteractionService {
	return &InteractionService{DB: db, Interactions: interatcions}
}

func (service *InteractionService) CreateUserInteraction() error {
	repo := _NewRepository(service.DB, service.Interactions)
	return repo.CreateUserInteraction()
}
