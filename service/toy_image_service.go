package service

import (
	"final-project/entity"
	"final-project/repository"
)

type IToyImageService interface {
	IBaseService[entity.ToyImage]
}

type ToyImageService struct {
	BaseService[entity.ToyImage]
}

func NewToyImageService(repo repository.IToyImageRepository) IToyImageService {
	return &ToyImageService{
		BaseService: BaseService[entity.ToyImage]{repository: repo},
	}
}
