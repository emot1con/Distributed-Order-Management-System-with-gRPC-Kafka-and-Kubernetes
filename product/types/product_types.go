package types

// import (
// 	"time"
// )

// type Product struct {
// 	ID          uint      `gorm:"primaryKey" json:"id"`
// 	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
// 	Description string    `gorm:"type:text" json:"description"`
// 	Price       float64   `gorm:"not null" json:"price"`
// 	Stock       int       `gorm:"not null" json:"stock"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// type ProductPayload struct {
// 	Name        string  `json:"name" binding:"required"`
// 	Description string  `json:"description"`
// 	Price       float64 `json:"price" binding:"required,gt=0"`
// 	Stock       int     `json:"stock" binding:"required,gte=0"`
// }
