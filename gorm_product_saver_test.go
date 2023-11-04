package main

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GORMProductSaverTestSuite struct {
	suite.Suite
}

func (s *GORMProductSaverTestSuite) TestSaveProduct() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	s.Run("product is empty, should return error", func() {
		saver := NewGORMProductSaver(db)

		err = saver.SaveProducts([]Product{})

		s.Error(err)
	})

	s.Run("product is not empty, should return nil", func() {
		saver := NewGORMProductSaver(db)

		err = saver.SaveProducts([]Product{
			{
				ID:    1,
				Name:  "product 1",
				Price: 1000,
			},
			{
				ID:    2,
				Name:  "product 2",
				Price: 2000,
			},
		})

		var products []productDAO

		db.Find(&products)

		s.NoError(err)
		s.Equal(2, len(products))

		db.Delete(products)
	})

	s.Run("product exists, should update product", func() {
		saver := NewGORMProductSaver(db)

		err = saver.SaveProducts([]Product{
			{
				ID:    1,
				Name:  "product 1",
				Price: 1000,
			},
			{
				ID:    2,
				Name:  "product 2",
				Price: 2000,
			},
		})

		err = saver.SaveProducts([]Product{
			{
				ID:    1,
				Name:  "product 1",
				Price: 1500,
			},
			{
				ID:    2,
				Name:  "product 2",
				Price: 3000,
			},
		})

		products := make([]productDAO, 0)

		db.Find(&products)

		s.NoError(err)
		s.Equal(2, len(products))

		s.Equal(int64(1500), products[0].Price)
		s.Equal(int64(3000), products[1].Price)

		db.Delete(products)
	})
}
