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
			SetValues: map[string]string{
				"saas.enabled":                           "true",
				"saas.credentials.clientId":              "clientId",
				"saas.credentials.clientSecret":          "clientSecret",
				"saas.credentials.authServer":            "authServer",
				"saas.credentials.zeebeRestAddress":      "zeebeRestAddress",
				"saas.credentials.authType":              "authType",
				"saas.credentials.zeebeGrpcAddress":      "zeebeGrpcAddress",
				"saas.credentials.authorizationAudience": "authorizationAudience",
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
				"saas.enabled":                    "true",
				"saas.credentials.existingSecret": "secret-name",
			},
		})
	}
}
