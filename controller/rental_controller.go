package controller

import (
	"errors"
	"final-project/entity"
	"final-project/service"
	"final-project/utils/helpers"
	"final-project/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type IRentalController interface {
	FindAll(c *gin.Context)
	FinById(c *gin.Context)
	Insert(c *gin.Context)
	UpdateById(c *gin.Context)
	DeleteById(c *gin.Context)
	ReturnRental(c *gin.Context)
}

type RentalController struct {
	RentalSvc service.IRentalService
}

func NewRentalController(rentalSvc service.IRentalService) IRentalController {
	return &RentalController{
		RentalSvc: rentalSvc,
	}
}

// FindAll godoc
// @Description Get all rentals
// @Tags Rental
// @Produce json
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} entity.Rental
// @Router /rental [get]
func (r *RentalController) FindAll(c *gin.Context) {
	var logger = helpers.Logger

	var page = c.DefaultQuery("page", "1")
	var pageInt = helpers.ParseToInt(page)

	var limit = c.DefaultQuery("limit", "10")
	var limitInt = helpers.ParseToInt(limit)

	var offset = (pageInt - 1) * limitInt

	data, totalData, err := r.RentalSvc.FindAll(c.Request.Context(), limitInt, offset)
	if err != nil {
		logger.Error("Failed to find all rentals: ", err)
		response.ResponseError(c, http.StatusInternalServerError, "Failed to find all rentals")
		return
	}

	metaData := response.Page{
		Limit:     limitInt,
		Total:     int(totalData),
		Page:      pageInt,
		TotalPage: int(totalData) / limitInt,
	}

	response.ResponseSuccess(c, http.StatusOK, data, metaData, "Success get all rentals")
}

// FinById godoc
// @Summary Get rental by id
// @Description Get rental by id
// @Tags Rental
// @Produce json
// @Param id path string true "Rental ID"
// @Success 200 {object} entity.Rental
// @Router /rental/{id} [get]
func (r *RentalController) FinById(c *gin.Context) {
	var logger = helpers.Logger

	var id = c.Param("id")
	if id == "" {
		logger.Error("Id is required")
		response.ResponseError(c, http.StatusBadRequest, "Id is required")
		return
	}

	data, err := r.RentalSvc.FindById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(fmt.Errorf("rental with id %s not found", id))
			response.ResponseError(c, http.StatusNotFound, "Rental not found")
			return
		}

		logger.Error(fmt.Errorf("failed to find rental by id %s: %v", id, err))
		response.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.ResponseSuccess(c, http.StatusOK, data, nil, "Success get rental")
}

// Insert godoc
// @Summary Insert rental
// @Description Insert rental
// @Tags Rental
// @Accept json
// @Produce json
// @Param rental body entity.CreateRentalRequest true "Rental"
// @Success 200 {object} entity.Rental
// @Router /rental [post]
func (r *RentalController) Insert(c *gin.Context) {
	var logger = helpers.Logger

	claims, exists := c.Get("claims")
	if !exists {
		logger.Error("Claims not found in context")
		response.ResponseError(c, http.StatusUnauthorized, "Claims not found in context")
		return
	}

	claimsData, ok := claims.(*helpers.ClaimsToken)
	if !ok {
		logger.Error("Invalid claims type")
		response.ResponseError(c, http.StatusUnauthorized, "Invalid claims type")
		return
	}

	var reqBody entity.CreateRentalRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Error("Failed to bind JSON: ", err)
		response.ResponseError(c, http.StatusBadRequest, "Failed to bind JSON")
		return
	}

	reqBody.UserID = claimsData.UserID

	rental, err := r.RentalSvc.CreateRental(c.Request.Context(), reqBody)
	if err != nil {
		logger.Error("Failed to insert rental: ", err)
		response.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.ResponseSuccess(c, http.StatusOK, rental, nil, "Success insert rental")
}

// UpdateById godoc
// @Description Update rental by id
// @Tags Rental
// @Accept json
// @Produce json
// @Param id path string true "Rental ID"
// @Param rental body entity.Rental true "Rental"
// @Success 200 {object} response.APISuccessResponse
// @Router /rental/{id} [put]
func (r *RentalController) UpdateById(c *gin.Context) {
}

// DeleteById godoc
// @Description Delete rental by id
// @Tags Rental
// @Param id path string true "Rental ID"
// @Success 200 {object} response.APISuccessResponse
// @Router /rental/{id} [delete]
func (r *RentalController) DeleteById(c *gin.Context) {
	var logger = helpers.Logger

	var id = c.Param("id")
	if id == "" {
		logger.Error("Id is required")
		response.ResponseError(c, http.StatusBadRequest, "Id is required")
		return
	}

	err := r.RentalSvc.DeleteById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(fmt.Errorf("rental with id %s not found", id))
			response.ResponseError(c, http.StatusNotFound, "Rental not found")
			return
		}

		logger.Error(fmt.Errorf("failed to delete rental by id %s: %v", id, err))
		response.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.ResponseSuccess(c, http.StatusOK, nil, nil, "Success delete rental")
}

// Return godoc
// @Summary Return rental
// @Description Memproses pengembalian rental dan update kondisi mainan
// @Tags Rental
// @Accept json
// @Produce json
// @Param id path string true "Rental ID"
// @Param request body entity.ReturnRentalRequest true "Data pengembalian rental"
// @Success 200 {object} entity.Rental
// @Router /rental/{id}/return [put]
func (r *RentalController) ReturnRental(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var request entity.ReturnRentalRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rental, err := r.RentalSvc.ReturnRental(ctx.Request.Context(), idStr, request)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "rental tidak ditemukan" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rental)
}
