package handler

import (
	"broker/auth"
	"broker/proto"
	"broker/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productRepo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{
		productRepo: repo,
	}
}

func (p *ProductHandler) RegisterRoutes(r *gin.Engine) {
	productRoutes := r.Group("/product")
	productRoutes.Use(auth.ProtectedEndpoint())

	productRoutes.POST("/", p.CreateProduct)
	productRoutes.GET("/", p.ListProducts)
	productRoutes.PUT("/", p.UpdateProduct)
	productRoutes.DELETE("/", p.DeleteProduct)
}

func (p *ProductHandler) CreateProduct(c *gin.Context) {
	var payload proto.ProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if _, err := p.productRepo.Create(&proto.ProductRequest{
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

	response, err := u.productRepo.ListProducts(&proto.Offset{Id: int32(page)})
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

	response, err := u.productRepo.UpdateProduct(&product)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

func (u *ProductHandler) DeleteProduct(c *gin.Context) {
	query := c.Query("id")
	ID, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(400, gin.H{"error": "InvalidID"})
		return
	}

	_, err = u.productRepo.DeleteProduct(&proto.GetProductRequest{Id: int32(ID)})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}
