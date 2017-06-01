package synthetics_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/dollarshaveclub/terraform-provider-nrs/pkg/synthetics"
)

func client() *synthetics.Client {
	conf := func(s *synthetics.Client) {
		if apiKey := os.Getenv("NEW_RELIC_API_KEY"); apiKey != "" {
			s.APIKey = apiKey
		} else {
			s.APIKey = "NEW_RELIC_API_KEY"
		}
	}
	client, err := synthetics.NewClient(conf)
	if err != nil {
		panic(err)
	}
	return client
}

func TestGetMonitor(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "asdf-asdf-asdf-asdf"

	monitor := &synthetics.Monitor{
		ID:           id,
		Name:         "test-monitor",
		Type:         "SCRIPT_BROWSER",
		Frequency:    60,
		URI:          "https://dollarshaveclub.com",
		Locations:    []string{"AWS_US_WEST_1"},
		Status:       "ENABLED",
		SLAThreshold: 7,
	}

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://synthetics.newrelic.com/synthetics/api/v3/monitors/%s", id),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, monitor)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	monitorResponse, err := client().GetMonitor(id)
	if err != nil {
		t.Fatal(err)
	}

	if monitor.ID != monitorResponse.ID ||
		monitor.Name != monitorResponse.Name ||
		monitor.Type != monitorResponse.Type ||
		monitor.Frequency != monitorResponse.Frequency ||
		monitor.URI != monitorResponse.URI ||
		monitor.Locations[0] != monitorResponse.Locations[0] ||
		monitor.Status != monitorResponse.Status ||
		monitor.SLAThreshold != monitorResponse.SLAThreshold {
		t.Fatal("monitor response is incorrect")
	}
}

func TestGetAllMonitors(t *testing.T) {
	t.Skip()

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
	// t.Skip()

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
