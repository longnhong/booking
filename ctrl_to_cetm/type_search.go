package ctrl_to_cetm

/* DATA BANKS */
type Data struct {
	Data   []*InfoBank `json:"data"`
	Status string      `json:"status"`
}

type DataBank struct {
	Data   *InfoBank `json:"data"`
	Status string    `json:"status"`
}

/* DATA BANKS */
type AddressBank struct {
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Address    string  `json:"address"`
	NameBranch string  `json:"name"`
}

/* DATA SERVICES */
type DataService struct {
	Data Service `json:"data"`
}

type DataServices struct {
	Data []*Service `json:"data"`
}

type LangCD struct {
	Eng string `json:"eng"`
	Es  string `json:"es"`
	Sp  string `json:"sp"`
	Vi  string `json:"vi"`
}
type Service struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	L10n LangCD `json:"l10n"`
}
