package controller

import (
	"errors"
	"final-project/service"
	"final-project/utils/helpers"
	"final-project/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type IToyController interface {
	FindAll(c *gin.Context)
	FinById(c *gin.Context)
	Insert(c *gin.Context)
	UpdateById(c *gin.Context)
	DeleteById(c *gin.Context)
}

type ToyController struct {
	toySvc service.IToyService
}

func NewToyController(toySvc service.IToyService) IToyController {
	return &ToyController{
		toySvc: toySvc,
	}
}

// FindAll godoc
// @Description Get all toys
// @Tags Toy
// @Produce json
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} entity.Toy
// @Router /toy [get]
func (t ToyController) FindAll(c *gin.Context) {
	var logger = helpers.Logger

	var page = c.DefaultQuery("page", "1")
	var pageInt = helpers.ParseToInt(page)

	var limit = c.DefaultQuery("limit", "10")
	var limitInt = helpers.ParseToInt(limit)

	var offset = (pageInt - 1) * limitInt

	data, totalData, err := t.toySvc.FindAll(c.Request.Context(), limitInt, offset)
	if err != nil {
		logger.Error("Failed to find all toys: ", err)
		response.ResponseError(c, http.StatusInternalServerError, "Failed to find all toys")
		return
	}

	metaData := response.Page{
		Limit:     limitInt,
		Total:     int(totalData),
		Page:      pageInt,
		TotalPage: int(totalData) / limitInt,
	}

	response.ResponseSuccess(c, http.StatusOK, data, metaData, "Success to find all toys")
}

// FindById godoc
// @Description Get toy by id
// @Tags Toy
// @Produce json
// @Param id path string true "Toy ID"
// @Success 200 {object} entity.Toy
// @Router /toy/{id} [get]
func (t ToyController) FinById(c *gin.Context) {
	var logger = helpers.Logger

	var id = c.Param("id")
	if id == "" {
		logger.Error("Id is required")
		response.ResponseError(c, http.StatusBadRequest, "Id is required")
		return
	}

	data, err := t.toySvc.FindById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(fmt.Errorf("toy with id %s not found", id))
			response.ResponseError(c, http.StatusNotFound, "Toy not found")
			return
		}

		logger.Error(fmt.Errorf("failed to find toy by id %s: %v", id, err))
		response.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.ResponseSuccess(c, http.StatusOK, data, nil, "Success to find toy")
}

// Insert godoc
// @Description Upload a new toy with images
// @Tags Toy
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Toy name"
// @Param description formData string true "Toy description"
// @Param age_recommendation formData string true "Age recommendation"
// @Param condition formData string true "Condition (new, excellent, good, fair, poor)"
// @Param rental_price formData number true "Rental price"
// @Param late_fee_per_day formData number true "Late fee per day"
// @Param replacement_price formData number true "Replacement price"
// @Param is_available formData boolean false "Is available"
// @Param stock formData int true "Stock"
// @Param categories formData []string true "Category IDs" collectionFormat(multi)
// @Param is_primary_index formData int false "Index of primary image"
// @Param images formData file true "Upload multiple images"
// @Success 200 {object} entity.Toy
// @Router /toy [post]
func (t ToyController) Insert(c *gin.Context) {

}

// UpdateById godoc
// @Summary Update toy by ID
// @Description Update the toy details by providing its ID and form data (including images)
// @Tags Toy
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Toy ID"
// @Param name formData string true "Toy Name"
// @Param description formData string true "Toy Description"
// @Param age_recommendation formData string true "Age Recommendation"
// @Param rental_price formData float64 true "Rental Price"
// @Param late_fee_per_day formData float64 true "Late Fee Per Day"
// @Param replacement_price formData float64 true "Replacement Price"
// @Param stock formData int true "Stock"
// @Param is_available formData bool true "Is Available"
// @Param categories formData []string true "Category IDs" collectionFormat(multi)
// @Param images formData file true "Upload multiple images"
// @Param deleted_images formData string false "Deleted Image IDs (comma separated)"
// @Success 200 {object} entity.Toy "Updated Toy"
// @Router /toy/{id} [put]
func (t ToyController) UpdateById(c *gin.Context) {
}

func (t ToyController) DeleteById(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
