package main

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

type SearchQuery struct {
	Category string
	Rows     int
	Start    int
	Page     int
}

type Shop struct {
	ID   int64
	Name string
}

type Product struct {
	ID       int64
	URL      string
	ImageURL string
	Name     string
	Price    int64
	Rating   float64
	Shop     Shop
}

type ProductList struct {
	Count    int
	Products []Product
}

type ProductSearcher interface {
	SearchProduct(SearchQuery) (ProductList, error)
}

type ProductSaver interface {
	SaveProducts(products []Product) error
}

type ScrapeQuery struct {
	Category    string
	NumProducts int
}

type ProductScraperConfig struct {
	NumSearchWorker int
}

type ProductScraper struct {
	searcher        ProductSearcher
	numSearchWorker int
	saver           ProductSaver
}

func NewProductScraper(cfg ProductScraperConfig, searcher ProductSearcher, saver ProductSaver) *ProductScraper {
	if cfg.NumSearchWorker <= 0 {
		cfg.NumSearchWorker = 1
	}

	return &ProductScraper{
		searcher:        searcher,
		numSearchWorker: cfg.NumSearchWorker,
		saver:           saver,
	}
}

func (s *ProductScraper) ScrapeProduct(query ScrapeQuery) error {
	if query.Category == "" {
		return fmt.Errorf("category must not be empty")
	}

	if query.NumProducts <= 0 {
		return fmt.Errorf("num products must be greater than 0")
	}

	g := new(errgroup.Group)
	productListQueue := make(chan ProductList, s.numSearchWorker)
	remainder := query.NumProducts % s.numSearchWorker

	for i := 0; i < s.numSearchWorker; i++ {
		extra := 0

		if remainder > 0 {
			extra = 1
			remainder--
		}

		g.Go(func() error {
			productList, err := s.searcher.SearchProduct(SearchQuery{
				Category: query.Category,
				Rows:     (query.NumProducts / s.numSearchWorker) + extra,
				Start:    1,
				Page:     i + 1,
			})
			if err != nil {
				return err
			}

			productListQueue <- productList

			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return fmt.Errorf("error when search product: %w", err)
	}

	close(productListQueue)

	products := make([]Product, 0)

	// save products to csv file
	for productList := range productListQueue {
		products = append(products, productList.Products...)
	}

	err = s.saver.SaveProducts(products)
	if err != nil {
		return err
	}

	return nil
}
