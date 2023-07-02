package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nextasy01/go-stripe/api/entity"
	"github.com/Nextasy01/go-stripe/api/handlers"
	"github.com/stripe/stripe-go/v74"
)

type prodConfig struct {
	*handlers.StripeConfig
}

func NewProdConfig(sc *handlers.StripeConfig) *prodConfig {
	return &prodConfig{sc}
}

func (sc *prodConfig) CreateProduct(w http.ResponseWriter, r *http.Request) {
	sc.Log.Println("Handling POST create product request")
	prod := entity.NewProduct()
	e := json.NewDecoder(r.Body)
	e.Decode(&prod)

	if prod.Name == "" {
		http.Error(w, "Please provide the name of the product!", http.StatusBadRequest)
		return
	}

	params := &stripe.ProductParams{
		Name:        stripe.String(prod.Name),
		Description: stripe.String(prod.Description),
		Images:      stripe.StringSlice(prod.Images),
	}
	p, err := sc.Api.Products.New(params)
	if err != nil {
		http.Error(w, "Unable to create a product", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Unable to marshal product", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	sc.Log.Println("Created a product: " + p.Name)
	fmt.Fprintf(w, "Successfuly created a product: %s", string(data))
}
