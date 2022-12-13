package controllers

import (
	"context"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1beta1 "github.com/Nedopro2022/wao-estimator/api/v1beta1"
	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
)

// EstimatorReconciler reconciles a Estimator object
type EstimatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	estimators *estimator.Estimators
}

//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Estimator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *EstimatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EstimatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := r.startEstimatorServer(); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Estimator{}).
		Complete(r)
}

func (r *EstimatorReconciler) startEstimatorServer() error {
	addr := net.JoinHostPort("", estimator.ServerDefaultPort)

	r.estimators = &estimator.Estimators{}

	h, err := estimator.NewServer(r.estimators).Handler(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer, middleware.Heartbeat("/healthz"))
	if err != nil {
		return err
	}

	go http.ListenAndServe(addr, h)

	return nil
}
