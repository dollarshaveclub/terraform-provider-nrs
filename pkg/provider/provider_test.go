package provider_test

import (
	"testing"

	"github.com/dollarshaveclub/terraform-provider-nrs/pkg/provider"
)

func TestProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
