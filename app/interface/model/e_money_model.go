package model

type GetBalance struct {
	UserId   string
	UserName string
	NoKartu  string
}

type BillerRequest struct {
	ID string
}

type TopUpRequest struct {
	UserId   string
	UserName string
	NoKartu  string
	Balance  string
}

type PaymentRequest struct {
	UserId   string
	UserName string
	NoKartu  string
	BillerId string
}
