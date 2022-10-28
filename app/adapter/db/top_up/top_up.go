package top_up

import "github.com/vins7/top-up-services/app/adapter/entity"

type TopUpRepo interface {
	UpdateBalance(interface{}) error
	InsertTopUp(in interface{}) error
	InsertTransactionHistory(req *entity.TransactionHistory) error
}
