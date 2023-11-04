package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/http2"
)

type TokopediaClient struct {
	client *http.Client
}

func NewTokopediaClient() *TokopediaClient {
	return &TokopediaClient{
		client: &http.Client{
			Transport: &http2.Transport{},
		},
	}
}

type productData struct {
	ProductList ProductList
}

type errorLocation struct {
	Line   int
	Column int
}

type errorExtensions struct {
	DeveloperMessage string
	MoreInfo         string
}

type Error struct {
	Message    string
	Locations  []errorLocation
	Extensions errorExtensions
}

type productResponse struct {
	Data   productData
	Errors []Error
}

var categoryMap = map[string]string{
	"handphone": "handphone-tablet_handphone",
}

func (c *TokopediaClient) SearchProduct(productQuery SearchQuery) (ProductList, error) {
	if productQuery.Category == "" {
		return ProductList{}, fmt.Errorf("category must not be empty")
	}

	var category string
	var ok bool

	if category, ok = categoryMap[productQuery.Category]; !ok {
		return ProductList{}, fmt.Errorf("category not found")
	}

	if productQuery.Rows <= 0 {
		return ProductList{}, fmt.Errorf("rows must be greater than 0")
	}

	if productQuery.Start <= 0 {
		return ProductList{}, fmt.Errorf("start must be greater than 0")
	}

	if productQuery.Page <= 0 {
		return ProductList{}, fmt.Errorf("page must be greater than 0")
	}

	httpRequest := &http.Request{
		Method: "POST",
		URL:    &url.URL{Scheme: "https", Host: "gql.tokopedia.com", Path: "/graphql/SearchProductQuery"},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
			"User-Agent":   []string{"Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0"},
			"Accept":       []string{"*/*"},
			"Origin":       []string{"https://www.tokopedia.com"},
			"Connection":   []string{"keep-alive"},
			"Referer":      []string{"https://www.tokopedia.com/search?st=product&q=handphone&navsource=home"},
		},
	}

	queryStatement := `
		query ($params: String) 
		{ 
			productList: searchProduct(params: $params) 
			{ 
				count 

				products 
				{ 
					id 
					url 
					imageUrl: image_url
					name 
					price: price_int
					rating 

					shop 
					{ 
						id 
						name
					}
				} 
			} 
		}`

	queryParams := []string{
		fmt.Sprintf("identifier=%s", categoryMap[category]),
		fmt.Sprintf("page=%d", productQuery.Page),
		fmt.Sprintf("rows=%d", productQuery.Rows),
		fmt.Sprintf("start=%d", productQuery.Start),
		"ob=23",
		"sc=24",
		"user_id=0",
		"source=directory",
		"device=desktop",
		"related=true",
		"st=product",
		"safe_search=false",
	}

	requestBody := map[string]any{
		"query": queryStatement,
		"variables": map[string]any{
			"params": strings.Join(queryParams, "&"),
		},
	}

	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return ProductList{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpRequest.Body = io.NopCloser(bytes.NewBuffer(jsonRequestBody))

	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return ProductList{}, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	var response productResponse

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return ProductList{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(response.Errors) > 0 {
		return ProductList{}, fmt.Errorf("failed to search product: %q", response.Errors[0].Message)
	}

	return response.Data.ProductList, nil
}
