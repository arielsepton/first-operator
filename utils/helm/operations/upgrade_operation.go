package operations

import (
	"github.com/arielsepton/first-operator/utils/helm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
)

type UpgradeParams struct {
	Namespace string
	Name      string
	Chart     *chart.Chart
}

type UpgradeOperation struct {
	Params UpgradeParams
}

func (uo UpgradeOperation) Execute(upgradeParams interface{}) error {
	params := upgradeParams.(PossibleParams)

	actionConfig, err := helm.GetActionConfig(params.Namespace, params.Name)
	if err != nil {
		return err
	}

	client := action.NewUpgrade(actionConfig)
	values, err := helm.YamlToMap("values.yaml")
	if err != nil {
		return err
	}

	_, err = client.Run(params.Name, params.Chart, values)
	return err
}
