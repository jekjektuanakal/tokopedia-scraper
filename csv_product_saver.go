package main

import (
	"encoding/csv"
	"fmt"
	"io"
)

type CSVProductSaver struct {
	writer io.Writer
}

func NewCSVProductSaver(writer io.Writer) *CSVProductSaver {
	return &CSVProductSaver{
		writer: writer,
	}
}

func (s *CSVProductSaver) SaveProducts(products []Product) error {
	if len(products) == 0 {
		return fmt.Errorf("product is empty")
	}

	csvWriter := csv.NewWriter(s.writer)

	defer csvWriter.Flush()

	header := []string{
		"ID",
		"URL",
		"ImageURL",
		"Name",
		"Price",
		"Rating",
		"ShopID",
		"ShopName",
	}

	err := csvWriter.Write(header)
	if err != nil {
		return err
	}

	defer csvWriter.Flush()

	for _, product := range products {
		err = csvWriter.Write([]string{
			fmt.Sprintf("%d", product.ID),
			product.URL,
			product.ImageURL,
			product.Name,
			fmt.Sprintf("%d", product.Price),
			fmt.Sprintf("%f", product.Rating),
			fmt.Sprintf("%d", product.Shop.ID),
			product.Shop.Name,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
