package main

import (
	"flag"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var saveToDB bool
	var saveToCSV bool
	var dbDSN string
	var csvFileName string
	var numWorker int
	var category string

	flag.BoolVar(&saveToDB, "db", false, "Save data to the database")
	flag.BoolVar(&saveToCSV, "csv", true, "Save data to a CSV file")
	flag.StringVar(&dbDSN, "dsn", "postgresql://postgres@localhost:5432", "Database DSN (Data Source Name)")
	flag.StringVar(&csvFileName, "csv-file", "products.csv", "CSV file name")
	flag.IntVar(&numWorker, "worker", 5, "Number of search worker")
	flag.StringVar(&category, "category", "handphone", "Category to search")

	flag.Parse()

	var saver ProductSaver

	if saveToDB && saveToCSV {
		fmt.Println("Error: Choose only one option, either 'db' or 'csv'")
	} else if saveToDB {
		if dbDSN == "" {
			fmt.Println("Error: DB DSN is required when saving to the database")
		} else {
			fmt.Println("Saving data to the database")

			db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
			if err != nil {
				panic(err)
			}

			saver = NewGORMProductSaver(db)
		}
	} else if saveToCSV {
		if csvFileName == "" {
			fmt.Println("Error: CSV file name is required when saving to a CSV file")

			panic(fmt.Errorf("csv file name is required"))
		} else {
			fmt.Println("Saving data to a CSV file with file name:", csvFileName)

			csvFile, err := os.OpenFile(csvFileName, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				panic(err)
			}

			saver = NewCSVProductSaver(csvFile)
		}

		fmt.Println("Saving data to a CSV file...")
	} else {
		fmt.Println("Error: Choose an option, either 'db' or 'csv'")
	}

	searcher := NewTokopediaClient()

	cfg := ProductScraperConfig{
		NumSearchWorker: numWorker,
	}

	scraper := NewProductScraper(cfg, searcher, saver)

	err := scraper.ScrapeProduct(ScrapeQuery{
		Category:    "handphone",
		NumProducts: 100,
	})

	if err != nil {
		panic(err)
	}

	fmt.Print("done")
}
