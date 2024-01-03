package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// curl -X POST -d "{\"username\":\"your_username\",\"password\":\"your_password\"}" -H "Content-Type: application/json" http://127.0.0.1:8000/api-token-auth/


/* 
Мне нужно сделать форму получения от продавца username and password, в ответ я буду выдавать отчет о их товарах,
 будет посылаться запрос на получение токена, который указан выше и далее запрос на API маркетплейса для получения 
 товаров. Затем будет происходить процесс генерации отчета, который в последствии будет отправлен обратно на маркетплейс 
*/ 

type Product struct {
	ID         int  `json:"id"`
	Name       string `json:"name"`
	ProductCode string `json:"product_code"`
	Price      string `json:"price"`
	ProductCount int  `json:"product_count"`
	Slug       string `json:"slug"`
	IsStock    bool `json:"is_stock"`
	SellerIDID int  `json:"seller_id_id"`
	Views      int  `json:"views"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func getToken(username, password string) (string, error) {
	credentials := &Credentials{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://127.0.0.1:8000/api-token-auth/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.Token, nil
}

func getProducts(token string) (*[]Product, error) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/api/products/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Token " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}


	return &products, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("create_a_report.html")
		tmpl.Execute(w, nil)
	} else {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		token, err := getToken(username, password)
		if err != nil {
			fmt.Println("Failed to get token:", err)
			return
		}

		products, err := getProducts(token)
		if err != nil {
			fmt.Println("Failed to get products:", err)
			return
		}

		fmt.Println("Products:", products)
	}
}

func main() {
	http.HandleFunc("/create_a_report", loginHandler)
	http.ListenAndServe(":8080", nil)
}
