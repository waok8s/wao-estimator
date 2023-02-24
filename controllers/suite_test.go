//go:build testOnExistingCluster

package controllers_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	v1beta1 "github.com/Nedopro2022/wao-estimator/api/v1beta1"
	"github.com/Nedopro2022/wao-estimator/controllers"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
		UseExistingCluster:    pointer.Bool(true),
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = v1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var (
	wait    = func() { time.Sleep(100 * time.Millisecond) }
	testNS  = "default"
	testEC1 = v1beta1.Estimator{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: testNS,
			Name:      "hoge",
		},
		Spec: v1beta1.EstimatorSpec{
			DefaultNodeConfig: &v1beta1.NodeConfig{
				NodeMonitor: &v1beta1.NodeMonitor{
					RefreshInterval: &metav1.Duration{
						Duration: v1beta1.DefaultNodeMonitorRefreshInterval,
					},
					Agents: []v1beta1.NodeMonitorAgent{},
				},
				PowerConsumptionPredictor: &v1beta1.PowerConsumptionPredictor{
					Type:     v1beta1.PowerConsumptionPredictorTypeNone,
					Endpoint: "",
				},
			},
		},
	}
	testEC2 = v1beta1.Estimator{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: testNS,
			Name:      "fuga",
		},
		Spec: testEC1.Spec,
	}

	testLabelNodeMonitor = "waofed.bitmedia.co.jp/node-monitor"
	testLabelPCPredictor = "waofed.bitmedia.co.jp/power-consumption-predictor"

	testEC3 = v1beta1.Estimator{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: testNS,
			Name:      "default",
		},
		Spec: v1beta1.EstimatorSpec{
			DefaultNodeConfig: &v1beta1.NodeConfig{
				NodeMonitor: &v1beta1.NodeMonitor{
					RefreshInterval: &metav1.Duration{
						Duration: v1beta1.DefaultNodeMonitorRefreshInterval,
					},
					Agents: []v1beta1.NodeMonitorAgent{},
				},
				PowerConsumptionPredictor: &v1beta1.PowerConsumptionPredictor{
					Type:     v1beta1.PowerConsumptionPredictorTypeNone,
					Endpoint: "",
				},
			},
			NodeConfigOverrides: map[string]*v1beta1.NodeConfig{
				"wao-estimator-test-worker": &v1beta1.NodeConfig{
					NodeMonitor: &v1beta1.NodeMonitor{
						Agents: []v1beta1.NodeMonitorAgent{{Type: v1beta1.NodeMonitorTypeFake}},
					},
					PowerConsumptionPredictor: &v1beta1.PowerConsumptionPredictor{
						Type: v1beta1.PowerConsumptionPredictorTypeFake,
					},
				},
				"wao-estimator-test-worker2": &v1beta1.NodeConfig{
					NodeMonitor: &v1beta1.NodeMonitor{
						Agents: []v1beta1.NodeMonitorAgent{{Type: v1beta1.NodeMonitorTypeFake}},
					},
					PowerConsumptionPredictor: &v1beta1.PowerConsumptionPredictor{
						Type: v1beta1.PowerConsumptionPredictorTypeFake,
					},
				},
			},
		},
	}
)

var _ = Describe("Estimator controller", func() {
	var cncl context.CancelFunc
	var estimatorReconciler *controllers.EstimatorReconciler

	BeforeEach(func() {
		ctx, cancel := context.WithCancel(context.Background())
		cncl = cancel

		var err error
		err = k8sClient.DeleteAllOf(ctx, &v1beta1.Estimator{}, client.InNamespace(testNS))
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			var objs v1beta1.EstimatorList
			err = k8sClient.List(ctx, &objs, client.InNamespace(testNS))
			Expect(err).NotTo(HaveOccurred())
			return len(objs.Items)
		}).Should(Equal(0))

		mgr, err := ctrl.NewManager(cfg, ctrl.Options{
			Scheme: scheme.Scheme,
		})
		Expect(err).NotTo(HaveOccurred())

		estimatorReconciler = &controllers.EstimatorReconciler{
			Client: k8sClient,
			Scheme: scheme.Scheme,
		}
		err = estimatorReconciler.SetupWithManager(mgr)
		Expect(err).NotTo(HaveOccurred())

		go func() {
			err := mgr.Start(ctx)
			if err != nil {
				panic(err)
			}
		}()
		wait()
	})

	AfterEach(func() {
		cncl() // stop the mgr
		wait()
	})

	It("should add/delete estimator.Estimator", func() {
		ctx := context.Background()

		// Estimators: empty
		Expect(estimatorReconciler.GetEstimators().Len()).To(Equal(0))

		// Estimators: default/hoge
		ec1 := testEC1
		err := k8sClient.Create(ctx, &ec1)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			return estimatorReconciler.GetEstimators().Len()
		}).Should(Equal(1))

		// Estimators: default/hoge, default/huga
		ec2 := testEC2
		err = k8sClient.Create(ctx, &ec2)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			return estimatorReconciler.GetEstimators().Len()
		}).Should(Equal(2))

		// Estimators: default/fuga
		err = k8sClient.Delete(ctx, &ec1)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			return estimatorReconciler.GetEstimators().Len()
		}).Should(Equal(1))

		// Estimators: empty
		err = k8sClient.Delete(ctx, &ec2)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			return estimatorReconciler.GetEstimators().Len()
		}).Should(Equal(0))
	})

	It("should update estimator.Estimator", func() {
		ctx := context.Background()

		// Estimators: empty
		Expect(estimatorReconciler.GetEstimators().Len()).To(Equal(0))

		// Estimators: default/hoge
		ec1 := testEC1
		err := k8sClient.Create(ctx, &ec1)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			return estimatorReconciler.GetEstimators().Len()
		}).Should(Equal(1))

		// Estimators: default/hoge(modified#1)
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKeyFromObject(&ec1), &ec1)
		}).Should(Succeed())
		Expect(ec1.Spec.DefaultNodeConfig.NodeMonitor.RefreshInterval.Duration).To(Equal(v1beta1.DefaultNodeMonitorRefreshInterval))
		op, err := ctrl.CreateOrUpdate(ctx, k8sClient, &ec1, func() error {
			ec1.Spec.DefaultNodeConfig.NodeMonitor.RefreshInterval = &metav1.Duration{Duration: time.Second}
			return nil
		})
		Expect(op).To(Equal(controllerutil.OperationResultUpdated))
		Expect(err).NotTo(HaveOccurred())

		// Estimators: default/hoge(modified#2)
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKeyFromObject(&ec1), &ec1)
		}).Should(Succeed())
		Expect(ec1.Spec.DefaultNodeConfig.NodeMonitor.RefreshInterval.Duration).To(Equal(time.Second))
		op, err = ctrl.CreateOrUpdate(ctx, k8sClient, &ec1, func() error {
			ec1.Spec.DefaultNodeConfig.NodeMonitor.RefreshInterval = &metav1.Duration{Duration: v1beta1.DefaultNodeMonitorRefreshInterval}
			return nil
		})
		Expect(op).To(Equal(controllerutil.OperationResultUpdated))
		Expect(err).NotTo(HaveOccurred())

		// Estimators: empty
		err = k8sClient.Delete(ctx, &ec1)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			return estimatorReconciler.GetEstimators().Len()
		}).Should(Equal(1))
	})

	It("should estimate", func() {
		ctx := context.Background()

		// Estimators: empty
		Expect(estimatorReconciler.GetEstimators().Len()).To(Equal(0))

		// Estimators: default/default
		ec := testEC3
		err := k8sClient.Create(ctx, &ec)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			return estimatorReconciler.GetEstimators().Len()
		}).Should(Equal(1))

		estm, ok := estimatorReconciler.GetEstimators().Get(client.ObjectKeyFromObject(&ec).String())
		Expect(ok).To(BeTrue())

		watts, err := estm.EstimatePowerConsumption(ctx, 500, 5)
		Expect(err).To(BeNil())
		// wattMatrix=[
		//	            [inf inf inf inf inf inf]  (control-plane)
		//              [100 105 110 115 120 125]  (worker,  base=100, 10W/Core)
		//              [ 50  60  70  80  90 100]  (worker2, base=50,  20W/Core)
		//            ]
		// wattDiffs= [
		//              [inf inf inf inf inf]
		//              [  5  10  15  20  25]
		//              [ 10  20  30  40  50]
		//            ]
		// watts=       [  5  10  15  20  25]
		Expect(watts).To(Equal([]float64{5, 10, 15, 20, 25}))
	})

})
