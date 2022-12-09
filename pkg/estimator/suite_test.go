package estimator_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Estimator Suite")
}

var _ = BeforeSuite(func() {
})

var _ = AfterSuite(func() {
})

var _ = Describe("Server/Client", func() {

	var es *estimator.Estimators
	var sv *estimator.Server
	var hsv *http.Server
	addr := net.JoinHostPort("localhost", estimator.ServerDefaultPort)
	httpAddr := "http://" + addr

	AfterEach(func() {
		err := hsv.Shutdown(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		es = nil
		sv = nil
		hsv = nil
	})

	It("add/delete; no authentication", func() {
		es = &estimator.Estimators{}
		sv = estimator.NewServer(es)
		h, err := sv.Handler()
		Expect(err).NotTo(HaveOccurred())
		hsv = &http.Server{Addr: addr, Handler: h}
		go func() {
			hsv.ListenAndServe()
		}()

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

	It("access; with X-API-KEY authentication", func() {

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

		// // empty
		testAccess(httpAddr, "", "default", "default", false)
		testAccess(httpAddr, "xxx", "default", "default", false)
		testAccess(httpAddr, key1, "default", "default", false)
		testAccess(httpAddr, key2, "default", "default", false)

		// // default/default
		ok := es.Add(client.ObjectKey{Namespace: "default", Name: "default"}.String(), estimator.NewEstimator(&estimator.Nodes{}))
		Expect(ok).To(BeTrue())
		testAccess(httpAddr, "", "default", "default", false)
		testAccess(httpAddr, "xxx", "default", "default", false)
		testAccess(httpAddr, key1, "default", "default", true)
		testAccess(httpAddr, key2, "default", "default", true)

		// empty
		es.Delete(client.ObjectKey{Namespace: "default", Name: "default"}.String())
		testAccess(httpAddr, "", "default", "default", false)
		testAccess(httpAddr, "xxx", "default", "default", false)
		testAccess(httpAddr, key1, "default", "default", false)
		testAccess(httpAddr, key2, "default", "default", false)
	})
})

func testAccess(httpAddr, apiKey, ns, name string, want bool) {
	var opts []estimator.ClientOption
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
