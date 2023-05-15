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
	operationFunc, found := operationMap[operation]
	if !found {
		return errors.New("unsupported operation")
	}

	return operationFunc.Execute(releaseNamespace, releaseName, chart)
}
