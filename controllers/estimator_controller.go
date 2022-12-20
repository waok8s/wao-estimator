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

	sv := &estimator.Server{Estimators: &estimator.Estimators{}}

	r.estimators = sv.Estimators
	h, err := sv.Handler(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer, middleware.Heartbeat("/healthz"))
	if err != nil {
		return err
	}

	go http.ListenAndServe(net.JoinHostPort("", fmt.Sprint(estimator.ServerDefaultPort)), h)

	return nil
}

//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=waofed.bitmedia.co.jp,resources=estimators/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch

// Reconcile moves the current state of the cluster closer to the desired state.
func (r *EstimatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lg := log.FromContext(ctx)
	lg.Info("Reconcile")

	// get Estimator
	var estConf v1beta1.Estimator
	err := r.Get(ctx, req.NamespacedName, &estConf)
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
	if !estConf.DeletionTimestamp.IsZero() {
		lg.Info("Estimator is being deleted")
		return ctrl.Result{}, nil
	}

	// setup estimator.Node
	estNodeList, err := r.reconcileEstimatorNodes(ctx, &estConf)
	if err != nil {
		return ctrl.Result{}, err
	}

	// update estimator.Estimator
	// TODO: check diff instead of replacing whole estimator.Estimator to reduce hardware resource usage
	estNodes := &estimator.Nodes{}
	for _, en := range estNodeList {
		if ok := estNodes.Add(en.Name, en); !ok {
			err := fmt.Errorf("r.estNodes.Add() returned false: %s", en.Name)
			lg.Error(err, "duplicate node name found")
		}
	}
	e := &estimator.Estimator{Nodes: estNodes}
	r.estimators.Delete(req.String())
	if ok := r.estimators.Add(req.String(), e); !ok {
		err := fmt.Errorf("r.estimators.Add() returned false: %s", req.String())
		lg.Error(err, "unable to add estimator.Estimator")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *EstimatorReconciler) reconcileEstimatorNodes(ctx context.Context, estConf *v1beta1.Estimator) ([]*estimator.Node, error) {
	lg := log.FromContext(ctx)
	lg.Info("reconcileEstimatorNodes")

	var nodeList corev1.NodeList
	if err := r.List(ctx, &nodeList); err != nil {
		return nil, err
	}

	var estNodeList []*estimator.Node

	for _, node := range nodeList.Items {
		name := node.Name

		// spec.nodeMonitor
		var nm estimator.NodeMonitor
		nmType := v1beta1.NodeMonitorType(getFieldValue(estConf.Spec.NodeMonitor.Type, &node))
		switch nmType {
		case v1beta1.NodeMonitorTypeNone:
		case v1beta1.NodeMonitorTypeFake:
			nm = setupFakeNodeMonitor(r.Client, client.ObjectKeyFromObject(&node))
		case v1beta1.NodeMonitorTypeIPMIExporter:
			lg.Info(fmt.Sprintf("NodeMonitorType=%v is not implemented", nmType))
		case v1beta1.NodeMonitorTypeRedfish:
			lg.Info(fmt.Sprintf("NodeMonitorType=%v is not implemented", nmType))
		default:
			lg.Info(fmt.Sprintf("NodeMonitorType=%v is not defined", nmType))
		}
		lg.Info(fmt.Sprintf("spec.nodeMonitor.Type=%v nm=%+v", nmType, nm))

		// spec.powerConsumptionPredictor
		var pcp estimator.PowerConsumptionPredictor
		pcpType := v1beta1.PowerConsumptionPredictorType(getFieldValue(estConf.Spec.PowerConsumptionPredictor.Type, &node))
		switch pcpType {
		case v1beta1.PowerConsumptionPredictorTypeNone:
		case v1beta1.PowerConsumptionPredictorTypeFake:
			pcp = setupFakePCPredictor(r.Client, client.ObjectKeyFromObject(&node))
		case v1beta1.PowerConsumptionPredictorTypeMLServer:
			lg.Info(fmt.Sprintf("PowerConsumptionPredictorType=%v is not implemented", pcpType))
		case v1beta1.PowerConsumptionPredictorTypeTFServing:
			lg.Info(fmt.Sprintf("PowerConsumptionPredictorType=%v is not implemented", pcpType))
		default:
			lg.Info(fmt.Sprintf("PowerConsumptionPredictorType=%v is not defined", pcpType))
		}
		lg.Info(fmt.Sprintf("spec.powerConsumptionPredictor.Type=%v pcp=%+v", pcpType, pcp))

		estNode := estimator.NewNode(name, nm, estConf.Spec.NodeMonitor.RefreshInterval.Duration, pcp)
		estNodeList = append(estNodeList, estNode)
	}

	return estNodeList, nil
}

func getFieldValue(f v1beta1.Field, node *corev1.Node) string {
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
