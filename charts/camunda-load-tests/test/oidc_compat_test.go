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
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func renderStarterDeployment(t *testing.T, setValues map[string]string) appsv1.Deployment {
	t.Helper()

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", "load-test-"+strings.ToLower(random.UniqueId())),
		SetValues:      setValues,
	}
	output := helm.RenderTemplate(t, options, chartPath, "load-test", []string{"templates/starter.yaml"})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)
	return deployment
}

func starterEnv(deployment appsv1.Deployment) map[string]string {
	env := make(map[string]string)
	for _, variable := range deployment.Spec.Template.Spec.Containers[0].Env {
		env[variable.Name] = variable.Value
	}
	return env
}

func starterSecretNames(deployment appsv1.Deployment) map[string]bool {
	names := make(map[string]bool)
	for _, variable := range deployment.Spec.Template.Spec.Containers[0].Env {
		if variable.ValueFrom != nil && variable.ValueFrom.SecretKeyRef != nil {
			names[variable.ValueFrom.SecretKeyRef.Name] = true
		}
	}
	return names
}

func TestLegacySaasKeysStillEnableOidc(t *testing.T) {
	deployment := renderStarterDeployment(t, map[string]string{
		"saas.enabled":                    "true",
		"saas.credentials.existingSecret": "legacy-secret",
	})

	env := starterEnv(deployment)
	assert.Equal(t, "saas", env["CAMUNDA_CLIENT_MODE"])
	assert.Equal(t, "oidc", env["CAMUNDA_CLIENT_AUTH_METHOD"])
	assert.Empty(t, env["CAMUNDA_CLIENT_GRPC_ADDRESS"])
	assert.True(t, starterSecretNames(deployment)["legacy-secret"])
}

func TestOidcKeysTakePrecedenceOverSaas(t *testing.T) {
	deployment := renderStarterDeployment(t, map[string]string{
		"saas.enabled":                    "true",
		"oidc.enabled":                    "true",
		"saas.credentials.existingSecret": "legacy-secret",
		"oidc.credentials.existingSecret": "new-secret",
	})

	env := starterEnv(deployment)
	assert.Equal(t, "saas", env["CAMUNDA_CLIENT_MODE"])
	assert.Equal(t, "oidc", env["CAMUNDA_CLIENT_AUTH_METHOD"])
	names := starterSecretNames(deployment)
	assert.True(t, names["new-secret"])
	assert.False(t, names["legacy-secret"])
}

func TestExplicitOidcFalseOverridesLegacySaasTrue(t *testing.T) {
	deployment := renderStarterDeployment(t, map[string]string{
		"saas.enabled": "true",
		"oidc.enabled": "false",
	})

	env := starterEnv(deployment)
	assert.Empty(t, env["CAMUNDA_CLIENT_MODE"])
	assert.NotEmpty(t, env["CAMUNDA_CLIENT_GRPC_ADDRESS"])
	assert.False(t, starterSecretNames(deployment)["load-test-credentials"])
}
