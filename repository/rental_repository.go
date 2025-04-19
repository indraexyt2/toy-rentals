package repository

import (
	"context"
	"final-project/entity"
	"gorm.io/gorm"
)

type IRentalRepository interface {
	IBaseRepository[entity.Rental]
	UpdateToyStock(ctx context.Context, toyID string, quantity int) error
	ReturnRental(ctx context.Context, rental *entity.Rental) error
	UpdateRentalItem(ctx context.Context, rentalItem *entity.RentalItem) error
}

type RentalRepository struct {
	BaseRepository[entity.Rental]
}

func NewRentalRepository(db *gorm.DB) IRentalRepository {
	return &RentalRepository{
		BaseRepository: BaseRepository[entity.Rental]{DB: db},
	}
}

func (r *RentalRepository) FindById(ctx context.Context, id string) (entity.Rental, error) {
	var model entity.Rental
	if err := r.DB.WithContext(ctx).Where("id = ?", id).
		Preload("RentalItems").
		First(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (r *RentalRepository) Insert(ctx context.Context, model *entity.Rental) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("RentalItems").Create(model).Error; err != nil {
			return err
		}

		for i := range model.RentalItems {
			model.RentalItems[i].RentalID = model.ID
			if err := tx.Create(&model.RentalItems[i]).Error; err != nil {
				return err
			}

			if err := tx.Model(&entity.Toy{}).Where("id = ?", model.RentalItems[i].ToyID).
				UpdateColumn("stock", gorm.Expr("stock - ?", model.RentalItems[i].Quantity)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RentalRepository) UpdateToyStock(ctx context.Context, toyID string, quantity int) error {
	return r.DB.WithContext(ctx).Model(&entity.Toy{}).Where("id = ?", toyID).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error
}

func (r *RentalRepository) ReturnRental(ctx context.Context, rental *entity.Rental) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(rental).
			Select("status", "actual_return_date", "late_fee", "damage_fee", "total_amount").
			Updates(rental).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *RentalRepository) UpdateRentalItem(ctx context.Context, rentalItem *entity.RentalItem) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(rentalItem).
			Select("condition_after", "damage_description", "damage_fee", "status").
			Updates(rentalItem).Error; err != nil {
			return err
		}

		if rentalItem.Status == "returned" {
			if err := tx.Model(&entity.Toy{}).
				Where("id = ?", rentalItem.ToyID).
				UpdateColumn("stock", gorm.Expr("stock + ?", rentalItem.Quantity)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
