package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTokopediaClient(t *testing.T) {
	suite.Run(t, new(TokopediaClientTestSuite))
}

func TestProductScraper(t *testing.T) {
	suite.Run(t, new(ProductScraperTestSuite))
}

func TestCSVProductSaver(t *testing.T) {
	suite.Run(t, new(CSVProductSaverTestSuite))
}

func TestGORMProductSaver(t *testing.T) {
	suite.Run(t, new(GORMProductSaverTestSuite))
}
