package article

type Meta struct {
	Took            int `json:"took"`
	Page            int `json:"page"`
	TotalPage       int `json:"totalPage"`
	TotalData       int `json:"totalData"`
	TotalDataOnPage int `json:"totalDataOnPage"`
}
