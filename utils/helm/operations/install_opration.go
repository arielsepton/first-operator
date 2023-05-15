package operations

import (
	"github.com/arielsepton/first-operator/utils/helm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
)

type InstallParams struct {
	BasicParams
	Chart *chart.Chart
	// Additional parameters specific to install operation
}

type InstallOperation struct {
	Params InstallParams
}

func (io InstallOperation) Execute(releaseNamespace string, releaseName string, chart *chart.Chart) error {
	actionConfig, err := helm.GetActionConfig(releaseNamespace, releaseName)
	if err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = releaseNamespace
	client.ReleaseName = releaseName

	values, err := helm.YamlToMap("values.yaml")
	if err != nil {
		return err
	}

	_, err = client.Run(chart, values)
	return err
}
