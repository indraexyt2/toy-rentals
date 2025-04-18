package service

import (
	"final-project/entity"
	"final-project/repository"
)

type IToyService interface {
	IBaseService[entity.Toy]
}

type ToyService struct {
	BaseService[entity.Toy]
}

func NewToyService(repo repository.IToyRepository) IToyService {
	return &ToyService{
		BaseService: BaseService[entity.Toy]{repository: repo},
	}
}
