package service

import (
	"context"
	"errors"
	"final-project/entity"
	"final-project/repository"
	"fmt"
	"github.com/gofrs/uuid/v5"
)

type IRentalService interface {
	IBaseService[entity.Rental]
	CreateRental(ctx context.Context, req entity.CreateRentalRequest) (*entity.Rental, error)
	ReturnRental(ctx context.Context, id string, req entity.ReturnRentalRequest) (*entity.Rental, error)
}

type RentalService struct {
	BaseService[entity.Rental]
	rentalRepo repository.IRentalRepository
	userRepo   repository.IUserRepository
	toyRepo    repository.IToyRepository
}

func NewRentalService(
	repo repository.IRentalRepository,
	userRepo repository.IUserRepository,
	toyRepo repository.IToyRepository,
) IRentalService {
	return &RentalService{
		BaseService: BaseService[entity.Rental]{repository: repo},
		rentalRepo:  repo,
		userRepo:    userRepo,
		toyRepo:     toyRepo,
	}
}

func (s *RentalService) CreateRental(ctx context.Context, req entity.CreateRentalRequest) (*entity.Rental, error) {
	rental := &entity.Rental{
		UserID:             req.UserID,
		Status:             "pending",
		RentalDate:         req.RentalDate,
		ExpectedReturnDate: req.ExpectedReturnDate,
		TotalRentalPrice:   0,
		PaymentStatus:      "unpaid",
		Notes:              req.Notes,
		RentalItems:        make([]entity.RentalItem, 0, len(req.Items)),
	}

	var totalPrice float64 = 0
	for _, item := range req.Items {
		toy, err := s.toyRepo.FindById(ctx, item.ToyID.String())
		if err != nil {
			return nil, errors.New("mainan tidak ditemukan: " + item.ToyID.String())
		}

		if toy.Stock < item.Quantity {
			return nil, errors.New("stok mainan tidak mencukupi: " + toy.Name)
		}

		pricePerUnit := toy.RentalPrice
		itemTotalPrice := float64(item.Quantity) * pricePerUnit
		totalPrice += itemTotalPrice

		rentalItem := entity.RentalItem{
			ToyID:           item.ToyID,
			Quantity:        item.Quantity,
			PricePerUnit:    pricePerUnit,
			ConditionBefore: item.ConditionBefore,
			ConditionAfter:  item.ConditionBefore,
			Status:          "rented",
		}

		rental.RentalItems = append(rental.RentalItems, rentalItem)
	}

	rental.TotalRentalPrice = totalPrice
	if err := s.repository.Insert(ctx, rental); err != nil {
		return nil, err
	}

	return rental, nil
}

func (s *RentalService) ReturnRental(ctx context.Context, id string, req entity.ReturnRentalRequest) (*entity.Rental, error) {
	// Get rental
	rental, err := s.repository.FindById(ctx, id)
	if err != nil {
		return nil, errors.New("rental tidak ditemukan")
	}

	for _, item := range rental.RentalItems {
		fmt.Println(item)
	}

	// Validasi status rental
	if rental.Status == entity.RentalStatusCompleted || rental.Status == entity.RentalStatusCancelled {
		return nil, errors.New("rental sudah selesai atau dibatalkan")
	}

	// Validasi tanggal
	if req.ActualReturnDate.Before(rental.RentalDate) {
		return nil, errors.New("tanggal pengembalian tidak boleh sebelum tanggal rental")
	}

	// Set tanggal pengembalian aktual
	rental.ActualReturnDate = &req.ActualReturnDate

	// Map untuk melacak ID rental items
	rentalItemMap := make(map[uuid.UUID]*entity.RentalItem)
	for i := range rental.RentalItems {
		rentalItemMap[rental.RentalItems[i].ID] = &rental.RentalItems[i]
	}

	fmt.Println(rentalItemMap)

	// Hitung late fee jika terlambat
	var totalLateFee float64 = 0
	//var isLate bool = false

	if req.ActualReturnDate.After(rental.ExpectedReturnDate) {
		days := int(req.ActualReturnDate.Sub(rental.ExpectedReturnDate).Hours()/48) + 1
		//isLate = true

		// Iterasi setiap item untuk menghitung late fee berdasarkan LateFeePerDay tiap mainan
		for _, rentalItem := range rental.RentalItems {
			// Ambil data mainan untuk mendapatkan LateFeePerDay
			toy, err := s.toyRepo.FindById(ctx, rentalItem.ToyID.String())
			if err != nil {
				return nil, errors.New("tidak dapat mendapatkan data mainan: " + rentalItem.ToyID.String())
			}

			// Hitung late fee: LateFeePerDay * jumlah hari terlambat * quantity
			itemLateFee := toy.LateFeePerDay * float64(days) * float64(rentalItem.Quantity)
			totalLateFee += itemLateFee
		}

		rental.Status = "overdue"
	} else {
		rental.Status = "completed"
	}

	rental.LateFee = totalLateFee

	// Proses setiap item
	var totalDamageFee float64 = 0
	//allReturned := true

	for _, itemReq := range req.Items {
		rentalItem, exists := rentalItemMap[itemReq.RentalItemID]
		if !exists {
			return nil, errors.New("item rental dengan ID " + itemReq.RentalItemID.String() + " tidak ditemukan")
		}

		// Validasi condition after
		validConditions := []string{"new", "excellent", "good", "fair", "poor", "damaged", "lost"}
		validCondition := false
		for _, c := range validConditions {
			if itemReq.ConditionAfter == c {
				validCondition = true
				break
			}
		}
		if !validCondition {
			return nil, errors.New("kondisi tidak valid: " + itemReq.ConditionAfter)
		}

		// Set kondisi setelah pengembalian
		rentalItem.ConditionAfter = itemReq.ConditionAfter
		rentalItem.DamageDescription = itemReq.DamageDescription

		// Hitung damage fee berdasarkan perubahan kondisi
		var damageFee float64 = 0

		// Ambil data mainan
		toy, _ := s.toyRepo.FindById(ctx, rentalItem.ToyID.String())

		if itemReq.ConditionAfter == "lost" {
			// Jika mainan hilang, kenakan biaya penuh replacement price
			damageFee = toy.ReplacementPrice * float64(rentalItem.Quantity)
			rentalItem.Status = "lost"
		} else if itemReq.ConditionAfter == "damaged" {
			// Jika rusak berat, kenakan 70% replacement price
			damageFee = toy.ReplacementPrice * 0.7 * float64(rentalItem.Quantity)
			rentalItem.Status = "damaged"
		} else {
			// Bandingkan kondisi sebelum dan sesudah
			conditionValues := map[string]int{
				"new":       5,
				"excellent": 4,
				"good":      3,
				"fair":      2,
				"poor":      1,
			}

			beforeValue := conditionValues[rentalItem.ConditionBefore]
			afterValue := conditionValues[itemReq.ConditionAfter]

			// Jika kondisi menurun, kenakan biaya berdasarkan selisih
			if afterValue < beforeValue {
				// Misalnya: 15% replacement price dikalikan selisih kondisi
				damageFee = toy.ReplacementPrice * 0.15 * float64(beforeValue-afterValue) * float64(rentalItem.Quantity)
			}

			rentalItem.Status = "returned"
		}

		rentalItem.DamageFee = damageFee
		totalDamageFee += damageFee

		// Update rental item
		if err := s.rentalRepo.UpdateRentalItem(ctx, rentalItem); err != nil {
			return nil, err
		}

		// Cek apakah semua item dikembalikan
		if rentalItem.Status != "returned" {
			//allReturned = false
		}
	}

	// Update total damage fee
	rental.DamageFee = totalDamageFee

	// Update notes jika ada
	if req.Notes != "" {
		rental.Notes = req.Notes
	}

	// Hitung total amount
	rental.TotalAmount = rental.TotalRentalPrice + rental.LateFee + rental.DamageFee

	// Simpan perubahan rental
	if err := s.rentalRepo.ReturnRental(ctx, &rental); err != nil {
		return nil, err
	}

	return &rental, nil
}
