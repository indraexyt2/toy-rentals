package entity

type ToyImage struct {
	BaseEntity
	ImageURL  string `gorm:"size:255;not null" json:"image_url"`
	IsPrimary bool   `gorm:"default:false" json:"is_primary"`

	Toy []Toy `gorm:"many2many:image_toys" json:"-"`
}

func (*ToyImage) TableName() string {
	return "toy_images"
}
