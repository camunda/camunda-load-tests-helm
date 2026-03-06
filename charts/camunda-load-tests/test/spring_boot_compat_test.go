// Copyright 2022 Camunda Services GmbH
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

func TestSpringBootEnvVarsStarter(t *testing.T) {
	// Verify that Spring Boot env vars are correctly set alongside old HOCON vars
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", "load-test-"+strings.ToLower(random.UniqueId())),
		SetValues: map[string]string{
			"starter.rate":      "200",
			"starter.processId": "my-process",
		},
	}

	output := helm.RenderTemplate(t, options, chartPath, "load-test", []string{"templates/starter.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	envVars := make(map[string]string)
	for _, env := range deployment.Spec.Template.Spec.Containers[0].Env {
		envVars[env.Name] = env.Value
	}

	// Spring Boot env vars
	assert.Equal(t, "200", envVars["LOAD_TESTER_STARTER_RATE"])
	assert.Equal(t, "my-process", envVars["LOAD_TESTER_STARTER_PROCESS_ID"])
	assert.Equal(t, "0", envVars["LOAD_TESTER_STARTER_DURATION_LIMIT"])
	assert.Equal(t, "false", envVars["CAMUNDA_CLIENT_PREFER_REST_OVER_GRPC"])
	assert.Equal(t, "http://camunda-gateway:26500", envVars["CAMUNDA_CLIENT_GRPC_ADDRESS"])
	assert.Equal(t, "http://camunda-gateway:8080", envVars["CAMUNDA_CLIENT_REST_ADDRESS"])

	// Old HOCON vars still present in JDK_JAVA_OPTIONS
	jdkOpts := envVars["JDK_JAVA_OPTIONS"]
	assert.Contains(t, jdkOpts, "-Dapp.starter.rate=200")
	assert.Contains(t, jdkOpts, "-Dapp.starter.processId=\"my-process\"")
}

func TestSpringBootEnvVarsWorker(t *testing.T) {
	// Verify that Spring Boot env vars are correctly set for workers
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", "load-test-"+strings.ToLower(random.UniqueId())),
		SetValues: map[string]string{
			"workers.myWorker.replicas":        "1",
			"workers.myWorker.capacity":        "50",
			"workers.myWorker.threads":         "5",
			"workers.myWorker.jobType":         "my-job",
			"workers.myWorker.completionDelay": "200ms",
		},
	}

	output := helm.RenderTemplate(t, options, chartPath, "load-test", []string{"templates/workers.yaml"})

	// Split output to find the myWorker deployment (there may be multiple)
	docs := strings.Split(output, "---")
	var myWorkerOutput string
	for _, doc := range docs {
		if strings.Contains(doc, "name: myWorker") {
			myWorkerOutput = doc
			break
		}
	}
	require.NotEmpty(t, myWorkerOutput, "myWorker deployment not found")

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, myWorkerOutput, &deployment)

	envVars := make(map[string]string)
	for _, env := range deployment.Spec.Template.Spec.Containers[0].Env {
		envVars[env.Name] = env.Value
	}

	// Spring Boot env vars
	assert.Equal(t, "50", envVars["LOAD_TESTER_WORKER_CAPACITY"])
	assert.Equal(t, "5", envVars["LOAD_TESTER_WORKER_THREADS"])
	assert.Equal(t, "my-job", envVars["LOAD_TESTER_WORKER_JOB_TYPE"])
	assert.Equal(t, "200ms", envVars["LOAD_TESTER_WORKER_COMPLETION_DELAY"])
	assert.Equal(t, "1ms", envVars["LOAD_TESTER_WORKER_POLLING_DELAY"])
	assert.Equal(t, "myWorker", envVars["LOAD_TESTER_WORKER_WORKER_NAME"])
	assert.Equal(t, "http://camunda-gateway:26500", envVars["CAMUNDA_CLIENT_GRPC_ADDRESS"])

	// Old HOCON vars still present
	jdkOpts := envVars["JDK_JAVA_OPTIONS"]
	assert.Contains(t, jdkOpts, "-Dapp.worker.capacity=50")
	assert.Contains(t, jdkOpts, "-Dapp.worker.jobType=\"my-job\"")
}

func TestSpringBootEnvVarsSaaSMode(t *testing.T) {
	// Verify SaaS mode sets CAMUNDA_CLIENT_MODE and CAMUNDA_CLIENT_AUTH_METHOD
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", "load-test-"+strings.ToLower(random.UniqueId())),
		SetValues: map[string]string{
			"saas.enabled":                    "true",
			"saas.credentials.existingSecret": "my-secret",
		},
	}

	// Test starter
	starterOutput := helm.RenderTemplate(t, options, chartPath, "load-test", []string{"templates/starter.yaml"})
	var starterDeployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, starterOutput, &starterDeployment)

	starterEnv := make(map[string]string)
	for _, env := range starterDeployment.Spec.Template.Spec.Containers[0].Env {
		starterEnv[env.Name] = env.Value
	}

	assert.Equal(t, "saas", starterEnv["CAMUNDA_CLIENT_MODE"])
	assert.Equal(t, "oidc", starterEnv["CAMUNDA_CLIENT_AUTH_METHOD"])
	assert.Empty(t, starterEnv["CAMUNDA_CLIENT_GRPC_ADDRESS"], "Should not set GRPC address in SaaS mode")
	assert.Empty(t, starterEnv["CAMUNDA_CLIENT_REST_ADDRESS"], "Should not set REST address in SaaS mode")

	// Test worker
	workerOutput := helm.RenderTemplate(t, options, chartPath, "load-test", []string{"templates/workers.yaml"})

	docs := strings.Split(workerOutput, "---")
	var workerDoc string
	for _, doc := range docs {
		if strings.Contains(doc, "kind: Deployment") {
			workerDoc = doc
			break
		}
	}
	require.NotEmpty(t, workerDoc)

	var workerDeployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, workerDoc, &workerDeployment)

	workerEnv := make(map[string]string)
	for _, env := range workerDeployment.Spec.Template.Spec.Containers[0].Env {
		workerEnv[env.Name] = env.Value
	}

	assert.Equal(t, "saas", workerEnv["CAMUNDA_CLIENT_MODE"])
	assert.Equal(t, "oidc", workerEnv["CAMUNDA_CLIENT_AUTH_METHOD"])
}
