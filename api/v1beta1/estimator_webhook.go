package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var estimatorlog = logf.Log.WithName("estimator-resource")

func (r *Estimator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-waofed-bitmedia-co-jp-v1beta1-estimator,mutating=true,failurePolicy=fail,sideEffects=None,groups=waofed.bitmedia.co.jp,resources=estimators,verbs=create;update,versions=v1beta1,name=mestimator.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Estimator{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Estimator) Default() {
	estimatorlog.Info("default", "name", r.Name)
	r.defaultDefaultNodeConfig()
	r.defaultNodeConfigOverrides()
}

func (r *Estimator) defaultDefaultNodeConfig() {
	if r.Spec.DefaultNodeConfig == nil {
		r.Spec.DefaultNodeConfig = &NodeConfig{}
	}
	// NodeMonitor
	if r.Spec.DefaultNodeConfig.NodeMonitor == nil {
		r.Spec.DefaultNodeConfig.NodeMonitor = &NodeMonitor{}
	}
	if r.Spec.DefaultNodeConfig.NodeMonitor.RefreshInterval == nil {
		r.Spec.DefaultNodeConfig.NodeMonitor.RefreshInterval = &metav1.Duration{Duration: DefaultNodeMonitorRefreshInterval}
	}
	if r.Spec.DefaultNodeConfig.NodeMonitor.Agents == nil {
		r.Spec.DefaultNodeConfig.NodeMonitor.Agents = []NodeMonitorAgent{}
	}
	// PowerConsumptionPredictor
	if r.Spec.DefaultNodeConfig.PowerConsumptionPredictor == nil {
		r.Spec.DefaultNodeConfig.PowerConsumptionPredictor = &PowerConsumptionPredictor{}
	}
	if r.Spec.DefaultNodeConfig.PowerConsumptionPredictor.Type == "" {
		r.Spec.DefaultNodeConfig.PowerConsumptionPredictor.Type = PowerConsumptionPredictorTypeNone
	}
}

func (r *Estimator) defaultNodeConfigOverrides() {}

// NOTE: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-waofed-bitmedia-co-jp-v1beta1-estimator,mutating=false,failurePolicy=fail,sideEffects=None,groups=waofed.bitmedia.co.jp,resources=estimators,verbs=create;update,versions=v1beta1,name=vestimator.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Estimator{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Estimator) ValidateCreate() error {
	estimatorlog.Info("validate create", "name", r.Name)
	return r.validateSpec()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Estimator) ValidateUpdate(old runtime.Object) error {
	estimatorlog.Info("validate update", "name", r.Name)
	return r.validateSpec()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Estimator) ValidateDelete() error {
	estimatorlog.Info("validate delete", "name", r.Name)
	// NOTE: No validations needed upon deletion.
	return nil
}

func (r *Estimator) validateSpec() error {
	return nil
}
