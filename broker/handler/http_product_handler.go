package handler

import (
	"broker/proto"
	"broker/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{
		repo: repo,
	}
}

func (p *ProductHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/product", p.CreateProduct)
	r.GET("/product", p.ListProducts)
	r.PUT("/product", p.UpdateProduct)
	r.DELETE("/product", p.DeleteProduct)
}

func (p *ProductHandler) CreateProduct(c *gin.Context) {
	var payload proto.ProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if _, err := p.repo.Create(&proto.ProductRequest{
		Payload: &payload,
	}); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product created successfully"})
}

func (u *ProductHandler) ListProducts(c *gin.Context) {
	query := c.Query("page")
	page, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid page"})
		return
	}

	response, err := u.repo.ListProducts(&proto.Offset{Id: int32(page)})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, proto.ProductList{
		Page:      response.Page,
		Total:     response.Total,
		TotalPage: response.TotalPage,
		Products:  response.Products,
	})
}

func (u *ProductHandler) UpdateProduct(c *gin.Context) {
	var product proto.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := c.Query("id")
	ID, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	product.Id = uint32(ID)

	response, err := u.repo.UpdateProduct(&product)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

func (u *ProductHandler) DeleteProduct(c *gin.Context) {
	query := c.Query("id")
	ID, err := strconv.Atoi(query)

	_, err = u.repo.DeleteProduct(&proto.GetProductRequest{Id: int32(ID)})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}
