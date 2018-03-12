package ticket

type DataBank struct {
	Data   *InfoBank `json:"data"`
	Status string    `json:"status"`
}

type InfoBank struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Address     string  `json:"address"`
	BranchID    string  `json:branch_id`
	CountPeople int     `json:"count_people"`
}

type DataAllTimeCreate struct {
}
