package helm

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func GetChart(path string) (*chart.Chart, error) {
	chart, err := loader.Load(path)
	if err != nil {
		return nil, err
	}

	return chart, nil
}

func GetActionConfig(releaseNamespace string, releaseName string) (*action.Configuration, error) {
	kubeconfigPath := "/root/.kube/config"

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(kube.GetConfig(kubeconfigPath, "", releaseNamespace), releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		log.Log.Info(fmt.Sprintf(format, v))
	}); err != nil {
		return nil, err
	}

	return actionConfig, nil
}

func YamlToMap(yamlFilePath string) (map[string]interface{}, error) {
	// Read the YAML file into a byte slice
	yamlData, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into a map[string]interface{}
	var m map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &m); err != nil {
		return nil, err
	}

	return m, nil
}
