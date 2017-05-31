package provider

import (
	"crypto/sha256"

	"github.com/dollarshaveclub/terraform-provider-nrs/pkg/synthetics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

// Provider returns a new New Relic Synthetics Terraform provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"new_relic_api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "An admin API key for New Relic",
				Sensitive:   true,
			},
		},
		ConfigureFunc: getClient,
		ResourcesMap: map[string]*schema.Resource{
			"nrs_monitor": NRSMonitorResource(),
		},
	}
}

func getClient(rd *schema.ResourceData) (interface{}, error) {
	apiKey, ok := rd.Get("new_relic_api_key").(string)
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

// NRSMonitorResource returns a Terraform schema for a New Relic
// Synthetics monitor.
func NRSMonitorResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The monitor's ID with New Relic",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"frequency": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				InputDefault: "60",
				Description:  "The monitor's checking frequency in minutes (one of 1, 5, 10, 15, 30, 60, 360, 720, or 1440",
			},
			"uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL to monitor",
			},
			"locations": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The locations to check from",
				Elem:        schema.TypeString,
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				InputDefault: "ENABLED",
				Description:  "The monitor's status (one of ENABLED, MUTED, DISABLED)",
			},
			"sla_threshold": &schema.Schema{
				Type:        schema.TypeFloat,
				Description: "The monitor's SLA threshold",
				Optional:    true,
			},
			"validation_string": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The monitor's validation string",
				Optional:    true,
			},
			"verify_ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Verify SSL",
				Optional:    true,
			},
			"bypass_head_request": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Bypass HEAD request",
				Optional:    true,
			},
			"treat_redirect_as_failure": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Treat redirect as failure",
				Optional:    true,
			},
			"script": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The script to execute",
				Optional:    true,
				StateFunc: func(i interface{}) string {
					s := i.(string)
					hash := sha256.New()
					hash.Write([]byte(s))
					return string(hash.Sum(nil))
				},
			},
			"script_locations": &schema.Schema{
				Type:        schema.TypeList,
				Description: "The private locations to execute the script from",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The name of the private location",
							Optional:    true,
						},
						"hmac": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The HMAC for the private location",
							Optional:    true,
						},
					},
				},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The type of monitor (one of SIMPLE, BROWSER, SCRIPT_API, SCRIPT_BROWSER)",
				Required:    true,
				ForceNew:    true,
			},
		},
		Create: func(resourceData *schema.ResourceData, meta interface{}) error {
			client := meta.(*synthetics.Client)

			args := &synthetics.CreateMonitorArgs{
				Name:         resourceData.Get("name").(string),
				Type:         resourceData.Get("type").(string),
				Frequency:    resourceData.Get("frequency").(uint),
				URI:          resourceData.Get("uri").(string),
				Locations:    resourceData.Get("locations").([]string),
				Status:       resourceData.Get("status").(string),
				SLAThreshold: resourceData.Get("sla_threshold").(float64),
			}

			if data, ok := resourceData.GetOk("validation_string"); ok {
				args.ValidationString = strPtr(data.(string))
			}
			if data, ok := resourceData.GetOk("verify_ssl"); ok {
				args.VerifySSL = boolPtr(data.(bool))
			}
			if data, ok := resourceData.GetOk("bypass_head_request"); ok {
				args.BypassHEADRequest = boolPtr(data.(bool))
			}
			if data, ok := resourceData.GetOk("treat_redirect_as_failure"); ok {
				args.TreatRedirectAsFailure = boolPtr(data.(bool))
			}

			monitor, err := client.CreateMonitor(args)
			if err != nil {
				return errors.Wrapf(err, "error: could not create monitor")
			}

			resourceData.SetId(monitor.ID)

			// Set script if it was provided.
			if data, ok := resourceData.GetOk("script"); ok {
				args := &synthetics.UpdateMonitorScriptArgs{
					ScriptText: data.(string),
				}

				// Set script locations
				if data, ok := resourceData.GetOk("script_locations"); ok {
					scriptLocations := data.([]map[string]interface{})
					for _, scriptLocation := range scriptLocations {
						args.ScriptLocations = append(
							args.ScriptLocations,
							&synthetics.ScriptLocation{
								Name: scriptLocation["name"].(string),
								HMAC: scriptLocation["hmac"].(string),
							},
						)
					}
				}

				if err := client.UpdateMonitorScript(monitor.ID, args); err != nil {
					return errors.Wrap(err, "error: could not update monitor script")
				}
			}

			return nil
		},
		Exists: func(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
			client := meta.(*synthetics.Client)

			if _, err := client.GetMonitor(resourceData.Id()); err != nil {
				if err == synthetics.ErrMonitorNotFound {
					return false, nil
				}
				return false, errors.Wrap(err, "error: could not get monitor")
			}

			return true, nil
		},
		Delete: func(resourceData *schema.ResourceData, meta interface{}) error {
			client := meta.(*synthetics.Client)

			if err := client.DeleteMonitor(resourceData.Id()); err != nil {
				return errors.Wrap(err, "error: could not delete monitor")
			}

			return nil
		},
		Read: func(resourceData *schema.ResourceData, meta interface{}) error {
			client := meta.(*synthetics.Client)

			monitor, err := client.GetMonitor(resourceData.Id())
			if err != nil {
				return errors.Wrap(err, "error: could not get monitor")
			}

			script, err := client.GetMonitorScript(resourceData.Id())
			switch err {
			case synthetics.ErrMonitorScriptNotFound:
				if err := resourceData.Set("script", nil); err != nil {
					return err
				}
				if err := resourceData.Set("script_locations", nil); err != nil {
					return err
				}
			case nil:
				if err := resourceData.Set("script", script); err != nil {
					return err
				}
			default:
				return errors.Wrap(err, "error: could not get monitor script")
			}

			if err := resourceData.Set("name", monitor.Name); err != nil {
				return err
			}
			if err := resourceData.Set("type", monitor.Type); err != nil {
				return err
			}
			if err := resourceData.Set("frequency", monitor.Frequency); err != nil {
				return err
			}
			if err := resourceData.Set("uri", monitor.URI); err != nil {
				return err
			}
			if err := resourceData.Set("locations", monitor.Locations); err != nil {
				return err
			}
			if err := resourceData.Set("status", monitor.Status); err != nil {
				return err
			}
			if err := resourceData.Set("sla_threshold", monitor.SLAThreshold); err != nil {
				return err
			}

			if monitor.ValidationString != nil {
				if err := resourceData.Set("validation_string", *monitor.ValidationString); err != nil {
					return err
				}
			} else {
				if err := resourceData.Set("validation_string", nil); err != nil {
					return err
				}
			}

			if monitor.VerifySSL != nil {
				if err := resourceData.Set("verify_ssl", *monitor.VerifySSL); err != nil {
					return err
				}
			} else {
				if err := resourceData.Set("verify_ssl", nil); err != nil {
					return err
				}
			}

			if monitor.BypassHEADRequest != nil {
				if err := resourceData.Set("bypass_head_request", *monitor.BypassHEADRequest); err != nil {
					return err
				}
			} else {
				if err := resourceData.Set("bypass_head_request", nil); err != nil {
					return err
				}
			}

			if monitor.TreatRedirectAsFailure != nil {
				if err := resourceData.Set("treat_redirect_as_failure", *monitor.TreatRedirectAsFailure); err != nil {
					return err
				}
			} else {
				if err := resourceData.Set("treat_redirect_as_failure", nil); err != nil {
					return err
				}
			}

			return nil
		},
		Update: func(resourceData *schema.ResourceData, meta interface{}) error {
			client := meta.(*synthetics.Client)

			args := &synthetics.UpdateMonitorArgs{
				Name:         resourceData.Get("name").(string),
				Frequency:    resourceData.Get("frequency").(uint),
				URI:          resourceData.Get("uri").(string),
				Locations:    resourceData.Get("locations").([]string),
				Status:       resourceData.Get("status").(string),
				SLAThreshold: resourceData.Get("sla_threshold").(float64),
			}

			if resourceData.HasChange("validation_string") {
				validationString := resourceData.Get("validation_string").(string)
				if validationString != "" {
					args.ValidationString = strPtr(validationString)
				}
			}
			if resourceData.HasChange("verify_ssl") {
				args.VerifySSL = boolPtr(resourceData.Get("verify_ssl").(bool))
			}
			if resourceData.HasChange("bypass_head_request") {
				args.BypassHEADRequest = boolPtr(resourceData.Get("bypass_head_request").(bool))
			}
			if resourceData.HasChange("treat_redirect_as_failure") {
				args.TreatRedirectAsFailure = boolPtr(resourceData.Get("treat_redirect_as_failure").(bool))
			}

			client.UpdateMonitor(resourceData.Id(), args)

			return nil
		},
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}
