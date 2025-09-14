package types

type ProductDTO struct {
	Id          uint32  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProductListDTO struct {
	Total     int64         `json:"total"`
	Page      int32         `json:"page"`
	TotalPage int32         `json:"total_page"`
	Products  []*ProductDTO `json:"products"`
}
