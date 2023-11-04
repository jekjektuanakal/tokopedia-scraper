package main

import (
	"os"

	"github.com/stretchr/testify/suite"
)

type CSVProductSaverTestSuite struct {
	suite.Suite
}

func (s *CSVProductSaverTestSuite) TestSaveProduct() {
	s.Run("product is empty, should return error", func() {
		csvFile, err := os.CreateTemp(".", "test_product_save.csv")
		if err != nil {
			panic(err)
		}

		defer os.Remove(csvFile.Name())

		saver := NewCSVProductSaver(csvFile)

		err = saver.SaveProducts([]Product{})

		s.Error(err)
	})

	s.Run("product is not empty, should return nil", func() {
		csvFile, err := os.CreateTemp(".", "test_product_save.csv")
		if err != nil {
			panic(err)
		}

		defer os.Remove(csvFile.Name())

		saver := NewCSVProductSaver(csvFile)

		err = saver.SaveProducts([]Product{
			{
				Name:  "product 1",
				Price: 1000,
			},
			{
				Name:  "product 2",
				Price: 2000,
			},
		})

		s.NoError(err)
	})
}
