package helm

import (
	"errors"
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type OperarionParams struct{}

type BasicParams struct {
	Namespace string
	Name      string
}

type InstallParams struct {
	BasicParams
	Chart *chart.Chart
	// Additional parameters specific to install operation
}

type UninstallParams struct {
	BasicParams
	// Additional parameters specific to uninstall operation
}

type UpgradeParams struct {
	BasicParams
	Chart *chart.Chart
	// Additional parameters specific to upgrade operation
}

type Operation1 interface {
	Execute(releaseNamespace string, releaseName string, chart *chart.Chart) error
}

type InstallOperation1 struct {
	Params InstallParams
}

func (io InstallOperation1) Execute(releaseNamespace string, releaseName string, chart *chart.Chart) error {
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
	// Access io.Params to use install-specific variables
}

type UninstallOperation struct {
	Params UninstallParams
}

func (uo UninstallOperation) Execute(releaseNamespace string, releaseName string, chart *chart.Chart) error {
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

type UpgradeOperation struct {
	Params UpgradeParams
}

func (uo UpgradeOperation) Execute(releaseNamespace string, releaseName string, chart *chart.Chart) error {
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

type OperationFunc func() error

var operationMap = map[string]Operation1{
	"install":   InstallOperation1{},
	"uninstall": UninstallOperation{},
	"upgrade":   UpgradeOperation{},
}

func RunOperation1(operation string, releaseNamespace string, releaseName string, chart *chart.Chart) error {
	operationFunc, found := operationMap[operation]
	if !found {
		return errors.New("unsupported operation")
	}

	return operationFunc.Execute(releaseNamespace, releaseName, chart)
}
