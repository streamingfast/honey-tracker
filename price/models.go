package price

type Response struct {
	Data struct {
		Items []*HistoricalPrice `json:"items"`
	} `json:"data"`
	Success bool `json:"success"`
}

type HistoricalPrice struct {
	UnixTime int64   `json:"unixTime"`
	Value    float64 `json:"value"`
}
