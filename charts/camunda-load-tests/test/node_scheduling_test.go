// Copyright 2026 Camunda Services GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

func TestGoldenNodeScheduling(t *testing.T) {
	// Test to validate that nodeSelector and tolerations are correctly rendered in templates
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"workers", "starter"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: "golden-node-scheduling-" + name,
			Templates:      []string{"templates/" + name + ".yaml"},
			SetValues: map[string]string{
				"global.nodeSelector.test":       "123",
				"global.tolerations[0].key":      "component",
				"global.tolerations[0].operator": "Equal",
				"global.tolerations[0].value":    "foobar",
				"global.tolerations[0].effect":   "NoSchedule",
			},
		})
	}
}
