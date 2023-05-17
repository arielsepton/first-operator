package operations

import (
	"github.com/arielsepton/first-operator/utils/helm"
	"helm.sh/helm/v3/pkg/action"
)

type UninstallParams struct {
	Namespace string
	Name      string
}

type UninstallOperation struct {
	Params UninstallParams
}

// TODO: This function doesnt use chart, fix it
func (uo UninstallOperation) Execute(uninstallParams interface{}) error {
	params := uninstallParams.(PossibleParams)

	actionConfig, err := helm.GetActionConfig(params.Namespace, params.Name)
	if err != nil {
		return err
	}

	client := action.NewUninstall(actionConfig)
	_, err = client.Run(params.Name)

	return err
}
