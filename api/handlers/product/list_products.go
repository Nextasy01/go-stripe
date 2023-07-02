package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stripe/stripe-go/v74"
)

func (sc *prodConfig) GetAll(w http.ResponseWriter, r *http.Request) {
	sc.Log.Println("Handling GET all products request")
	p := []*stripe.Product{}
	params := &stripe.ProductListParams{}
	iter := sc.Api.Products.List(params)

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

// func convert(products []*stripe.Product) []*entity.Product {
// 	p := []*entity.Product{}

// 	for _, v := range products {
// 		temp := new(entity.Product)
// 		temp.ID = v.ID
// 		temp.Name = v.Name
// 		temp.Description = v.Description
// 		temp.CreatedAt = time.Unix(v.Created, 0)
// 		temp.UpdatedAt = time.Unix(v.Updated, 0)
// 		if len(v.Images) != 0 {
// 			temp.Images = append(temp.Images, v.Images...)
// 		} else {
// 			temp.Images = nil
// 		}

// 		p = append(p, temp)
// 	}

// 	return p
// }
