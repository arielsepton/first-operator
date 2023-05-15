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

type Operation interface {
	execute(releaseNamespace string, releaseName string, chart *chart.Chart) error
}

type InstallOperation struct{}

func (i InstallOperation) execute(releaseNamespace string, releaseName string, chart *chart.Chart) error {
	actionConfig, err := GetActionConfig(releaseNamespace, releaseName)
	if err != nil {
		return err
	}

	Client := action.NewInstall(actionConfig)
	Client.Namespace = releaseNamespace
	Client.ReleaseName = releaseName

	values, err := yamlToMap("values.yaml")
	if err != nil {
		return err
	}

	rel, err := Client.Run(chart, values)
	if err != nil {
		log.Log.Info(fmt.Sprintln("err: ", err))
		return err
	}

	log.Log.Info(fmt.Sprintln("Successfully installed release: ", rel.Name))
	return nil
}

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

func yamlToMap(yamlFilePath string) (map[string]interface{}, error) {
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
