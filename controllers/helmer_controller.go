/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/arielsepton/first-operator/utils/helm"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/arielsepton/first-operator/api/v1alpha1"
)

// HelmerReconciler reconciles a Helmer object
type HelmerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.domain,resources=helmers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.domain,resources=helmers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.domain,resources=helmers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Helmer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *HelmerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	helmer := &apiv1alpha1.Helmer{}
	err := r.Get(ctx, req.NamespacedName, helmer)
	if err != nil {
		return ctrl.Result{}, nil
	}

	releaseNamespace := "default"
	chartPath := helmer.Spec.Chart
	releaseName := helmer.Spec.ReleaseName
	operation := helmer.Spec.Operation

	if helmer.Status.Status == "" {
		log.Log.Info(fmt.Sprintln("Hello this is the chart: ", chartPath))
		log.Log.Info(fmt.Sprintln("Hello this is the operation: ", operation))

		chart, err := helm.GetChart(chartPath)
		if err != nil {
			log.Log.Info(fmt.Sprintln("Hello this is the func: ", "GetChart"))
			updateCRStatusFAIL(r, ctx, req, err)
			return ctrl.Result{}, err
		}

		if err = helm.RunOperation1(operation, releaseNamespace, releaseName, chart); err != nil {
			log.Log.Info(fmt.Sprintln("Hello this is the func: ", "RunOperation"))
			updateCRStatusFAIL(r, ctx, req, err)
			return ctrl.Result{}, err
		}

		err = updateCRStatusSUCCESS(r, ctx, req)
		if err != nil {
			log.Log.Info(fmt.Sprintln("Hello this is the err: ", err))
		}
	}

	return ctrl.Result{}, nil
}

func updateCRStatusSUCCESS(r *HelmerReconciler, ctx context.Context, req ctrl.Request) error {
	helmer := &apiv1alpha1.Helmer{}
	err := r.Get(ctx, req.NamespacedName, helmer)
	if err != nil {
		return err
	}

	status := apiv1alpha1.HelmerStatus{}
	status.Status = apiv1alpha1.SUCCESS
	status.Error = ""

	helmer.Status = status
	err = r.Status().Update(ctx, helmer)
	if err != nil {
		log.Log.Info(fmt.Sprintln("Hello this is the error on update: ", err))

		return err
	}

	log.Log.Info(fmt.Sprintln("Hello this is the status: ", helmer.Status.Status))
	return nil
}

func updateCRStatusFAIL(r *HelmerReconciler, ctx context.Context, req ctrl.Request, err error) error {
	helmer := &apiv1alpha1.Helmer{}
	r.Get(ctx, req.NamespacedName, helmer)

	helmer.Status.Status = apiv1alpha1.FAILED
	helmer.Status.Error = err.Error()
	r.Status().Update(ctx, helmer)

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelmerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Helmer{}).
		Complete(r)
}
