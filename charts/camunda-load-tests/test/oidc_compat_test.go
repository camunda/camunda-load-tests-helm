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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func renderOptions(t *testing.T, setValues map[string]string, valuesFiles ...string) (*helm.Options, string) {
	t.Helper()

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	return &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", "load-test-"+strings.ToLower(random.UniqueId())),
		SetValues:      setValues,
		ValuesFiles:    valuesFiles,
	}, chartPath
}

func renderStarterDeployment(t *testing.T, setValues map[string]string, valuesFiles ...string) appsv1.Deployment {
	t.Helper()

	options, chartPath := renderOptions(t, setValues, valuesFiles...)
	output := helm.RenderTemplate(t, options, chartPath, "load-test", []string{"templates/starter.yaml"})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)
	return deployment
}

func renderWorkerDeployment(t *testing.T, setValues map[string]string, valuesFiles ...string) appsv1.Deployment {
	t.Helper()

	options, chartPath := renderOptions(t, setValues, valuesFiles...)
	output := helm.RenderTemplate(t, options, chartPath, "load-test", []string{"templates/workers.yaml"})

	var workerDoc string
	for _, doc := range strings.Split(output, "---") {
		if strings.Contains(doc, "kind: Deployment") {
			workerDoc = doc
			break
		}
	}
	require.NotEmpty(t, workerDoc)

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, workerDoc, &deployment)
	return deployment
}

func renderCredentialsSecret(t *testing.T, setValues map[string]string, valuesFiles ...string) corev1.Secret {
	t.Helper()

	options, chartPath := renderOptions(t, setValues, valuesFiles...)
	output := helm.RenderTemplate(t, options, chartPath, "load-test", []string{"templates/credentials.yaml"})

	var secret corev1.Secret
	helm.UnmarshalK8SYaml(t, output, &secret)
	return secret
}

func renderCredentialsSecretRaw(t *testing.T, setValues map[string]string, valuesFiles ...string) (string, error) {
	t.Helper()

	options, chartPath := renderOptions(t, setValues, valuesFiles...)
	return helm.RenderTemplateE(t, options, chartPath, "load-test", []string{"templates/credentials.yaml"})
}

func writeValuesFile(t *testing.T, content string) string {
	t.Helper()

	file := filepath.Join(t.TempDir(), "values.yaml")
	require.NoError(t, os.WriteFile(file, []byte(content), 0o600))
	return file
}

func containerEnv(deployment appsv1.Deployment) map[string]string {
	env := make(map[string]string)
	for _, variable := range deployment.Spec.Template.Spec.Containers[0].Env {
		env[variable.Name] = variable.Value
	}
	return env
}

func secretNames(deployment appsv1.Deployment) map[string]bool {
	names := make(map[string]bool)
	for _, variable := range deployment.Spec.Template.Spec.Containers[0].Env {
		if variable.ValueFrom != nil && variable.ValueFrom.SecretKeyRef != nil {
			names[variable.ValueFrom.SecretKeyRef.Name] = true
		}
	}
	return names
}

func TestLegacySaasKeysStillEnableOidc(t *testing.T) {
	values := map[string]string{
		"saas.enabled":                    "true",
		"saas.credentials.existingSecret": "legacy-secret",
	}

	for name, deployment := range map[string]appsv1.Deployment{
		"starter": renderStarterDeployment(t, values),
		"worker":  renderWorkerDeployment(t, values),
	} {
		t.Run(name, func(t *testing.T) {
			env := containerEnv(deployment)
			assert.Equal(t, "saas", env["CAMUNDA_CLIENT_MODE"])
			assert.Equal(t, "oidc", env["CAMUNDA_CLIENT_AUTH_METHOD"])
			assert.Empty(t, env["CAMUNDA_CLIENT_GRPC_ADDRESS"])
			assert.True(t, secretNames(deployment)["legacy-secret"])
			assert.NotEmpty(t, deployment.Spec.Template.Annotations["checksum/saas-credentials"])
		})
	}
}

func TestOidcKeysTakePrecedenceOverSaas(t *testing.T) {
	values := map[string]string{
		"saas.enabled":                    "true",
		"oidc.enabled":                    "true",
		"saas.credentials.existingSecret": "legacy-secret",
		"oidc.credentials.existingSecret": "new-secret",
	}

	for name, deployment := range map[string]appsv1.Deployment{
		"starter": renderStarterDeployment(t, values),
		"worker":  renderWorkerDeployment(t, values),
	} {
		t.Run(name, func(t *testing.T) {
			env := containerEnv(deployment)
			assert.Equal(t, "saas", env["CAMUNDA_CLIENT_MODE"])
			assert.Equal(t, "oidc", env["CAMUNDA_CLIENT_AUTH_METHOD"])
			names := secretNames(deployment)
			assert.True(t, names["new-secret"])
			assert.False(t, names["legacy-secret"])
		})
	}
}

func TestExplicitOidcFalseOverridesLegacySaasTrue(t *testing.T) {
	values := map[string]string{
		"saas.enabled": "true",
		"oidc.enabled": "false",
	}

	for name, deployment := range map[string]appsv1.Deployment{
		"starter": renderStarterDeployment(t, values),
		"worker":  renderWorkerDeployment(t, values),
	} {
		t.Run(name, func(t *testing.T) {
			env := containerEnv(deployment)
			assert.Empty(t, env["CAMUNDA_CLIENT_MODE"])
			assert.NotEmpty(t, env["CAMUNDA_CLIENT_GRPC_ADDRESS"])
			assert.False(t, secretNames(deployment)["load-test-credentials"])
			assert.Empty(t, deployment.Spec.Template.Annotations["checksum/saas-credentials"])
		})
	}

	_, err := renderCredentialsSecretRaw(t, values)
	assert.Error(t, err, "credentials.yaml should render empty when oidc.enabled=false overrides saas.enabled=true")
}

func TestOidcDefaultCredentialValuesAndPartialMerge(t *testing.T) {
	secret := renderCredentialsSecret(t, map[string]string{
		"oidc.enabled":              "true",
		"saas.credentials.clientId": "from-saas",
		"oidc.credentials.authType": "CUSTOM",
	})

	assert.Equal(t, "from-saas", secret.StringData["clientId"])
	assert.Equal(t, "CUSTOM", secret.StringData["authType"])
	assert.Equal(t, "https://login.cloud.ultrawombat.com/oauth/token", secret.StringData["authServer"])
	assert.Equal(t, "zeebe.ultrawombat.com", secret.StringData["authorizationAudience"])
	assert.Empty(t, secret.StringData["clientSecret"])
}

func TestNumericSetCredentialsDoNotTurnIntoScientificNotation(t *testing.T) {
	raw, err := renderCredentialsSecretRaw(t, map[string]string{
		"oidc.enabled":              "true",
		"oidc.credentials.clientId": "12345678",
	})
	require.NoError(t, err)
	assert.Contains(t, raw, "clientId: 12345678")
	assert.NotContains(t, raw, "1.2345678e+07")
}

func TestNullYamlOverridesFallThroughToDefaultsAndSaas(t *testing.T) {
	t.Run("null credentials preserves defaults", func(t *testing.T) {
		file := writeValuesFile(t, "oidc:\n  enabled: true\n  credentials: null\n")
		secret := renderCredentialsSecret(t, nil, file)

		assert.Equal(t, "https://login.cloud.ultrawombat.com/oauth/token", secret.StringData["authServer"])
		assert.Equal(t, "OAUTH", secret.StringData["authType"])
		assert.Equal(t, "zeebe.ultrawombat.com", secret.StringData["authorizationAudience"])
	})

	t.Run("null oidc.enabled falls back to saas.enabled", func(t *testing.T) {
		file := writeValuesFile(t, "saas:\n  enabled: true\n  credentials:\n    existingSecret: legacy-secret\noidc:\n  enabled: null\n")
		starter := renderStarterDeployment(t, nil, file)

		env := containerEnv(starter)
		assert.Equal(t, "saas", env["CAMUNDA_CLIENT_MODE"])
		assert.Equal(t, "oidc", env["CAMUNDA_CLIENT_AUTH_METHOD"])
		assert.True(t, secretNames(starter)["legacy-secret"])
	})
}
