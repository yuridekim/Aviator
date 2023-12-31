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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ncputil "github.com/cloud-club/Aviator-service/pkg"

	vmv1 "vm.cloudclub.io/api/v1"
)

// ProvisionReconciler reconciles a Provision object
type ProvisionReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	ncpService *ncputil.NcpService
}

func NewProvisionReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	ncpService *ncputil.NcpService,
) *ProvisionReconciler {
	return &ProvisionReconciler{
		Client:     client,
		Scheme:     scheme,
		ncpService: ncpService,
	}
}

//+kubebuilder:rbac:groups=vm.cloudclub.io,resources=provisions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vm.cloudclub.io,resources=provisions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vm.cloudclub.io,resources=provisions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Provision object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ProvisionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.V(0).Info("Reconciling Provision request", "Request", req)

	original := &vmv1.Provision{}
	err := r.Get(ctx, req.NamespacedName, original)
	if err != nil {
		if errors.IsNotFound(err) {
			log.V(0).Info("Provision resource not found. Ignoring reconciliation.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Provision resource")
		return ctrl.Result{}, err
	}

	switch original.Status.Phase {
	case "", vmv1.ProvisionPhaseCreate:
		log.V(0).Info("Creating a new VM")
		// create
		err := r.ncpService.Server.Create(apiUrlCreate)
		if err != nil {
			log.Error(err, "Failed to create VM")
			return ctrl.Result{}, err
		}
	case vmv1.ProvisionPhaseUpdate:
		// update
		log.V(0).Info("Updating an existing VM")
		err := r.ncpService.Server.Update(apiUrlUpdate)
		if err != nil {
			log.Error(err, "Failed to update VM")
			return ctrl.Result{}, err
		}
	case vmv1.ProvisionPhaseStop:
		// delete
		log.V(0).Info("Stopping an existing VM")
		err := r.ncpService.Server.Stop(apiUrlDelete)
		if err != nil {
			log.Error(err, "Failed to stop VM")
			return ctrl.Result{}, err
		}
	case vmv1.ProvisionPhaseDelete:
		// delete
		log.V(0).Info("Deleting an existing VM")
		err := r.ncpService.Server.Delete(apiUrlDelete)
		if err != nil {
			log.Error(err, "Failed to delete VM")
			return ctrl.Result{}, err
		}
	case vmv1.ProvisionPhaseGet:
		// get info
		log.V(0).Info("Getting information for an existing VM")
		err := r.ncpService.Server.Get(apiUrlGet)
		if err != nil {
			log.Error(err, "Failed to get VM information")
			return ctrl.Result{}, err
		}
	default:
		log.V(0).Info("No action defined for the current phase")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller  with the Manager.
func (r *ProvisionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vmv1.Provision{}).
		Complete(r)
}
