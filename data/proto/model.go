package proto

import "time"

type Payment struct {
	TrxId     string
	Date      time.Time
	ToAddress string
	Amount    float64
}

type SplitPayment struct {
	TrxId         string
	Date          time.Time
	FleetAddress  string
	FleetAmount   float64
	DriverAddress string
	DriverAmount  float64
}

type Transfer struct {
	TrxId            string
	Date             time.Time
	FromOwnerAddress string
	FromAddress      string
	ToOwnerAddress   string
	ToAddress        string
	Amount           float64
}
