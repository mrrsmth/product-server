package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Model   string `gorm:"not null" json:"model"`
	Company string `gorm:"not null" json:"company"`
	Price   int    `json:"price"`
}

func main() {
	dsn := "root:123456@tcp(localhost:3306)/productdb?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Миграции таблицы Product
	db.AutoMigrate(&Product{})

	e := echo.New()

	// Подключение базы данных к обработчикам запросов
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	// Маршруты
	e.GET("/products", getProducts)
	e.GET("/products/:id", getProduct)
	e.POST("/products", createProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)

	// Запуск веб-сервера
	e.Logger.Fatal(e.Start(":8080"))
}

func getProducts(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	var products []Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, products)
}

func getProduct(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

func createProduct(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	var product Product
	if err := c.Bind(&product); err != nil {
		return err
	}

	if err := db.Create(&product).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, product)
}

func updateProduct(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	fmt.Println(id)

	var product Product
	if err := c.Bind(&product); err != nil {
		return err
	}

	if err := db.Save(&product).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

func deleteProduct(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := db.Delete(&Product{}, id).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
