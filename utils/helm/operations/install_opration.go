package operations

import (
	"github.com/arielsepton/first-operator/utils/helm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
)

type InstallParams struct {
	Namespace string
	Name      string
	Chart     *chart.Chart
}

type InstallOperation struct {
}

func (io InstallOperation) Execute(installParams interface{}) error {
	params := installParams.(PossibleParams)

	actionConfig, err := helm.GetActionConfig(params.Namespace, params.Namespace)
	if err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = params.Namespace
	client.ReleaseName = params.Name

	values, err := helm.YamlToMap("values.yaml")
	if err != nil {
		return err
	}

	_, err = client.Run(params.Chart, values)
	return err
}
