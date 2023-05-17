package operations

import (
	"helm.sh/helm/v3/pkg/chart"
)

type PossibleParams struct {
	Namespace string
	Name      string
	Chart     *chart.Chart
}

type Operation interface {
	Execute(interface{}) error
}
