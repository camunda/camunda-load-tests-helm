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

func TestGoldenDefaults(t *testing.T) {

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"clients-service", "starter"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: name,
			Templates:      []string{"templates/" + name + ".yaml"},
			SetValues:      map[string]string{},
		})
	}
}

func TestGoldenWorkers(t *testing.T) {

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"workers"}

	values := map[string]string{
		"workers.otherWorker.jobType":  "otherWorker-job",
		"workers.otherWorker.replicas": "1",
		"workers.otherWorker.capacity": "10",
	}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: name,
			Templates:      []string{"templates/" + name + ".yaml"},
			SetValues:      values,
		})
	}
}

func TestGoldenExtendedStarter(t *testing.T) {
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"starter"}

	values := map[string]string{
		"starter.logLevel":          "INFO",
		"starter.payloadPath":       "empty.json",
		"starter.processId":         "real",
		"starter.bpmnXmlPath":       "bpmn/real.bpmn",
		"starter.extraResources[0]": "bpmn/extra.bpmn",
		"starter.extraResources[1]": "bpmn/extra.dmn",
		"starter.businessKey":       "customerId",
	}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "load-test",
			Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: name + "-extended",
			Templates:      []string{"templates/" + name + ".yaml"},
			SetValues:      values,
		})
	}
}

func TestGoldenWorkerWithMessage(t *testing.T) {
	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	suite.Run(t, &golden.TemplateGoldenTest{
		ChartPath:      chartPath,
		Release:        "load-test",
		Namespace:      "load-test-" + strings.ToLower(random.UniqueId()),
		GoldenFileName: "workers-with-message",
		Templates:      []string{"templates/workers.yaml"},
		SetValues: map[string]string{
			"workers.worker.replicas":                    "1",
			"workers.worker.capacity":                    "30",
			"workers.worker.jobType":                     "benchmark-task",
			"workers.worker.completionDelay":             "300ms",
			"workers.worker.message.name":                "my-message",
			"workers.worker.message.correlationVariable": "orderId",
		},
	})
}
