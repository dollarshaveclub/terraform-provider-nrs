package provider_test

import (
	"testing"

	"github.com/dollarshaveclub/terraform-provider-nrs/pkg/provider"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := provider.Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
