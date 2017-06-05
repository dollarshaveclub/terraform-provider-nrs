package provider

import (
	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

// Provider returns a new New Relic Synthetics Terraform provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"newrelic_api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "An admin API key for New Relic",
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_KEY", "key"),
			},
		},
		ConfigureFunc: getClient,
		ResourcesMap: map[string]*schema.Resource{
			"nrs_monitor":         NRSMonitorResource(),
			"nrs_alert_condition": NRSAlertConditionResource(),
		},
	}
}

func getClient(rd *schema.ResourceData) (interface{}, error) {
	apiKey, ok := rd.Get("newrelic_api_key").(string)
	if !ok {
		return nil, errors.New("invalid type for new relic api key")
	}

	conf := func(s *synthetics.Client) {
		s.APIKey = apiKey
	}
	client, err := synthetics.NewClient(conf)
	if err != nil {
		return nil, errors.Wrap(err, "error: could not instantiate synthetics client")
	}

	return client, nil
}
