package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Nextasy01/go-stripe/api/handlers"
	"github.com/Nextasy01/go-stripe/api/validators"
	"github.com/stripe/stripe-go/v74"
)

type product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	Currency    string    `json:"currency" validate:"required,currency"`
	Price       int64     `json:"price" validate:"required"`
	Images      []string  `json:"images,omitempty"`
	CreatedAt   time.Time `json:"created"`
}

type srvConfig struct {
	*handlers.StripeConfig
}

func NewSrvConfig(sc *handlers.StripeConfig) *srvConfig {
	return &srvConfig{sc}
}

func (sc *srvConfig) TestShortLink(w http.ResponseWriter, r *http.Request) {
	link, err := shortenLink("https://buy.stripe.com/test_4gw6q1fpx9kPbtubIK")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}
	sc.Log.Printf("https://%s", link)
	fmt.Fprintf(w, "Here is the payment link: %s", fmt.Sprintf("https://%s", link))

}

func (sc *srvConfig) CreateProduct(w http.ResponseWriter, r *http.Request) {
	sc.Log.Println("Handling POST create product request")
	prod := &product{}
	e := json.NewDecoder(r.Body)
	e.Decode(&prod)

	if prod.Description == "" {
		prod.Description = "none"
	}

	if err := validators.ValidateCurrency(prod); err != nil {
		http.Error(w, "Please provide the price of the product!", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	prodParams := &stripe.ProductParams{
		Name:        stripe.String(prod.Name),
		Description: stripe.String(prod.Description),
	}
	stripeProd, err := sc.Api.Products.New(prodParams)
	if err != nil {
		http.Error(w, "Unable to create a product", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	priceParams := &stripe.PriceParams{
		Currency:   stripe.String(prod.Currency),
		Product:    stripe.String(stripeProd.ID),
		UnitAmount: stripe.Int64(prod.Price),
	}
	stripePrice, err := sc.Api.Prices.New(priceParams)
	if err != nil {
		http.Error(w, "Unable to create a price of a product", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	link, err := sc.generatePaymentLink(stripePrice.ID)
	if err != nil {
		http.Error(w, "Unable to create a payment link", http.StatusBadRequest)
		sc.Log.Printf("Error occured: %v", err)
		return
	}

	newLink, err := shortenLink(link)
	if err != nil {
		sc.Log.Printf("Couldn't shorten a link: %v", err)
	} else {
		link = fmt.Sprintf("https://%s", newLink)
	}

	sc.Log.Println("Successfully created a product!")
	fmt.Fprintf(w, "Successfuly created a product! Here is the payment link: %s", link)
}

func (sc *srvConfig) generatePaymentLink(price_id string) (string, error) {

	params := &stripe.PaymentLinkParams{
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				Price:    stripe.String(price_id),
				Quantity: stripe.Int64(1),
			},
		},
	}
	result, err := sc.Api.PaymentLinks.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

func shortenLink(url string) (string, error) {
	body, err := json.Marshal(struct {
		URL string `json:"url"`
	}{url})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://mini-sh.onrender.com/api/v1", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Println("Response from shortener service: " + string(respBody))

	var shortLink struct {
		Data struct {
			Url string `json:"short"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &shortLink)
	if err != nil {
		return "", err
	}

	log.Println(shortLink)

	return shortLink.Data.Url, err
}
