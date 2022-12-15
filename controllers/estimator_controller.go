package controllers

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators/finalizers,verbs=update

// Reconcile moves the current state of the cluster closer to the desired state.
func (r *EstimatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lg := log.FromContext(ctx)
	lg.Info("Reconcile")

	// get Estimator
	var estimatorConf v1beta1.Estimator
	err := r.Get(ctx, req.NamespacedName, &estimatorConf)
	if errors.IsNotFound(err) {
		lg.Info("Estimator is deleted")

		// delete estimator.Estimator
		r.estimators.Delete(req.String())

		return ctrl.Result{}, nil
	}
	if err != nil {
		lg.Error(err, "unable to get Estimator")
		return ctrl.Result{}, err
	}
	if !estimatorConf.DeletionTimestamp.IsZero() {
		lg.Info("Estimator is being deleted")
		return ctrl.Result{}, nil
	}

	// init estimator.Estimator
	// TODO: set Nodes
	e := estimator.NewEstimator(&estimator.Nodes{})
	if ok := r.estimators.Add(req.String(), e); !ok {
		err := fmt.Errorf("r.estimators.Add() returned false: %s", req.String())
		lg.Error(err, "unable to add estimator.Estimator")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func GetFieldValue(f v1beta1.Field, node *corev1.Node) string {
	switch {
	case f.Override != nil && f.Override.Label != nil && node != nil:
		v, ok := node.Labels[*f.Override.Label]
		if !ok {
			return f.Default
		}
		return v
	default:
		return f.Default
	}
}
