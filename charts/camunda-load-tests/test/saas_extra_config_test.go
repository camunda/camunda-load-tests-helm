// Copyright 2026 Camunda Services GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// TestGoldenSaasWithExtraConfig exercises the `saas.enabled=true` +
// `global.extraConfig` combination, which previously caused the
// camunda-load-tests.annotations helper to emit both checksum annotations
// on a single line, breaking YAML parsing of the rendered Deployments.
func TestGoldenSaasWithExtraConfig(t *testing.T) {
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"workers", "starter"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: "golden-saas-with-extra-config-" + name,
			Templates:      []string{"templates/" + name + ".yaml"},
			// checksum/saas-credentials hashes the rendered credentials.yaml, which embeds
			// the chart version via the helm.sh/chart label. Ignore it so a chart version
			// bump doesn't require regenerating this golden file.
			IgnoredLines: []string{`checksum/saas-credentials:.*`},
			SetValues: map[string]string{
				"saas.enabled":                                             "true",
				"saas.credentials.clientId":                                "clientId",
				"saas.credentials.clientSecret":                            "clientSecret",
				"saas.credentials.authServer":                              "authServer",
				"saas.credentials.zeebeRestAddress":                        "zeebeRestAddress",
				"saas.credentials.authType":                                "authType",
				"saas.credentials.zeebeGrpcAddress":                        "zeebeGrpcAddress",
				"saas.credentials.authorizationAudience":                   "authorizationAudience",
				"global.extraConfig.load-tester.monitor-data-availability": "false",
			},
		})
	}
}
