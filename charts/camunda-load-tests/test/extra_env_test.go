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

func TestGoldenExtraEnvs(t *testing.T) {
	// Test to validate that extra environments are correctly rendered in templates
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"workers", "starter"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: "golden-extra-envs-" + name,
			Templates:      []string{"templates/" + name + ".yaml"},
			SetValues: map[string]string{
				"global.extraEnvVars[0].name":  "APP_MONITORDATAAVAILABILITY",
				"global.extraEnvVars[0].value": "false",
			},
		})
	}
}
