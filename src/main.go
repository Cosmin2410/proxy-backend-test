package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/Cosmin2410/proxy-backend-test/src/handler"
	"github.com/Cosmin2410/proxy-backend-test/src/limit"
	"github.com/Cosmin2410/proxy-backend-test/src/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("SQL_DATABASE_PASSWORD")

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	err = db.AutoMigrate(&model.SaveLog{})
	if err != nil {
		log.Fatal("Auto migrate failed", err)
	}

	proxy, err := NewProxy("https://jsonplaceholder.typicode.com")
	if err != nil {
		panic(err)
	}

	http.Handle("/", limit.RateLimiter(ProxyRequestHandler(proxy)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		req.Host = req.URL.Host
		originalDirector(req)
		modifyRequest(req)
	}

	h := &handler.DBCreate{DB: db}
	proxy.ModifyResponse = h.ModifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}

func modifyRequest(req *http.Request) {
	req.Header.Del("If-Modified-Since")
	req.Header.Del("If-None-Match")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
