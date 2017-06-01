.PHONY: terraform-provider-nrs test

terraform-provider-nrs:
	go install github.com/dollarshaveclub/terraform-provider-nrs/cmd/terraform-provider-nrs

test:
	go test -v github.com/dollarshaveclub/terraform-provider-nrs/pkg/... \
		github.com/dollarshaveclub/terraform-provider-nrs/cmd/...
