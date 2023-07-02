package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Nextasy01/go-stripe/api/handlers"
	"github.com/Nextasy01/go-stripe/api/handlers/product"
	"github.com/Nextasy01/go-stripe/api/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	l              = log.New(os.Stdout, "go-stripe-api ", log.LstdFlags)
	stripeHandler  handlers.StripeConfig
	productHandler = product.NewProdConfig(&stripeHandler)
	srv            = services.NewSrvConfig(&stripeHandler)
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	stripeHandler = handlers.NewStripeConfig(os.Getenv("Stripe_KEY"), l)
}

func main() {
	r := mux.NewRouter()

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/products", productHandler.GetAll)
	getRouter.HandleFunc("/api/products/search", productHandler.SearchProduct)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/api/products", productHandler.CreateProduct)
	postRouter.HandleFunc("/api/products/{id:[a-zA-Z0-9_]+}", productHandler.UpdateProduct)

	postRouter.HandleFunc("/api/v2/products", srv.CreateProduct)
	postRouter.HandleFunc("/api/v2/products/test", srv.TestShortLink)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/api/products/{id:[a-zA-Z0-9_]+}", productHandler.DeleteProduct)

	l.Println("Starting server at 8000 port")
	log.Fatal(http.ListenAndServe(":8000", r))

}
