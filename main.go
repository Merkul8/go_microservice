package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Данный микросервис выводит товары продавцов с наибольшим количеством просмотров

type Product struct {
	ID           int  `json:"id"`
	Name         string `json:"name"`
	ProductCode  string `json:"product_code"`
	Price        string `json:"price"`
	ProductCount int  `json:"product_count"`
	Slug         string `json:"slug"`
	IsStock      bool `json:"is_stock"`
	SellerIDID   int  `json:"seller_id_id"`
	Views        int  `json:"views"`
}

func main() {
	file, err := os.Open("../marketplace/products.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var products []Product
	err = json.Unmarshal(bytes, &products)
	if err != nil {
		panic(err)
	}
	max := 0
	var name string
	for _, product := range products {
		if max < product.Views {
			max = product.Views
			name = product.Name
		}
	}
	fmt.Println("Товар -", name, "наиболее просматриваемый, просмотров -", max)
}
