package domain

type Remains struct {
	GoodsName  string  `json:"goods_name" db:"goods_name"`
	Producer   string  `json:"producer" db:"producer"`
	Country    string  `json:"country" db:"country"`
	MNN        string  `json:"mnn" db:"mnn"`
	Price      float32 `json:"price" db:"price"`
	Remain     string  `json:"remain" db:"remain"`
	Contractor string  `json:"contractor" db:"contractor"`
	Series     string  `json:"series" db:"series"`
	BestBefore string  `json:"best_before" db:"best_before"`
	Store      string  `json:"store" db:"store"`
}

type RemainRequest struct {
	GoodsName string `json:"goods_name" db:"goods_name"`
	Producer  string `json:"producer" db:"producer"`
	MNN       string `json:"mnn" db:"mnn"`
}
