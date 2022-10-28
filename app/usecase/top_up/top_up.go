package top_up

type TopUp interface {
	TopUpBalance(interface{}) error
	Payment(interface{}) error
}
