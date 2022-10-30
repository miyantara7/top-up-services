package top_up

import (
	"strconv"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	client "github.com/vins7/top-up-services/app/adapter/client"
	db "github.com/vins7/top-up-services/app/adapter/db/top_up"
	"github.com/vins7/top-up-services/app/adapter/entity"
	"github.com/vins7/top-up-services/app/interface/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TopUpUsecase struct {
	client client.EMoneyRepo
	db     db.TopUpRepo
}

func NewTopUpUsecase(client client.EMoneyRepo, db db.TopUpRepo) TopUp {
	return &TopUpUsecase{
		client: client,
		db:     db,
	}
}

func (u *TopUpUsecase) TopUpBalance(in interface{}) error {
	var req *model.TopUpRequest
	var errTrxHist error
	var errUpdateBal error
	wt := new(sync.WaitGroup)

	if err := mapstructure.Decode(in, &req); err != nil {
		return err
	}

	b, err := u.client.GetBalance(&model.GetBalance{
		UserId:   req.UserId,
		UserName: req.UserName,
		NoKartu:  req.NoKartu,
	})
	if err != nil {
		return err
	}

	saldo, err := strconv.Atoi(b.Balance)
	if err != nil {
		return err
	}

	amount, err := strconv.Atoi(req.Balance)
	if err != nil {
		return err
	}

	vBalance := saldo + amount

	dataBalance := &entity.EMoney{
		NoKartu: req.NoKartu,
		UserId:  req.UserId,
		Balance: &vBalance,
	}
	wt.Add(2)
	go func() {
		defer wt.Done()
		err := make(chan error, 1)
		UpdateBalance(u.db, dataBalance, err)
		errUpdateBal = <-err
	}()

	go func() {
		defer wt.Done()
		err := make(chan error, 1)
		InsertTransactionHistory(u.db, &entity.TransactionHistory{
			UserId:      req.UserId,
			NoKartu:     req.NoKartu,
			CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
			UpdateDate:  time.Now().Format("2006-01-02 15:04:05"),
			Setor:       req.Balance,
			Tarik:       "0",
			Balance:     strconv.Itoa(vBalance),
		}, err)
		errTrxHist = <-err
	}()

	wt.Wait()
	if errTrxHist != nil {
		return errTrxHist
	}
	if errUpdateBal != nil {
		return errUpdateBal
	}

	return nil
}

func (u *TopUpUsecase) Payment(in interface{}) error {
	var req *model.PaymentRequest
	var errTrxHist error
	var errUpdateBal error
	var errInsertTopUp error
	wt := new(sync.WaitGroup)

	if err := mapstructure.Decode(in, &req); err != nil {
		return err
	}

	biller, err := u.client.DetailBiller(&model.BillerRequest{
		ID: req.BillerId,
	})
	if err != nil {
		return err
	}

	b, err := u.client.GetBalance(&model.GetBalance{
		UserId:   req.UserId,
		UserName: req.UserName,
		NoKartu:  req.NoKartu,
	})
	if err != nil {
		return err
	}

	saldo, err := strconv.Atoi(b.Balance)
	if err != nil {
		return err
	}

	amount, err := strconv.Atoi(biller.Price)
	if err != nil {
		return err
	}

	fee, err := strconv.Atoi(biller.Fee)
	if err != nil {
		return err
	}

	vTotal := amount + fee
	if vTotal > saldo {
		return status.Errorf(codes.Internal, "Saldo anda tidak mencukupi!")
	}

	vBalance := saldo - vTotal
	dataBalance := &entity.EMoney{
		NoKartu: req.NoKartu,
		UserId:  req.UserId,
		Balance: &vBalance,
	}

	wt.Add(3)
	go func() {
		defer wt.Done()
		err := make(chan error, 1)
		UpdateBalance(u.db, dataBalance, err)
		errUpdateBal = <-err
	}()

	go func() {
		defer wt.Done()
		err := make(chan error, 1)
		InsertTopUp(u.db, &entity.TopUp{
			UserId:  req.UserId,
			NoKartu: req.NoKartu,
			Product: biller.Product,
			Price:   biller.Price,
			Fee:     biller.Fee,
		}, err)
		errInsertTopUp = <-err
	}()

	go func() {
		defer wt.Done()
		err := make(chan error, 1)
		InsertTransactionHistory(u.db, &entity.TransactionHistory{
			UserId:      req.UserId,
			NoKartu:     req.NoKartu,
			CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
			UpdateDate:  time.Now().Format("2006-01-02 15:04:05"),
			Setor:       "0",
			Tarik:       strconv.Itoa(vTotal),
			Balance:     strconv.Itoa(*dataBalance.Balance),
		}, err)
		errTrxHist = <-err
	}()

	wt.Wait()
	if errTrxHist != nil {
		return errTrxHist
	}
	if errInsertTopUp != nil {
		return errInsertTopUp
	}
	if errUpdateBal != nil {
		return errUpdateBal
	}

	return nil
}

func InsertTransactionHistory(dbs db.TopUpRepo, req *entity.TransactionHistory, errChan chan error) {
	errChan <- dbs.InsertTransactionHistory(req)
}

func UpdateBalance(dbs db.TopUpRepo, req *entity.EMoney, errChan chan error) {
	errChan <- dbs.UpdateBalance(req)
}

func InsertTopUp(dbs db.TopUpRepo, req *entity.TopUp, errChan chan error) {
	errChan <- dbs.InsertTopUp(req)
}
