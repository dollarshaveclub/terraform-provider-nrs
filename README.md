# terraform-provider-nrs

A Terraform provider for New Relic Synthetics. This provider can be
used to manage the lifecycle of Synthetics monitors and alert
condititions with Terraform.

# Installing

First, grab the latest release from Github and extract the binaries:

```
$ wget https://github.com/dollarshaveclub/terraform-provider-nrs/releases/download/0.1.0/terraform-provider-nrs-0.1.0.tar.gz
$ tar -xf terraform-provider-nrs-0.1.0.tar.gz
```

Within the release, there are two binaries, one for Linux and one for
MacOS:
- `terraform-provider-nrs/darwin-amd64/terraform-provider-nrs`
- `terraform-provider-nrs/linux-amd64/terraform-provider-nrs`

Take note of which is appropriate for the architecture of your system.

Copy the binary to `/usr/local/bin`:

```
// Substitute $ARCH with either darwin-amd64 or linux-amd64
$ cp terraform-provider-nrs/$ARCH/terraform-provider-nrs /usr/local/bin
```

Modify `~/.terraformrc` so that it has these contents:

```
providers {
    nrs = "/usr/local/bin/terraform-provider-nrs
}
```

The provider should now be useable from within Terraform!

# Example

```
provider "nrs" {
  newrelic_api_key = "REDACTED"
}

resource "nrs_monitor" "new_monitor" {
  name = "monitor_name"

  // The monitor's checking frequency in minutes (one of 1, 5, 10,
  // 15, 30, 60, 360, 720, or 1440).
  frequency = 60

  // The monitoring locations. A list can be found at the endpoint:
  // https://synthetics.newrelic.com/synthetics/api/v1/locations
  locations = ["AWS_US_WEST_1"]

  status = "ENABLED"

  // The type of monitor (one of SIMPLE, BROWSER, SCRIPT_API,
  // SCRIPT_BROWSER)
  type = "SCRIPT_BROWSER"

  sla_threshold = 7

  // The URI to check. This only applies to SIMPLE and BROWSER
  // monitors.
  uri = "https://www.dollarshaveclub.com"

  // The API or browser script to execute. This only applies to
  // SCRIPT_API or SCRIPT_BROWSER monitors. Docs can be found here:
  // https://docs.newrelic.com/docs/synthetics/new-relic-synthetics/scripting-monitors/write-scripted-browsers
  script = "console.log('this is a check!')"
}

resource "nrs_alert_condition" "new_condition" {
  name = "test-condition"
  monitor_id = "${nrs_monitor.new_monitor.id}"
  enabled = true
  policy_id = "${newrelic_alert_policy.new_policy.id}"
}
```

# Import

In case of using import for alerts condition be aware that provider needs two value to do correct import.

```
terraform import name_of_resource policy_id:condition_id
terraform import nrs_alert_condition.alert1 123456:567890
```

For importing monitors only id of monitor is needed:
```
terraform import name_of_resource monitor_id
terraform import nrs_monitor.monitor1 d02c69d5-bac8-4243-91f4-4f9c62a7c71c
```
