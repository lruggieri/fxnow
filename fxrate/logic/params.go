package logic

type GetRateRequest struct {
	Pairs []string
}

type GetRateResponseRate struct {
	Pair      string
	Rate      float64
	Timestamp int64
}

type GetRateResponse struct {
	Rates []GetRateResponseRate
}
