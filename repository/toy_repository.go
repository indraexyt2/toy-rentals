package repository

import (
	"context"
	"final-project/entity"
	"gorm.io/gorm"
)

type IToyRepository interface {
	IBaseRepository[entity.Toy]
}

type ToyRepository struct {
	BaseRepository[entity.Toy]
}

func NewToyRepository(db *gorm.DB) IToyRepository {
	return &ToyRepository{
		BaseRepository: BaseRepository[entity.Toy]{DB: db},
	}
}

func (r *ToyRepository) Update(ctx context.Context, toy *entity.Toy) error {
	return r.DB.WithContext(ctx).Model(&entity.Toy{}).Where("id = ?", toy.ID).Updates(map[string]interface{}{
		"name":               toy.Name,
		"description":        toy.Description,
		"age_recommendation": toy.AgeRecommendation,
		"condition":          toy.Condition,
		"rental_price":       toy.RentalPrice,
		"late_fee_per_day":   toy.LateFeePerDay,
		"replacement_price":  toy.ReplacementPrice,
		"is_available":       toy.IsAvailable,
		"stock":              toy.Stock,
	}).Error
}

func (r *ToyRepository) FindAll(ctx context.Context, limit int, offset int) ([]entity.Toy, int64, error) {
	var entities []entity.Toy
	if err := r.DB.WithContext(ctx).
		Preload("Categories").
		Preload("Images").
		Limit(limit).Offset(offset).
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	var totalData int64
	if err := r.DB.WithContext(ctx).Model(new(entity.Toy)).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}
	return entities, totalData, nil
}
