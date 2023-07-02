package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stripe/stripe-go/v74"
)

func (sc *prodConfig) SearchProduct(w http.ResponseWriter, r *http.Request) {
	sc.Log.Println("Handling GET search product request")
	if r.URL.Query().Get("name") == "" {
		http.Error(w, "Please specify the name of the product you want to find", http.StatusBadRequest)
		return
	}
	p := []*stripe.Product{}
	params := &stripe.ProductSearchParams{}
	params.Query = *stripe.String(fmt.Sprintf("name: '%s'", r.URL.Query().Get("name")))

	iter := sc.Api.Products.Search(params)
	for iter.Next() {
		p = append(p, iter.Product())
	}

	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Unable to marshal product list", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}
	sc.Log.Println("Returned list of products")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Your products: %s", string(data))
}
