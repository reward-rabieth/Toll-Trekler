package types

type OBUData struct {
	OBUID int     `json:"obuid"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
	//Vin   string  `json:"vin"`
}

type Distance struct {
	Value float64 `json:"value"`
	OBUID int     `json:"OBUID"`
	Unix  int64   `json:"unix"`
}

type Invoice struct {
	OBUID         int     `json:"OBUID"`
	TotalDistance float64 `json:"total_distance"`
	TotalAmount   float64 `json:"total_amount"`
}
