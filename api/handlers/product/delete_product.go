package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (sc *prodConfig) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	sc.Log.Println("Handling DELETE product request")
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

	p, err := sc.Api.Products.Del(id, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when deleting the product: %v", err), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Unable to marshal product", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Deleted product: %s", string(data))
	sc.Log.Println("Deleted product id = " + id)
}
