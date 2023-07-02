package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v74"
)

func (sc *prodConfig) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	sc.Log.Println("Handling PUT product request")
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Please provide the id of the product!", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(id, "prod_") {
		http.Error(w, "Please provide valid id", http.StatusBadRequest)
		return
	}

	params := &stripe.ProductParams{}
	e := json.NewDecoder(r.Body)
	e.Decode(&params)

	p, err := sc.Api.Products.Update(id, params)
	if err != nil {
		http.Error(w, "Unable to update a product", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Unable to marshal product", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Updated product: %s", string(data))
	sc.Log.Println("Updated product id = " + id)
}
