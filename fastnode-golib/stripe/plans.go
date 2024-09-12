package stripe

import (
	"log"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/errors"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/licensing"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/plan"
)

// PlanID is a string
type PlanID = string

// Plans contains the information about the different plan offered and a mutex to protect concurrent access
type Plans map[licensing.Plan]PlanID

// DefaultPlanIDs lists which Stripe IDs are preferred in case of duplicate fastnodePlanIDs
var DefaultPlanIDs = map[string]struct{}{
	"XXXXXXX": {}, // prod pro-yearly
	"XXXXXXX": {}, // prod pro-monthly
	"XXXXXXX": {}, // test pro-yearly
	"XXXXXXX": {}, // test pro-monthly
}

// PlansFromStripe fetches the plan mapping from Stripe
func PlansFromStripe() (Plans, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	plans := make(Plans)

	stripePlans := plan.List(&stripe.PlanListParams{
		Active: stripe.Bool(true),
	})
	for stripePlans.Next() {
		stripePlan := stripePlans.Plan()

		fastnodePlanID, ok := stripePlan.Metadata["fastnode_plan"]
		if !ok {
			continue
		}

		fastnodePlan := licensing.Plan(fastnodePlanID)
		if !fastnodePlan.IsValid() {
			return nil, errors.Errorf("invalid fastnode_plan %s", fastnodePlan)
		}
		if !fastnodePlan.IsPaid() {
			return nil, errors.Errorf("fastnode_plan %s is not a paid plan", fastnodePlan)
		}
		if _, ok := DefaultPlanIDs[stripePlan.ID]; ok {
			plans[fastnodePlan] = stripePlan.ID
		}
	}

	for _, fastnodePlan := range licensing.PaidPlans {
		if _, ok := plans[fastnodePlan]; !ok {
			log.Printf("no Stripe plan with fastnode_plan %s", fastnodePlan)
		}
	}

	return plans, nil
}

func checkInit() error {
	if stripe.Key == "" {
		return errors.New("Please init stripe before calling this function")
	}
	if octobatPublicKey == "" || octobatSecretKey == "" {
		return errors.New("Please make sure that both octobat public and private key are initialized before calling this function (using InitOctobat function)")
	}
	return nil
}
