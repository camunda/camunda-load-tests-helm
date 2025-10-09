package test

import (
	"camunda-load-tests/charts/camunda-load-test/test/golden"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestGoldenElasticExporterDefaults(t *testing.T) {
	// Test which allows to verify also parent chart templates
	// This makes sure that properties are correctly set
	// OR configurations have been changed

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"servicemonitor", "deployment"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test",
			GoldenFileName: "es-exporter-" + name,
			Templates:      []string{"charts/prometheus-elasticsearch-exporter/templates/" + name + ".yaml"},
		})
	}
}
