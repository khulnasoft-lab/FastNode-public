package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/community/account"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/community/account/checkout"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/envutil"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/stripe"
)

func main() {
	stripeSecret := envutil.MustGetenv("STRIPE_SECRET")
	stripeWebhookSecret := envutil.MustGetenv("STRIPE_WEBHOOK_SECRET")
	octobatSecret := envutil.MustGetenv("OCTOBAT_SECRET")
	octobatPublishable := envutil.MustGetenv("OCTOBAT_PUBLISHABLE")
	beanieConfigID := envutil.MustGetenv("BEANIE_CONFIG_ID")

	account.InitStripe(stripeSecret, stripeWebhookSecret, octobatSecret, octobatPublishable, beanieConfigID)
	plans, err := stripe.PlansFromStripe()
	if err != nil {
		log.Fatalf("Error while fetching plan information from stripe: %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	checkoutRouter := router.PathPrefix("/api/checkout").Subrouter()
	checkout.SetupRoutes(plans, checkoutRouter)

	const port int = 9090
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), router))
}
