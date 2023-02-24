package estimator_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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
	SetDefaultEventuallyTimeout(1 * time.Second)
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

		testAccess(addr, "", "default", "default", nil, true)
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
		testAccess(httpAddr, "", "default", "default", estimator.ErrServerEstimatorNotFound, false)
		testAccess(httpAddr, "", "hoge", "fuga", estimator.ErrServerEstimatorNotFound, false)

		// default/default
		ok := es.Add(estimator.RequestToEstimatorName("default", "default"), &estimator.Estimator{Nodes: nil})
		Expect(ok).To(BeTrue())
		testAccess(httpAddr, "", "default", "default", estimator.ErrEstimatorNoNodesAvailable, false)
		testAccess(httpAddr, "", "hoge", "fuga", estimator.ErrServerEstimatorNotFound, false)

		// default/default, hoge/fuga
		ok = es.Add(estimator.RequestToEstimatorName("hoge", "fuga"), &estimator.Estimator{Nodes: nil})
		Expect(ok).To(BeTrue())
		testAccess(httpAddr, "", "default", "default", estimator.ErrEstimatorNoNodesAvailable, false)
		testAccess(httpAddr, "", "hoge", "fuga", estimator.ErrEstimatorNoNodesAvailable, false)

		// default/default, hoge/fuga
		es.Delete(estimator.RequestToEstimatorName("foo", "bar"))
		testAccess(httpAddr, "", "default", "default", estimator.ErrEstimatorNoNodesAvailable, false)
		testAccess(httpAddr, "", "hoge", "fuga", estimator.ErrEstimatorNoNodesAvailable, false)

		// default/default
		es.Delete(estimator.RequestToEstimatorName("hoge", "fuga"))
		testAccess(httpAddr, "", "default", "default", estimator.ErrEstimatorNoNodesAvailable, false)
		testAccess(httpAddr, "", "hoge", "fuga", estimator.ErrServerEstimatorNotFound, false)

		// empty
		es.Delete(estimator.RequestToEstimatorName("default", "default"))
		testAccess(httpAddr, "", "default", "default", estimator.ErrServerEstimatorNotFound, false)
		testAccess(httpAddr, "", "hoge", "fuga", estimator.ErrServerEstimatorNotFound, false)
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
		testAccess(httpAddr, "", "default", "default", estimator.ErrClientUnauthorized, false)
		testAccess(httpAddr, "xxx", "default", "default", estimator.ErrClientUnauthorized, false)
		testAccess(httpAddr, key1, "default", "default", estimator.ErrServerEstimatorNotFound, false)
		testAccess(httpAddr, key2, "default", "default", estimator.ErrServerEstimatorNotFound, false)

		// Estimators: default/default
		ok := es.Add(estimator.RequestToEstimatorName("default", "default"), &estimator.Estimator{Nodes: nil})
		Expect(ok).To(BeTrue())
		testAccess(httpAddr, "", "default", "default", estimator.ErrClientUnauthorized, false)
		testAccess(httpAddr, "xxx", "default", "default", estimator.ErrClientUnauthorized, false)
		testAccess(httpAddr, key1, "default", "default", estimator.ErrEstimatorNoNodesAvailable, false)
		testAccess(httpAddr, key2, "default", "default", estimator.ErrEstimatorNoNodesAvailable, false)

		// Estimators: empty
		es.Delete(estimator.RequestToEstimatorName("default", "default"))
		testAccess(httpAddr, "", "default", "default", estimator.ErrClientUnauthorized, false)
		testAccess(httpAddr, "xxx", "default", "default", estimator.ErrClientUnauthorized, false)
		testAccess(httpAddr, key1, "default", "default", estimator.ErrServerEstimatorNotFound, false)
		testAccess(httpAddr, key2, "default", "default", estimator.ErrServerEstimatorNotFound, false)
	})

	It("should request", func() {

		ns := "default"
		name := "default"

		// client
		opts := []estimator.ClientOption{}
		cl, err := estimator.NewClient(httpAddr, ns, name, opts...)
		Expect(err).NotTo(HaveOccurred())
		// server
		es = &estimator.Estimators{}
		sv = &estimator.Server{Estimators: es}
		h, err := sv.Handler()
		Expect(err).NotTo(HaveOccurred())
		hsv = &http.Server{Addr: addr, Handler: h}
		go func() {
			hsv.ListenAndServe()
		}()
		wait()
		// estimator
		est := &estimator.Estimator{}
		sv.Estimators.Add(estimator.RequestToEstimatorName(ns, name), est)

		// test: no nodes
		testRequest(cl, &estimator.PowerConsumption{CpuMilli: 500, NumWorkloads: 5}, nil, estimator.ErrEstimatorNoNodesAvailable)

		// test: n0 (no StatusMonitor, no PCPredictor)
		intv := 300 * time.Millisecond
		n0 := estimator.NewNode("n0", nil, intv, nil)
		est.Nodes.Add(n0.Name, n0)
		testRequest(cl, &estimator.PowerConsumption{
			CpuMilli: 500, NumWorkloads: 5,
		}, &estimator.PowerConsumption{
			CpuMilli: 500, NumWorkloads: 5, WattIncreases: &[]float64{math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1)},
		}, nil)
		testRequest(cl, &estimator.PowerConsumption{
			CpuMilli: 1000, NumWorkloads: 1,
		}, &estimator.PowerConsumption{
			CpuMilli: 1000, NumWorkloads: 1, WattIncreases: &[]float64{math.Inf(1)},
		}, nil)
		testRequest(cl, &estimator.PowerConsumption{
			CpuMilli: 1000, NumWorkloads: 0,
		}, &estimator.PowerConsumption{
			CpuMilli: 1000, NumWorkloads: 0, WattIncreases: &[]float64{},
		}, nil)

		// test: n0, n1 (fake)
		nm1 := &estimator.FakeNodeMonitor{FetchFunc: func(context.Context, *estimator.NodeStatus) error { return nil }}
		pcp1 := &estimator.FakePCPredictor{PredictFunc: func(_ context.Context, requestCPUMilli int, _ *estimator.NodeStatus) (watt float64, err error) {
			// 100mCPU/W
			return float64(requestCPUMilli) / 100, nil
		}}
		n1 := estimator.NewNode("n1", []estimator.NodeMonitor{nm1}, intv, pcp1)
		est.Nodes.Add(n1.Name, n1)
		// n0: [inf inf inf inf]
		// n1: [  5  10  15  20]
		testRequest(cl, &estimator.PowerConsumption{
			CpuMilli: 500, NumWorkloads: 4,
		}, &estimator.PowerConsumption{
			CpuMilli: 500, NumWorkloads: 4, WattIncreases: &[]float64{5, 10, 15, 20},
		}, nil)

		// test: n0, n1, n2 (fake)
		nm2 := &estimator.FakeNodeMonitor{FetchFunc: func(context.Context, *estimator.NodeStatus) error { return nil }}
		pcp2 := &estimator.FakePCPredictor{PredictFunc: func(_ context.Context, requestCPUMilli int, _ *estimator.NodeStatus) (watt float64, err error) {
			// 200mCPU/W
			return float64(requestCPUMilli) / 200, nil
		}}
		n2 := estimator.NewNode("n2", []estimator.NodeMonitor{nm2}, intv, pcp2)
		est.Nodes.Add(n2.Name, n2)
		// n0: [inf inf inf inf]
		// n1: [  5  10  15  20]
		// n2: [2.5   5 7.5  10]
		testRequest(cl, &estimator.PowerConsumption{
			CpuMilli: 500, NumWorkloads: 4,
		}, &estimator.PowerConsumption{
			CpuMilli: 500, NumWorkloads: 4, WattIncreases: &[]float64{2.5, 5, 7.5, 10},
		}, nil)

	})

})

func testAccess(httpAddr, apiKey, ns, name string, wantAPIErr error, wantErr bool) {
	opts := []estimator.ClientOption{estimator.ClientOptionGetRequestAsCurl(GinkgoWriter)}
	if apiKey != "" {
		opts = append(opts, estimator.ClientOptionAddRequestHeader(estimator.AuthFnAPIKeyRequestHeader, apiKey))
	}
	cl, err := estimator.NewClient(httpAddr, ns, name, opts...)
	Expect(err).NotTo(HaveOccurred())

	testFn := func() error {
		pc, apiErr, err := cl.EstimatePowerConsumption(context.Background(), 500, 5)

		if wantErr {
			if err == nil || pc != nil || apiErr != nil {
				return fmt.Errorf("wantErr=%v but got %v, pc=%v apiErr=%v", wantErr, err, pc, apiErr)
			}
			return nil
		}

		if wantAPIErr != nil {
			if apiErr == nil || pc != nil || err != nil {
				return fmt.Errorf("wantAPIErr=%v but got %v, pc=%v err=%v", wantAPIErr, apiErr, pc, err)
			}
			if !errors.Is(estimator.GetErrorFromCode(*apiErr), wantAPIErr) {
				return fmt.Errorf("wantAPIErr=%v but got %v, pc=%v err=%v", wantAPIErr, apiErr, pc, err)
			}
			return nil
		}

		if !wantErr && wantAPIErr == nil {
			if pc == nil || apiErr != nil || err != nil {
				return fmt.Errorf("want response but got nil, pc=%v, apiErr=%v, err=%v", pc, apiErr, err)
			}
			return nil
		}

		return fmt.Errorf("unexpected wantAPIErr=%v wantErr=%v pc=%v, apiErr=%v, err=%v", wantAPIErr, wantErr, pc, apiErr, err)
	}

	Eventually(testFn).Should(Succeed())
}

func testRequest(cl *estimator.Client, req, want *estimator.PowerConsumption, wantAPIErr error) {
	testFn := func() error {
		if req == nil {
			return fmt.Errorf("req must not be nil")
		}

		pc, apiErr, err := cl.EstimatePowerConsumption(context.Background(), req.CpuMilli, req.NumWorkloads)

		if err != nil {
			return fmt.Errorf("err=%v", err)
		}

		if wantAPIErr != nil {
			if apiErr == nil || pc != nil || err != nil {
				return fmt.Errorf("wantAPIErr=%v but got %v, pc=%v err=%v", wantAPIErr, apiErr, pc, err)
			}
			if !errors.Is(estimator.GetErrorFromCode(*apiErr), wantAPIErr) {
				return fmt.Errorf("wantAPIErr=%v but got %v, pc=%v err=%v", wantAPIErr, apiErr, pc, err)
			}
			return nil
		}

		if !reflect.DeepEqual(pc, want) {
			return fmt.Errorf("want=%v (WattIncreases=%v) but got %v (WattIncreases=%v)", want, want.WattIncreases, pc, pc.WattIncreases)
		}

		return nil
	}

	Eventually(testFn).Should(Succeed())
}

var _ = Describe("Node/Nodes", func() {

	ctx := context.Background()

	nm0 := &estimator.FakeNodeMonitor{
		FetchFunc: func(ctx context.Context, base *estimator.NodeStatus) error {
			time.Sleep(50 * time.Millisecond)
			setNodeStatus(base, 10.0, 20.0, 0.0)
			return nil
		},
	}
	nm1 := &estimator.FakeNodeMonitor{
		FetchFunc: func(ctx context.Context, base *estimator.NodeStatus) error {
			time.Sleep(50 * time.Millisecond)
			return estimator.ErrNodeMonitor
		},
	}
	pcp0 := &estimator.FakePCPredictor{
		PredictFunc: estimator.PredictPCFnDummy,
	}
	pcp1 := &estimator.FakePCPredictor{
		PredictFunc: func(context.Context, int, *estimator.NodeStatus) (float64, error) {
			return 0.0, estimator.ErrPCPredictor
		},
	}
	intv := 300 * time.Millisecond

	nodes := &estimator.Nodes{}

	// n0: normal case
	n0 := estimator.NewNode("n0", []estimator.NodeMonitor{nm0}, intv, pcp0)
	Eventually(func() time.Time {
		return n0.GetStatus().Timestamp()
	}).ShouldNot(Equal(time.Time{}))

	ok := nodes.Add("n0", n0) // nodes.Add() calls Node.start()
	Expect(ok).To(BeTrue())

	Eventually(func() time.Time {
		return n0.GetStatus().Timestamp()
	}).ShouldNot(Equal(time.Time{}))

	watt, err := n0.Predict(ctx, 2000, setNodeStatus(nil, 30.0, 20.0, 0.0))
	Expect(err).NotTo(HaveOccurred())
	Expect(watt).To(Equal(70.0))

	// n1: error case
	n1 := estimator.NewNode("n1", []estimator.NodeMonitor{nm1}, intv, pcp1)
	Eventually(func() time.Time {
		return n1.GetStatus().Timestamp()
	}).ShouldNot(Equal(time.Time{}))

	ok = nodes.Add("n1", n1)
	Expect(ok).To(BeTrue())

	time.Sleep(intv)
	Eventually(func() time.Time {
		return n1.GetStatus().Timestamp()
	}).ShouldNot(Equal(time.Time{}))

	_, err = n1.Predict(ctx, 2000, setNodeStatus(nil, 30.0, 20.0, 0.0))
	Expect(err).To(HaveOccurred())
})

func setNodeStatus(base *estimator.NodeStatus, cpuUsage, ambientTemp, staticPressureDiff float64) *estimator.NodeStatus {
	if base == nil {
		base = estimator.NewNodeStatus()
	}
	estimator.NodeStatusSetCPUUsage(base, cpuUsage)
	estimator.NodeStatusSetAmbientTemp(base, ambientTemp)
	estimator.NodeStatusSetStaticPressureDiff(base, staticPressureDiff)
	return base
}
