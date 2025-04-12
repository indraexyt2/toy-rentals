package entity

import (
	"github.com/gofrs/uuid/v5"
)

type ToyImage struct {
	BaseEntity
	ToyID     uuid.UUID `gorm:"type:uuid;not null" json:"toy_id"`
	ImageURL  string    `gorm:"size:255;not null" json:"image_url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`

	Toy Toy `gorm:"foreignKey:ToyID" json:"-"`
}

func (*ToyImage) TableName() string {
	return "toy_images"
}
