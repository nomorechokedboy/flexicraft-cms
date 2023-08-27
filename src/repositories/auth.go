package repositories

import (
	apperr "api/src/app_error"
	"api/src/entities"
	"net/http"

	"gorm.io/gorm"
)

type (
	AuthRepo interface {
		Insert(entities.CreateAuth) (*entities.AuthEntity, error)
		FindOne(entities.AuthEntity) (*entities.AuthEntity, error)
	}

	GormAuthRepo struct {
		db *gorm.DB
	}
)

func New(db *gorm.DB) AuthRepo {
	return &GormAuthRepo{db}
}

var _ AuthRepo = (*GormAuthRepo)(nil)

func (gap *GormAuthRepo) Insert(payload entities.CreateAuth) (*entities.AuthEntity, error) {
	authRes := &entities.AuthEntity{
		Identifier: payload.Identifier,
		Password:   payload.Password,
	}
	if err := gap.db.Create(&authRes).Error; err != nil {
		return nil, apperr.New(
			"101003",
			http.StatusInternalServerError,
			"Unknow error",
			"Internal server error",
			err,
		)
	}

	return authRes, nil
}

func (gap *GormAuthRepo) FindOne(payload entities.AuthEntity) (*entities.AuthEntity, error) {
	if err := gap.db.Take(&payload).Error; err != nil {
		return nil, apperr.New("101004", http.StatusNotFound, "Entity not found", "Not found", err)
	}

	return &payload, nil
}
