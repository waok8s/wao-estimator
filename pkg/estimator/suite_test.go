package estimator_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Estimator Suite")
}

var _ = BeforeSuite(func() {
})

var _ = AfterSuite(func() {
})

var (
	wait = func() { time.Sleep(100 * time.Millisecond) }
)

var _ = Describe("Server/Client", func() {

	var es *estimator.Estimators
	var sv *estimator.Server
	var hsv *http.Server
	addr := net.JoinHostPort("localhost", estimator.ServerDefaultPort)
	httpAddr := "http://" + addr

	AfterEach(func() {
		err := hsv.Shutdown(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		wait()
		es = nil
		sv = nil
		hsv = nil
	})

	It("should fail to access with invalid URL", func() {
		es = &estimator.Estimators{}
		sv = estimator.NewServer(es)
		h, err := sv.Handler()
		Expect(err).NotTo(HaveOccurred())
		hsv = &http.Server{Addr: addr, Handler: h}
		go func() {
			hsv.ListenAndServe()
		}()
		wait()

		testAccess(addr, "", "default", "default", false)
	})

	It("shoud add/delete; no authentication", func() {
		es = &estimator.Estimators{}
		sv = estimator.NewServer(es)
		h, err := sv.Handler()
		Expect(err).NotTo(HaveOccurred())
		hsv = &http.Server{Addr: addr, Handler: h}
		go func() {
			hsv.ListenAndServe()
		}()
		wait()

		// empty
		testAccess(httpAddr, "", "default", "default", false)
		testAccess(httpAddr, "", "hoge", "fuga", false)

		// default/default
		ok := es.Add(client.ObjectKey{Namespace: "default", Name: "default"}.String(), estimator.NewEstimator(&estimator.Nodes{}))
		Expect(ok).To(BeTrue())
		testAccess(httpAddr, "", "default", "default", true)
		testAccess(httpAddr, "", "hoge", "fuga", false)

		// default/default, hoge/fuga
		ok = es.Add(client.ObjectKey{Namespace: "hoge", Name: "fuga"}.String(), estimator.NewEstimator(&estimator.Nodes{}))
		Expect(ok).To(BeTrue())
		testAccess(httpAddr, "", "default", "default", true)
		testAccess(httpAddr, "", "hoge", "fuga", true)

		// default/default, hoge/fuga
		es.Delete(client.ObjectKey{Namespace: "foo", Name: "bar"}.String())
		testAccess(httpAddr, "", "default", "default", true)
		testAccess(httpAddr, "", "hoge", "fuga", true)

		// default/default
		es.Delete(client.ObjectKey{Namespace: "hoge", Name: "fuga"}.String())
		testAccess(httpAddr, "", "default", "default", true)
		testAccess(httpAddr, "", "hoge", "fuga", false)

		// empty
		es.Delete(client.ObjectKey{Namespace: "default", Name: "default"}.String())
		testAccess(httpAddr, "", "default", "default", false)
		testAccess(httpAddr, "", "hoge", "fuga", false)
	})

	It("should access; with X-API-KEY authentication", func() {

		key1 := "foobar"
		key2 := "hogefuga"

		es = &estimator.Estimators{}
		sv = estimator.NewServer(es)
		h, err := sv.HandlerWithAuthFn(estimator.AuthFnAPIKey(map[string]struct{}{key1: {}, key2: {}}))
		Expect(err).NotTo(HaveOccurred())
		hsv = &http.Server{Addr: addr, Handler: h}
		go func() {
			hsv.ListenAndServe()
		}()
		wait()

		// Estimators: empty
		testAccess(httpAddr, "", "default", "default", false)
		testAccess(httpAddr, "xxx", "default", "default", false)
		testAccess(httpAddr, key1, "default", "default", false)
		testAccess(httpAddr, key2, "default", "default", false)

		// Estimators: default/default
		ok := es.Add(client.ObjectKey{Namespace: "default", Name: "default"}.String(), estimator.NewEstimator(&estimator.Nodes{}))
		Expect(ok).To(BeTrue())
		testAccess(httpAddr, "", "default", "default", false)
		testAccess(httpAddr, "xxx", "default", "default", false)
		testAccess(httpAddr, key1, "default", "default", true)
		testAccess(httpAddr, key2, "default", "default", true)

		// Estimators: empty
		es.Delete(client.ObjectKey{Namespace: "default", Name: "default"}.String())
		testAccess(httpAddr, "", "default", "default", false)
		testAccess(httpAddr, "xxx", "default", "default", false)
		testAccess(httpAddr, key1, "default", "default", false)
		testAccess(httpAddr, key2, "default", "default", false)
	})
})

func testAccess(httpAddr, apiKey, ns, name string, want bool) {
	opts := []estimator.ClientOption{estimator.ClientOptionGetRequestAsCurl(GinkgoWriter)}
	if apiKey != "" {
		opts = append(opts, estimator.ClientOptionAddRequestHeader(estimator.AuthFnAPIKeyRequestHeader, apiKey))
	}
	cl, err := estimator.NewClient(httpAddr, ns, name, opts...)
	Expect(err).NotTo(HaveOccurred())

	testFn := func() error {
		pc, err := cl.EstimatePowerConsumption(context.TODO(), 500, 5)
		if err != nil {
			return fmt.Errorf("err != nil : %w", err)
		}
		if pc == nil {
			return errors.New("pc == nil")
		}
		return nil
	}

	if want {
		Eventually(testFn).Should(Succeed())
	} else {
		Eventually(testFn).ShouldNot(Succeed())
	}
}

var _ = Describe("Node/Nodes", func() {

	var nm = &estimator.FakeNodeMonitor{
		GetFunc: func(ctx context.Context) (*estimator.NodeStatus, error) {
			time.Sleep(50 * time.Millisecond)
			return &estimator.NodeStatus{
				Timestamp:      time.Now(),
				CPUSockets:     2,
				CPUCores:       4,
				CPUUsages:      [][]float64{{10.0, 10.0, 10.0, 10.0}, {10.0, 10.0, 10.0, 10.0}},
				CPUTemps:       [][]float64{{30.0, 30.0, 30.0, 30.0}, {30.0, 30.0, 30.0, 30.0}},
				AmbientSensors: 2,
				AmbientTemps:   []float64{20.0, 20.0},
			}, nil
		},
	}
	var pcp = &estimator.FakePCPredictor{
		PredictFunc: estimator.PredictPCFnDummy,
	}

	n0 := estimator.NewNode("n0", nm, 300*time.Millisecond, pcp)

	status := n0.GetStatus()
	Expect(status).To(BeNil())

	// nodes.Add() calls Node.start()
	nodes := &estimator.Nodes{}
	ok := nodes.Add("n0", n0)
	Expect(ok).To(BeTrue())

	Eventually(func() *estimator.NodeStatus {
		return n0.GetStatus()
	}).ShouldNot(BeNil())
})
