package synthetics_test

import (
	"os"
	"testing"

	"github.com/dollarshaveclub/terraform-provider-nrs/pkg/synthetics"
)

func client() *synthetics.Client {
	conf := func(s *synthetics.Client) {
		s.APIKey = os.Getenv("NEW_RELIC_API_KEY")
	}
	client, err := synthetics.NewClient(conf)
	if err != nil {
		panic(err)
	}

	return client
}

func TestGetAllMonitors(t *testing.T) {
	response, err := client().GetAllMonitors()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetAllMonitors response: %#v", response)

	t.Logf("Count: %d", response.Count)
	for _, monitor := range response.Monitors {
		t.Logf("Monitor: %#v", monitor)
	}
}

func TestIntegration(t *testing.T) {
	args := &synthetics.CreateMonitorArgs{
		Name:         "david-test-1",
		Type:         "SCRIPT_BROWSER",
		Frequency:    60,
		URI:          "https://www.dollarshaveclub.com",
		Locations:    []string{"AWS_US_WEST_1"},
		Status:       "ENABLED",
		SLAThreshold: 7,
	}

	createMonitor, err := client().CreateMonitor(args)
	if err != nil {
		t.Fatal(err)
	}

	getMonitor, err := client().GetMonitor(createMonitor.ID)
	if err != nil {
		t.Fatal(err)
	}

	updateMonitor, err := client().UpdateMonitor(getMonitor.ID, &synthetics.UpdateMonitorArgs{
		Frequency: 60,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := client().UpdateMonitorScript(updateMonitor.ID, &synthetics.UpdateMonitorScriptArgs{
		ScriptText: "for {}",
	}); err != nil {
		t.Fatal(err)
	}

	script, err := client().GetMonitorScript(updateMonitor.ID)
	if err != nil {
		t.Fatal(err)
	}
	if script != "for {}" {
		t.Fatalf("invalid script returned: %s", script)
	}

	if err := client().DeleteMonitor(updateMonitor.ID); err != nil {
		t.Fatal(err)
	}
}
