.PHONY: test

test:
	go test -v github.com/dollarshaveclub/terraform-provider-nrs/pkg/... \
		github.com/dollarshaveclub/terraform-provider-nrs/cmd/...
