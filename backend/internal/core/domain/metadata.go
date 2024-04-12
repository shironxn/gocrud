package domain

type Metadata struct {
	Sort         string `query:"sort" json:"sort" validate:"oneof=id user_id name title created_at updated_at"`
	Order        string `query:"order" json:"order" validate:"oneof=asc desc"`
	TotalRecords int64  `json:"total_records"`
	TotalPage    int    `json:"total_page"`
	Limit        int    `query:"limit" json:"limit"`
	Page         int    `query:"page" json:"page"`
}
