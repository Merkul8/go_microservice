package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"log"
	"html/template"
	"net/http"
	"github.com/xuri/excelize/v2"
)

// curl -X POST -d "{\"username\":\"your_username\",\"password\":\"your_password\"}" -H "Content-Type: application/json" http://127.0.0.1:8000/api-token-auth/


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

		// Генерируем отчет после получения продуктов
		err = generateReport(*products)
		if err != nil {
			fmt.Println("Failed to generate report:", err)
			return
		}

		// Редирект на страницу /download_report
		http.Redirect(w, r, "/download_report", 302)
	}
}

func generateReport(products []Product) error {
	file := excelize.NewFile()
	headers := []string{"ID", "Name", "Product Code", "Price", "Product Count", "Slug", "Is Stock", "Seller ID", "Views"}
	for i, header := range headers {
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string('A'+i), 1), header)
	}

	for i, product := range products {
		dataRow := i + 2
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", dataRow), product.ID)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", dataRow), product.Name)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", dataRow), product.ProductCode)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", dataRow), product.Price)
		file.SetCellValue("Sheet1", fmt.Sprintf("E%d", dataRow), product.ProductCount)
		file.SetCellValue("Sheet1", fmt.Sprintf("F%d", dataRow), product.Slug)
		file.SetCellValue("Sheet1", fmt.Sprintf("G%d", dataRow), product.IsStock)
		file.SetCellValue("Sheet1", fmt.Sprintf("H%d", dataRow), product.SellerIDID)
		file.SetCellValue("Sheet1", fmt.Sprintf("I%d", dataRow), product.Views)
	}

	if err := file.SaveAs("report.xlsx"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func downloadReportHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("report.xlsx")
	if err != nil {
		http.Error(w, "File not found", 404)
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	http.ServeContent(w, r, "report.xlsx", time.Now(), file)

}



func main() {
	http.HandleFunc("/create_a_report", loginHandler)
	http.HandleFunc("/download_report", downloadReportHandler)
	http.ListenAndServe(":8080", nil)
}
