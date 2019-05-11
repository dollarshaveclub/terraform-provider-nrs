package provider

import (
	"fmt"
	"strconv"
	"strings"

	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

// NRSAlertConditionResource returns a Terraform schema for a New
// Relic Synthetics alert condition.
func NRSAlertConditionResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alert condition",
			},
			"monitor_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the monitor",
				ForceNew:    true,
			},
			"runbook_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL to a runbook for addressing the alert",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the alert condition is enabled",
			},
			"policy_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the policy to attach the alert condition to",
				ForceNew:    true,
			},
		},
		Create: NRSAlertConditionCreate,
		Exists: NRSAlertConditionExists,
		Delete: NRSAlertConditionDelete,
		Read:   NRSAlertConditionRead,
		Update: NRSAlertConditionUpdate,
		Importer: &schema.ResourceImporter{
			State: NRSAlertConditionImportState,
		},
	}
}

func schemaId(resourceData *schema.ResourceData) (int, error) {
	sresid := resourceData.Id()

	iresid, err := strconv.Atoi(sresid)
	if (err != nil) {
		return -1, fmt.Errorf("error: could not determine/convert id from resourceData.Id()=\"%s\" to int.", sresid)
	}

	return iresid, nil
}


// NRSAlertConditionCreate creates a Synthetics alert condition using
// Terraform configuration.
func NRSAlertConditionCreate(resourceData *schema.ResourceData, meta interface{}) error {
	client := meta.(*synthetics.Client)

	args := &synthetics.CreateAlertConditionArgs{
		Name:      resourceData.Get("name").(string),
		MonitorID: resourceData.Get("monitor_id").(string),
		Enabled:   resourceData.Get("enabled").(bool),
	}
	if data, ok := resourceData.GetOk("runbook_url"); ok {
		args.RunbookURL = data.(string)
	}

	alertCondition, err := client.CreateAlertCondition(uint(resourceData.Get("policy_id").(int)), args)
	if err != nil {
		return errors.Wrapf(err, "error: could not create alert condition")
	}

	resourceData.SetId(fmt.Sprintf("%d", alertCondition.ID))

	return nil
}

// NRSAlertConditionImportState imports given condition to Terraform state
// using policy_id and condition_id from New Relic alerts API
func NRSAlertConditionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		/*
		   In New Relic alert condition , we need both policy_id  and condition_id to get given user data.
		*/
		return nil, fmt.Errorf("Import resource ID should consist of policy_id:condition_id")
	}

	val, _ := strconv.Atoi(s[0])
	d.SetId(s[1])
	d.Set("policy_id", val)

	return []*schema.ResourceData{d}, nil
}

// NRSAlertConditionExists checks whether an alert condition exists
// using Terraform configuration.
func NRSAlertConditionExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*synthetics.Client)

	iresid, err := schemaId(resourceData)
	if err != nil {
		return false, err
	}

	_, err = client.GetAlertCondition(uint(resourceData.Get("policy_id").(int)), uint(iresid))
	if err == synthetics.ErrAlertConditionNotFound {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "error: could not find alert condition")
	}

	return true, nil
}

// NRSAlertConditionDelete deletes a Synthetics alert condition using
// Terraform configuration.
func NRSAlertConditionDelete(resourceData *schema.ResourceData, meta interface{}) error {
	client := meta.(*synthetics.Client)

	iresid, err := schemaId(resourceData)
	if err != nil {
		return err
	}

	if err = client.DeleteAlertCondition(uint(iresid)); err != nil {
		return errors.Wrap(err, "error: could not delete alert condition")
	}

	return nil
}

// NRSAlertConditionRead refreshes alert condition information using
// Terraform configuration.
func NRSAlertConditionRead(resourceData *schema.ResourceData, meta interface{}) error {
	client := meta.(*synthetics.Client)

	iresid, err := schemaId(resourceData)
	if err != nil {
		return err
	}

	ac, err := client.GetAlertCondition(uint(resourceData.Get("policy_id").(int)), uint(iresid))
	if err != nil {
		return errors.Wrapf(err, "error: could not find alert condition")
	}

	if err := resourceData.Set("name", ac.Name); err != nil {
		return err
	}
	if err := resourceData.Set("monitor_id", ac.MonitorID); err != nil {
		return err
	}
	if err := resourceData.Set("runbook_url", ac.RunbookURL); err != nil {
		return err
	}
	if err := resourceData.Set("enabled", ac.Enabled); err != nil {
		return err
	}

	return nil
}

// NRSAlertConditionUpdate updates a Synthetics alert condition using
// Terraform configuration.
func NRSAlertConditionUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	client := meta.(*synthetics.Client)

	args := &synthetics.UpdateAlertConditionArgs{
		Name:      resourceData.Get("name").(string),
		MonitorID: resourceData.Get("monitor_id").(string),
		Enabled:   resourceData.Get("enabled").(bool),
	}
	if resourceData.HasChange("runbook_url") {
		args.RunbookURL = resourceData.Get("runbook_url").(string)
	}

	_, err := client.UpdateAlertCondition(uint(resourceData.Get("policy_id").(int)), args)
	if err != nil {
		return errors.Wrapf(err, "error: could not update alert condition")
	}

	return nil
}
