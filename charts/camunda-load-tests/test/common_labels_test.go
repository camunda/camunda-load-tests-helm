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

func TestGoldenCommonLabels(t *testing.T) {
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	values := map[string]string{
		"global.commonLabels.camunda\\.io/created-by": "test-user",
		"global.commonLabels.camunda\\.io/purpose":    "load-test",
	}

	cases := []struct {
		goldenFile string
		templates  []string
	}{
		{"golden-common-labels-starter", []string{"templates/starter.yaml"}},
		{"golden-common-labels-workers", []string{"templates/workers.yaml"}},
		{"golden-common-labels-clients-service", []string{"templates/clients-service.yaml"}},
	}

	for _, tc := range cases {
		tc := tc
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: tc.goldenFile,
			Templates:      tc.templates,
			SetValues:      values,
		})
	}
}
