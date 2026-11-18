package test

import (
	"camunda-load-tests/charts/camunda-load-tests/test/golden"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestGoldenCredentials(t *testing.T) {
	// Test which allows to verify also parent chart templates
	// This makes sure that properties are correctly set
	// OR configurations have been changed

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"credentials", "workers", "starter"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: "golden-credentials-" + name,
			Templates:      []string{"templates/" + name + ".yaml"},
			// checksum/saas-credentials hashes the rendered credentials.yaml, which embeds
			// the chart version via the helm.sh/chart label. Ignore it so a chart version
			// bump doesn't require regenerating this golden file.
			IgnoredLines: []string{`checksum/saas-credentials:.*`},
			SetValues: map[string]string{
				"oidc.enabled":                           "true",
				"oidc.credentials.clientId":              "clientId",
				"oidc.credentials.clientSecret":          "clientSecret",
				"oidc.credentials.authServer":            "authServer",
				"oidc.credentials.zeebeRestAddress":      "zeebeRestAddress",
				"oidc.credentials.authType":              "authType",
				"oidc.credentials.zeebeGrpcAddress":      "zeebeGrpcAddress",
				"oidc.credentials.authorizationAudience": "authorizationAudience",
			},
		})
	}
}

func TestGoldenCredentialsExistingSecret(t *testing.T) {
	// Test which allows to verify also parent chart templates
	// This makes sure that properties are correctly set
	// OR configurations have been changed

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"workers", "starter"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: "golden-existing-credential-secret-" + name,
			Templates:      []string{"templates/" + name + ".yaml"},
			SetValues: map[string]string{
				"oidc.enabled":                    "true",
				"oidc.credentials.existingSecret": "secret-name",
			},
		})
	}
}
