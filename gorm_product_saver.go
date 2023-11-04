package main

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productDAO struct {
	gorm.Model
	URL      string
	ImageURL string
	Name     string
	Price    int64
	Rating   float64
	ShopID   int64
	ShopName string
}

type GORMProductSaver struct {
	db *gorm.DB
}

func NewGORMProductSaver(db *gorm.DB) *GORMProductSaver {
	err := db.AutoMigrate(&productDAO{})
	if err != nil {
		panic(err)
	}

	return &GORMProductSaver{db: db}
}

func (s *GORMProductSaver) SaveProducts(products []Product) error {
	if len(products) == 0 {
		return fmt.Errorf("product is empty")
	}

	daos := make([]productDAO, 0, len(products))

	for _, product := range products {
		daos = append(daos, productDAO{
			Model:    gorm.Model{ID: uint(product.ID), DeletedAt: gorm.DeletedAt{}},
			URL:      product.URL,
			ImageURL: product.ImageURL,
			Name:     product.Name,
			Price:    product.Price,
			Rating:   product.Rating,
			ShopID:   product.Shop.ID,
			ShopName: product.Shop.Name,
		})
	}

	return s.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"url",
				"image_url",
				"name",
				"price",
				"rating",
				"shop_id",
				"shop_name",
				"updated_at",
				"deleted_at",
			}),
		}).
		Create(&daos).Error
}
