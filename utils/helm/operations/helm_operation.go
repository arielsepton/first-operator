package operations

import (
	"helm.sh/helm/v3/pkg/chart"
)

type BasicParams struct {
	Namespace string
	Name      string
}

type Operation interface {
	Execute(releaseNamespace string, releaseName string, chart *chart.Chart) error
}
