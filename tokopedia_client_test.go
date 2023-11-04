package main

import (
	"github.com/stretchr/testify/suite"
)

type TokopediaClientTestSuite struct {
	suite.Suite
}

func (s *TokopediaClientTestSuite) TestSearchProduct() {
	// Create new tokopedia client
	client := NewTokopediaClient()

	s.Run("category is empty, should return error", func() {
		_, err := client.SearchProduct(SearchQuery{
			Rows:  20,
			Start: 1,
			Page:  1,
		})

		s.Error(err)
	})

	s.Run("category is not found, should return error", func() {
		_, err := client.SearchProduct(SearchQuery{
			Category: "not_found",
			Rows:     20,
			Start:    1,
			Page:     1,
		})

		s.Error(err)
	})

	s.Run("rows is empty, should return error", func() {
		_, err := client.SearchProduct(SearchQuery{
			Category: "handphone",
			Start:    1,
			Page:     1,
		})

		s.Error(err)
	})

	s.Run("rows is less than 0, should return error", func() {
		_, err := client.SearchProduct(SearchQuery{
			Category: "handphone",
			Rows:     -1,
			Start:    1,
			Page:     1,
		})

		s.Error(err)
	})

	s.Run("start is empty, should return error", func() {
		_, err := client.SearchProduct(SearchQuery{
			Category: "handphone",
			Rows:     20,
			Page:     1,
		})

		s.Error(err)
	})

	s.Run("start is less than 0, should return error", func() {
		_, err := client.SearchProduct(SearchQuery{
			Category: "handphone",
			Rows:     20,
			Start:    -1,
			Page:     1,
		})

		s.Error(err)
	})

	s.Run("page is empty, should return error", func() {
		_, err := client.SearchProduct(SearchQuery{
			Category: "handphone",
			Rows:     20,
			Start:    1,
		})

		s.Error(err)
	})

	s.Run("page is less than 0, should return error", func() {
		_, err := client.SearchProduct(SearchQuery{
			Category: "handphone",
			Rows:     20,
			Start:    1,
			Page:     -1,
		})

		s.Error(err)
	})

	s.Run("category, rows, start, and page is not empty, should return product list", func() {
		products, err := client.SearchProduct(SearchQuery{
			Category: "handphone",
			Rows:     20,
			Start:    1,
			Page:     1,
		})

		s.NoError(err)
		s.NotEmpty(products)
	})
}
