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

func TestGoldenExtraConfigurations(t *testing.T) {
	// Test to validate that extra configurations are correctly rendered in templates
	// This makes sure that properties are correctly set in configmap entries but also
	// mounted on the workers and starter deployments

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"load-test-config", "workers", "starter"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: "golden-extra-configs-" + name,
			Templates:      []string{"templates/" + name + ".yaml"},
			SetValues: map[string]string{
				"global.extraConfig.app.starter.processId":       "foo",
				"global.extraConfig.app.monitorDataAvailability": "true",
			},
		})
	}
}
