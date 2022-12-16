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

func init() {
	SetDefaultEventuallyTimeout(3 * time.Second)
	SetDefaultEventuallyPollingInterval(100 * time.Millisecond)
}

var (
	wait     = func() { time.Sleep(100 * time.Millisecond) }
	testHost = "localhost"
	testPort = fmt.Sprint(estimator.ServerDefaultPort + 1)
)

var _ = Describe("Server/Client", func() {

	var es *estimator.Estimators
	var sv *estimator.Server
	var hsv *http.Server
	addr := net.JoinHostPort(testHost, testPort)
	httpAddr := "http://" + addr

	AfterEach(func() {
		err := hsv.Shutdown(context.Background())
		Expect(err).NotTo(HaveOccurred())
		wait()
		es = nil
		sv = nil
		hsv = nil
	})

	It("should fail to access with invalid URL", func() {
		sv = &estimator.Server{}
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
		sv = &estimator.Server{Estimators: es}
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
		ok := es.Add(client.ObjectKey{Namespace: "default", Name: "default"}.String(), &estimator.Estimator{Nodes: nil})
		Expect(ok).To(BeTrue())
		testAccess(httpAddr, "", "default", "default", true)
		testAccess(httpAddr, "", "hoge", "fuga", false)

		// default/default, hoge/fuga
		ok = es.Add(client.ObjectKey{Namespace: "hoge", Name: "fuga"}.String(), &estimator.Estimator{Nodes: nil})
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
		sv = &estimator.Server{Estimators: es}
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
		ok := es.Add(client.ObjectKey{Namespace: "default", Name: "default"}.String(), &estimator.Estimator{Nodes: nil})
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
		pc, err := cl.EstimatePowerConsumption(context.Background(), 500, 5)
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

	ctx := context.Background()

	nm0 := &estimator.FakeNodeMonitor{
		FetchFunc: func(ctx context.Context) (estimator.NodeStatus, error) {
			time.Sleep(50 * time.Millisecond)
			return estimator.NodeStatus{
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
	nm1 := &estimator.FakeNodeMonitor{
		FetchFunc: func(ctx context.Context) (estimator.NodeStatus, error) {
			time.Sleep(50 * time.Millisecond)
			return estimator.NodeStatus{}, estimator.ErrNodeMonitor
		},
	}
	pcp0 := &estimator.FakePCPredictor{
		PredictFunc: estimator.PredictPCFnDummy,
	}
	pcp1 := &estimator.FakePCPredictor{
		PredictFunc: func(context.Context, int, estimator.NodeStatus) (float64, error) {
			return 0.0, estimator.ErrPCPredictor
		},
	}
	intv := 300 * time.Millisecond

	nodes := &estimator.Nodes{}

	// n0: normal case
	n0 := estimator.NewNode("n0", nm0, intv, pcp0)
	status := n0.GetStatus()
	Expect(status.Timestamp).To(Equal(time.Time{}))

	ok := nodes.Add("n0", n0) // nodes.Add() calls Node.start()
	Expect(ok).To(BeTrue())

	Eventually(func() time.Time {
		return n0.GetStatus().Timestamp
	}).ShouldNot(Equal(time.Time{}))

	watt, err := n0.Predict(ctx, 2000, estimator.NodeStatus{
		CPUUsages:    [][]float64{{30.0}},
		AmbientTemps: []float64{20.0},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(watt).To(Equal(70.0))

	// n1: error case
	n1 := estimator.NewNode("n1", nm1, intv, pcp1)
	status = n1.GetStatus()
	Expect(status.Timestamp).To(Equal(time.Time{}))

	ok = nodes.Add("n1", n1)
	Expect(ok).To(BeTrue())

	time.Sleep(intv)
	Eventually(func() time.Time {
		return n1.GetStatus().Timestamp
	}).Should(Equal(time.Time{}))

	_, err = n1.Predict(ctx, 2000, estimator.NodeStatus{
		CPUUsages:    [][]float64{{30.0}},
		AmbientTemps: []float64{20.0},
	})
	Expect(err).To(HaveOccurred())
})
