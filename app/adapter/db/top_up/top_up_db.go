package top_up

import (
	"errors"

	"github.com/vins7/top-up-services/app/adapter/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type TopUpDB struct {
	db       *gorm.DB
	dbEMoney *gorm.DB
}

func NewTopUpDB(db *gorm.DB, dbEMoney *gorm.DB) *TopUpDB {
	return &TopUpDB{
		db:       db,
		dbEMoney: dbEMoney,
	}
}

func (u *TopUpDB) UpdateBalance(in interface{}) error {

	req, ok := in.(*entity.EMoney)
	if !ok {
		return status.Errorf(codes.NotFound, "Error while casting")
	}

	data := &entity.EMoney{}
	if err := u.dbEMoney.Debug().Where("user_id = ? and no_kartu = ?", req.UserId, req.NoKartu).
		First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		}
		return status.Errorf(codes.Internal, err.Error())
	}

	data.Balance = req.Balance
	if err := u.dbEMoney.Debug().Save(&data).Error; err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

func (u *TopUpDB) InsertTopUp(in interface{}) error {

	req, ok := in.(*entity.TopUp)
	if !ok {
		return status.Errorf(codes.NotFound, "Error while casting")
	}

	if err := u.db.Debug().Create(&req).Error; err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

func (u *TopUpDB) InsertTransactionHistory(req *entity.TransactionHistory) error {

	if err := u.dbEMoney.Debug().Create(req).Error; err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}
