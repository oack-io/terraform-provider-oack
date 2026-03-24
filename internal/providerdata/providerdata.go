package providerdata

import oack "github.com/oack-io/oack-go"

// Data holds the shared client and account context for all resources and datasources.
type Data struct {
	Client    *oack.Client
	AccountID string
}
