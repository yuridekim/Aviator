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
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ncputil "github.com/cloud-club/Aviator-service/pkg"
	server "github.com/cloud-club/Aviator-service/types/server"

	vmv1 "vm.cloudclub.io/api/v1"
)

var provisionReconcileMap map[string]func(*ProvisionReconciler, logr.Logger, string, *vmv1.Provision, interface{}) error

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
	initProvisionReconcileMap()
	return &ProvisionReconciler{
		Client:     client,
		Scheme:     scheme,
		ncpService: ncpService,
	}
}

func initProvisionReconcileMap() {
	provisionReconcileMap = make(map[string]func(*ProvisionReconciler, logr.Logger, string, *vmv1.Provision, interface{}) error)
	provisionReconcileMap["provision"] = provision
	provisionReconcileMap["deProvision"] = deProvision
	provisionReconcileMap["update"] = update
	provisionReconcileMap["get"] = get
	provisionReconcileMap["stop"] = stop
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
	log.V(ErrorLevelIsInfo).Info("Reconciling Provision request", "Request", req)

	original := &vmv1.Provision{}
	err := r.Get(ctx, req.NamespacedName, original)
	if err != nil {
		if errors.IsNotFound(err) {
			log.V(ErrorLevelIsInfo).Info("Provision resource not found. Ignoring reconciliation.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Provision resource")
		return ctrl.Result{}, err
	}

	switch original.Spec.Phase {
	case "", vmv1.ProvisionPhaseCreate:
		if v, ok := provisionReconcileMap["provision"]; ok {
			if err = v(r, log, apiUrlCreate, original, nil); err != nil {
				log.Error(err, "Failed to create VM")
				return ctrl.Result{}, err
			}
		}
	case vmv1.ProvisionPhaseUpdate:
		if v, ok := provisionReconcileMap["update"]; ok {
			if err = v(r, log, apiUrlUpdate, original, nil); err != nil {
				log.Error(err, "Failed to update VM")
				return ctrl.Result{}, err
			}
		}
	case vmv1.ProvisionPhaseStop:
		if v, ok := provisionReconcileMap["stop"]; ok {
			if err = v(r, log, apiUrlStop, original, nil); err != nil {
				log.Error(err, "Failed to stop VM")
				return ctrl.Result{}, err
			}
		}
	case vmv1.ProvisionPhaseDelete:
		if v, ok := provisionReconcileMap["deProvision"]; ok {
			if err = v(r, log, apiUrlDelete, original, nil); err != nil {
				log.Error(err, "Failed to delete VM")
				return ctrl.Result{}, err
			}
		}
	case vmv1.ProvisionPhaseGet:
		if v, ok := provisionReconcileMap["get"]; ok {
			if err = v(r, log, apiUrlGet, original, nil); err != nil {
				log.Error(err, "Failed to get VM information")
				return ctrl.Result{}, err
			}
		}
	default:
		log.V(ErrorLevelIsAnError).Error(err, "No action defined for the current phase",
			"reconcile phase", original.Spec.Phase, "namespace", req.NamespacedName)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller  with the Manager.
func (r *ProvisionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vmv1.Provision{}).
		Complete(r)
}

func provision(r *ProvisionReconciler, log logr.Logger, url string, original *vmv1.Provision, payload interface{}) error {
	log.V(ErrorLevelIsInfo).Info("Creating a new VM")
	csr := &server.CreateServerRequest{
		ServerImageProductCode:    original.Spec.Server.ImageProductCode,
		VpcNo:                     original.Spec.VpcNo,
		SubnetNo:                  original.Spec.SubnetNo,
		NetworkInterfaceOrder:     original.Spec.NetworkInterface.Order,
		AccessControlGroupNoListN: original.Spec.AccessControlGroupNoListN,
		ServerProductCode:         original.Spec.Server.ProductCode,
	}
	createServerResponse, err := r.ncpService.Server.Create(ncputil.API_URL+ncputil.CREATE_SERVER_INSTANCE_PATH, csr, []int{1, 1})
	fmt.Println(createServerResponse)
	return err
}

func deProvision(r *ProvisionReconciler, log logr.Logger, url string, original *vmv1.Provision, payload interface{}) error {
	log.V(ErrorLevelIsInfo).Info("Deleting an existing VM")
	dsr := &server.DeleteServerRequest{ServerNo: original.Spec.ServerNo}
	deleteServerResponse, err := r.ncpService.Server.Delete(ncputil.API_URL+ncputil.DELETE_SERVER_INSTANCE_PATH, dsr)
	fmt.Println(deleteServerResponse)
	return err
}

func update(r *ProvisionReconciler, log logr.Logger, url string, original *vmv1.Provision, payload interface{}) error {
	log.V(ErrorLevelIsInfo).Info("Updating an existing VM")
	usr := &server.UpdateServerRequest{
		ServerInstanceNo:  original.Spec.ServerInstanceNo,
		ServerProductCode: original.Spec.Server.ProductCode,
	}
	updateServerResponse, err := r.ncpService.Server.Update(ncputil.API_URL+ncputil.UPDATE_SERVER_INSTANCE_PATH, usr)
	fmt.Println(updateServerResponse)
	return err
}

func stop(r *ProvisionReconciler, log logr.Logger, url string, original *vmv1.Provision, payload interface{}) error {
	log.V(ErrorLevelIsInfo).Info("Stopping an existing VM")
	ssr := &server.StopServerRequest{ServerNo: original.Spec.ServerNo}
	stopServerResponse, err := r.ncpService.Server.Stop(ncputil.API_URL+ncputil.STOP_SERVER_INSTANCE_PATH, ssr)
	fmt.Println(stopServerResponse)
	return err
}

func get(r *ProvisionReconciler, log logr.Logger, url string, original *vmv1.Provision, payload interface{}) error {
	log.V(ErrorLevelIsInfo).Info("Getting information for an existing VM")
	lsr := &server.ListServerRequest{RegionCode: original.Spec.RegionCode}
	serverListResponse, err := r.ncpService.Server.List(ncputil.API_URL+ncputil.GET_SERVER_INSTANCE_PATH, lsr)
	fmt.Println(serverListResponse)
	return err
}
