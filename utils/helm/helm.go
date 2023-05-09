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

func InstallChart(releaseNamespace string, releaseName string, chart *chart.Chart) error {
	actionConfig, err := GetActionConfig(releaseNamespace, releaseName)
	if err != nil {
		return err
	}

	Client := action.NewInstall(actionConfig)
	Client.Namespace = releaseNamespace
	Client.ReleaseName = releaseName

	// rel, err := Client.Run(chart, nil)
	// if err != nil {
	// 	return err
	// }

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

func UninstallChart(releaseNamespace string, releaseName string) error {
	actionConfig, err := GetActionConfig(releaseNamespace, releaseName)
	if err != nil {
		return err
	}

	Client := action.NewUninstall(actionConfig)
	rel, err := Client.Run(releaseName)
	if err != nil {
		return err
	}

	log.Log.Info(fmt.Sprintln("Successfully uninstalled release: ", rel.Release.Name))

	return nil
}

func UpgradeChart(releaseNamespace string, releaseName string, chart *chart.Chart) error {
	actionConfig, err := GetActionConfig(releaseNamespace, releaseName)
	if err != nil {
		return err
	}

	values, err := yamlToMap("values.yaml")
	if err != nil {
		return err
	}

	Client := action.NewUpgrade(actionConfig)
	rel, err := Client.Run(releaseName, chart, values)
	if err != nil {
		return err
	}

	log.Log.Info(fmt.Sprintln("Successfully upgraded release: ", rel.Name))

	return nil
}

func RunOperation(operation string, releaseNamespace string, releaseName string, chart *chart.Chart) error {
	if operation == "install" {
		InstallChart(releaseNamespace, releaseName, chart)
	}

	if operation == "uninstall" {
		UninstallChart(releaseNamespace, releaseName)
	}

	if operation == "upgrade" {
		UpgradeChart(releaseNamespace, releaseName, chart)
	}

	if operation == "done" {
		log.Log.Info(fmt.Sprintln("Successfully done!"))
	}

	return nil
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

// func updateCR() error {
// 	log.Log.Info(fmt.Sprintln("In update CR"))

// 	cr := &unstructured.Unstructured{}
// 	cr.SetGroupVersionKind(schema.GroupVersionKind{
// 		Group: "my.domain",
// 		Version: "v1alpha1",
// 		Kind:    "Helmer",
// 	})

// 	err := client.Client.Get(context.TODO(), types.NamespacedName{Name: "helmer-sample", Namespace: "default"}, cr)
// 	if err != nil {
// 		return err
// 	}

// 	spec := cr.Spec
// 	log.Log.Info(fmt.Sprintln("spec: ", spec))

// 	spec.operation = "done"
// 	log.Log.Info(fmt.Sprintln("spec: ", spec))

// 	cr.Spec = spec

// 	err = client.Client.Update(context.TODO(), cr)
// 	if err != nil {
// 		return err
// 	}

// 	log.Log.Info(fmt.Sprintln("CR updated successfully"))
// 	return nil
// }
