package operations

import (
	"github.com/arielsepton/first-operator/utils/helm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
)

type UninstallParams struct {
	BasicParams
	// Additional parameters specific to uninstall operation
}

type UninstallOperation struct {
	Params UninstallParams
}

// TODO: This function doesnt use chart, fix it
func (uo UninstallOperation) Execute(releaseNamespace string, releaseName string, chart *chart.Chart) error {
	actionConfig, err := helm.GetActionConfig(releaseNamespace, releaseName)
	if err != nil {
		return err
	}

	client := action.NewUninstall(actionConfig)
	_, err = client.Run(releaseName)

	return err
}
