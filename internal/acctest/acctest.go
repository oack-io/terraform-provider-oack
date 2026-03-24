package acctest

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/oack-io/terraform-provider-oack/internal/provider"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"oack": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("OACK_API_KEY") == "" {
		t.Fatal("OACK_API_KEY must be set for acceptance tests")
	}
	if os.Getenv("OACK_ACCOUNT_ID") == "" {
		t.Fatal("OACK_ACCOUNT_ID must be set for acceptance tests")
	}
}
