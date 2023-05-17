package executer

import (
	"errors"

	"github.com/arielsepton/first-operator/utils/helm/operations"
	"helm.sh/helm/v3/pkg/chart"
)

var operationMap = map[string]operations.Operation{
	"install":   operations.InstallOperation{},
	"uninstall": operations.UninstallOperation{},
	"upgrade":   operations.UpgradeOperation{},
}

func RunOperation(operation string, releaseNamespace string, releaseName string, chart *chart.Chart) error {
	params := operations.PossibleParams{
		Namespace: releaseNamespace,
		Name:      releaseName,
		Chart:     chart,
	}

	operationFunc, found := operationMap[operation]
	if !found {
		return errors.New("unsupported operation")
	}
	// operationFunc.SetParams()

	return operationFunc.Execute(params)
}
