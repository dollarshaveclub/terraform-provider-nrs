.PHONY: terraform-provider-nrs test

VERSION = "0.1.0"

terraform-provider-nrs:
	go install github.com/dollarshaveclub/terraform-provider-nrs/cmd/terraform-provider-nrs

test:
	go test -v github.com/dollarshaveclub/terraform-provider-nrs/pkg/... \
		github.com/dollarshaveclub/terraform-provider-nrs/cmd/...

create-release:
	git diff-index --quiet HEAD -- || echo "All changes must be commited first." && exit 1

	rm -rf build releases
	mkdir -p build/terraform-provider-nrs/linux-amd64
	mkdir -p build/terraform-provider-nrs/darwin-amd64

	GOOS=linux GOARCH=amd64 go build -o build/terraform-provider-nrs/linux-amd64/terraform-provider-nrs \
		github.com/dollarshaveclub/terraform-provider-nrs/cmd/terraform-provider-nrs

	GOOS=darwin GOARCH=amd64 go build -o build/terraform-provider-nrs/darwin-amd64/terraform-provider-nrs \
		github.com/dollarshaveclub/terraform-provider-nrs/cmd/terraform-provider-nrs

	mkdir releases
	tar -c build/terraform-provider-nrs -C build | gzip -c > releases/terraform-provider-nrs-${VERSION}.tar.gz

	git tag ${VERSION} HEAD
