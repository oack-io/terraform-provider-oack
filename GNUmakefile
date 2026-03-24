default: build

build:
	go build -o terraform-provider-oack

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/oack-io/oack/0.1.0/darwin_arm64
	cp terraform-provider-oack ~/.terraform.d/plugins/registry.terraform.io/oack-io/oack/0.1.0/darwin_arm64/terraform-provider-oack_v0.1.0

clean:
	rm -f terraform-provider-oack

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS)

.PHONY: default build install clean testacc
