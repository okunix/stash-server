package dto

type Page struct {
	Limit  uint  `json:"limit"`
	Offset uint  `json:"offset"`
	Total  int64 `json:"total"`
}
