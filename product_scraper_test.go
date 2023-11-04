package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	NumSearchWorker = 5
)

type ProductScraperTestSuite struct {
	suite.Suite
}

func (s *ProductScraperTestSuite) TestScrapeProduct() {
	searcher := &productSearcherFake{}
	saver := &productSaverFake{}
	scraper := NewProductScraper(ProductScraperConfig{NumSearchWorker: NumSearchWorker}, searcher, saver)

	s.Run("category is empty, should return error", func() {
		searcher.callCount = 0
		searcher.err = nil

		saver.callCount = 0
		saver.err = nil

		err := scraper.ScrapeProduct(ScrapeQuery{
			NumProducts: 20,
		})

		s.Error(err)
	})

	s.Run("num products is less than 1, should return error", func() {
		searcher.callCount = 0
		searcher.err = nil

		saver.callCount = 0
		saver.err = nil

		err := scraper.ScrapeProduct(ScrapeQuery{
			Category:    "handphone",
			NumProducts: 0,
		})

		s.Error(err)
	})

	s.Run("searcher return error, saver is not executed", func() {
		searcher.callCount = 0
		searcher.err = fmt.Errorf("some error")

		saver.callCount = 0
		saver.err = nil

		err := scraper.ScrapeProduct(ScrapeQuery{
			Category:    "handphone",
			NumProducts: 5,
		})

		s.Error(err)
		s.GreaterOrEqual(searcher.callCount, NumSearchWorker)
		s.Equal(0, saver.callCount)
	})

	s.Run("searcher return success but saver return error, return error", func() {
		searcher.callCount = 0
		searcher.err = nil

		saver.callCount = 0
		saver.err = fmt.Errorf("some error")

		err := scraper.ScrapeProduct(ScrapeQuery{
			Category:    "handphone",
			NumProducts: 5,
		})

		s.Error(err)
		s.GreaterOrEqual(searcher.callCount, NumSearchWorker)
		s.Equal(1, saver.callCount)
	})

	s.Run("searcher return success and saver return success, return nil", func() {
		searcher.callCount = 0
		searcher.err = nil

		saver.callCount = 0
		saver.err = nil

		err := scraper.ScrapeProduct(ScrapeQuery{
			Category:    "handphone",
			NumProducts: 23,
		})

		s.NoError(err)
		s.GreaterOrEqual(searcher.callCount, NumSearchWorker)
		s.Equal(1, saver.callCount)
		s.Equal(23, len(saver.lastSavedProducts))

		var mapProduct = make(map[int64]bool)

		for i := 0; i < 23; i++ {
			if _, ok := mapProduct[saver.lastSavedProducts[i].ID]; ok {
				s.Fail("Duplicate product ID")
			}
		}
	})

	s.Run("run multiple searcher with one searcher return error, return error", func() {
		searcher.callCount = 0
		searcher.err = nil
		searcher.randomErr = fmt.Errorf("some error")

		saver.callCount = 0
		saver.err = nil

		err := scraper.ScrapeProduct(ScrapeQuery{
			Category:    "handphone",
			NumProducts: 5,
		})

		s.Error(err)
		s.GreaterOrEqual(searcher.callCount, NumSearchWorker)
		s.Equal(0, saver.callCount)
	})
}

type productSearcherFake struct {
	err       error
	randomErr error
	callCount int
	mu        sync.Mutex
}

func (s *productSearcherFake) SearchProduct(productQuery SearchQuery) (ProductList, error) {
	s.mu.Lock()
	s.callCount++
	s.mu.Unlock()

	time.Sleep(time.Duration(rand.Int63()%50) * time.Millisecond) // simulate network latency

	if s.err != nil {
		return ProductList{}, s.err
	}

	if s.randomErr != nil {
		if s.callCount%5 == 0 {
			return ProductList{}, s.randomErr
		}
	}

	products := make([]Product, 0, productQuery.Rows)

	for i := 0; i < productQuery.Rows; i++ {
		products = append(products, Product{
			ID:       int64(productQuery.Page*productQuery.Rows) + int64(i) + int64(productQuery.Start),
			URL:      fmt.Sprintf("https://www.tokopedia.com/%d", i),
			ImageURL: fmt.Sprintf("https://ecs7.tokopedia.net/img/cache/some-image-%d.jpg", i),
			Name:     fmt.Sprintf("Product %d", i),
			Price:    rand.Int63n(1000000) + 1000000,
			Rating:   float64(i),
			Shop: Shop{
				ID:   int64(i),
				Name: fmt.Sprintf("Shop %d", i),
			},
		})
	}

	return ProductList{
			Count:    productQuery.Rows,
			Products: products,
		},
		nil
}

type productSaverFake struct {
	lastSavedProducts []Product
	err               error
	callCount         int
}

func (s *productSaverFake) SaveProducts(products []Product) error {
	s.callCount++

	time.Sleep(time.Duration(rand.Int63()%50) * time.Millisecond) // simulate network latency

	if s.err == nil {
		s.lastSavedProducts = products
	} else {
		s.lastSavedProducts = nil
	}

	return s.err
}
