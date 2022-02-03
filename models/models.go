package models

type ExpectedYield_Averange struct {
	Currency string `json:"curency"`
	Value    uint   `json:"value"`
}

// type AverangePositionPrice struct {
// 	Currency string `json:"curency"`
// 	Value    uint   `json:"value"`
// }

// type AverangePositionPriceNoNkd struct {
// 	Currency string `json:"curency"`
// 	Value    uint   `json:"value"`
// }

type Position struct {
	Figi                       string `json:"figi"`
	Ticker                     string `json:"ticker"`
	Isin                       string `json:"isin"`
	InstrumentType             string `json:"instumentType"`
	Balance                    uint   `json:"balance"`
	Blocked                    uint   `json:"blocked"`
	ExpectedYield              ExpectedYield_Averange
	Lots                       uint `json:"lots"`
	AverangePositionPrice      ExpectedYield_Averange
	AverangePositionPriceNoNkd ExpectedYield_Averange
	Name                       string `json:"name"`
}

type Instruments struct {
	Figi              string  `json:"figi"`
	Ticker            string  `json:"ticker"`
	Isin              string  `json:"isin"`
	MinPriceIncrement float32 `json:"minPriceIncrement"`
	Lot               uint    `json:"lot"`
	MinQuantity       float32 `json:"minQuantity"`
	Currency          string  `json:"currency"`
	Name              string  `json:"name"`
	Type              string  `json:"type"`
}

type Bids_Asks struct {
	Price    float32
	Quantity float32 `json:"Quantity"`
}

type PayloadPortfolio struct {
	Position Position
}

type PayloadStocks struct {
	Instruments []Instruments
	Total       uint `json:"total"`
}

type PayloadStockByFigi struct {
	Figi              string `json:"figi"`
	Depth             uint8
	Bids              []Bids_Asks
	Asks              []Bids_Asks
	TradeStatus       string
	MinPriceIncrement float32 `json:"minPriceIncrement"`
	FaceValue         float32
	LastPrice         float32
	ClosePrice        float32
	LimitUp           float32
	LimitDown         float32
}

type PayloadError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Portfolio struct {
	TrackingID string `json:"trackingId"`
	Status     string `json:"status"`
	Payload    PayloadPortfolio
}

type Stocks struct {
	TrackingID string `json:"trackingId"`
	Payload    PayloadStocks
	Status     string `json:"status"`
}

type StocksByFigi struct {
	TrackingID string `json:"trackingId"`
	Payload    PayloadStockByFigi
	Status     string `json:"status"`
}

type NotFoundError struct {
	TrackingID   string       `json:"trackingId"`
	Status       string       `json:"status"`
	PayloadError PayloadError `json:"payload"`
}
