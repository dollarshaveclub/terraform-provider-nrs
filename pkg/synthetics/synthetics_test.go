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

func TestGetMonitor(t *testing.T) {
	monitor, err := client().GetMonitor("25f88215-8905-40d6-93df-4c34720531ca")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Monitor: %#v", monitor)
}

func TestCreateMonitor(t *testing.T) {
	t.Skip()

	args := &synthetics.CreateMonitorArgs{
		Name:         "david-test-1",
		Type:         "SIMPLE",
		Frequency:    60,
		URI:          "https://www.dollarshaveclub.com",
		Locations:    []string{"AWS_US_WEST_1"},
		Status:       "ENABLED",
		SLAThreshold: 7,
	}

	monitor, err := client().CreateMonitor(args)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Monitor: %#v", monitor)
}
