package logic

type GetRateRequest struct {
	FromCurrency string
	ToCurrency   string
}

type GetRateResponse struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	Timestamp    int64
}
