package price

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nextasy01/go-stripe/api/entity"
	"github.com/Nextasy01/go-stripe/api/handlers"
	"github.com/Nextasy01/go-stripe/api/validators"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"
)

type priceConfig struct {
	*handlers.StripeConfig
}

func NewPriceConfig(sc *handlers.StripeConfig) *priceConfig {
	return &priceConfig{sc}
}

func (sc *priceConfig) CreatePrice(w http.ResponseWriter, r *http.Request) {
	sc.Log.Println("Handling POST create product request")
	price_req := entity.NewPrice()
	e := json.NewDecoder(r.Body)
	e.Decode(&price_req)

	if err := validators.ValidateCurrency(price_req); err != nil {
		http.Error(w, "Please provide the price of the product!", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	params := &stripe.PriceParams{
		Currency:   stripe.String(price_req.Currency),
		Product:    stripe.String(price_req.Product),
		UnitAmount: stripe.Int64(int64(*price_req.Amount)),
	}
	p, err := price.New(params)

	if err != nil {
		http.Error(w, "Unable to create a price of a product", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Unable to marshal price", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	sc.Log.Println("Created price of a product: " + p.ID)
	fmt.Fprintf(w, "Successfuly created price of a product: %s", string(data))

}
