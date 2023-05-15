package operations

import (
	"github.com/arielsepton/first-operator/utils/helm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
)

type UpgradeParams struct {
	BasicParams
	Chart *chart.Chart
	// Additional parameters specific to upgrade operation
}

type UpgradeOperation struct {
	Params UpgradeParams
}

func (uo UpgradeOperation) Execute(releaseNamespace string, releaseName string, chart *chart.Chart) error {
	actionConfig, err := helm.GetActionConfig(releaseNamespace, releaseName)
	if err != nil {
		return err
	}

	client := action.NewUpgrade(actionConfig)
	values, err := helm.YamlToMap("values.yaml")
	if err != nil {
		return err
	}

	_, err = client.Run(releaseName, chart, values)
	return err

}
