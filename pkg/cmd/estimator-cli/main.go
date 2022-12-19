package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
)

type msgWriter struct {
	Writer io.Writer
	Prefix string
	Suffix string
}

func (w *msgWriter) Write(p []byte) (int, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString(w.Prefix)
	if err != nil {
		return 0, fmt.Errorf("could not write to buf: %w", err)
	}
	_, err = buf.Write(p)
	if err != nil {
		return 0, fmt.Errorf("could not write to buf: %w", err)
	}
	_, err = buf.WriteString(w.Suffix)
	if err != nil {
		return 0, fmt.Errorf("could not write to buf: %w", err)
	}
	return w.Writer.Write(buf.Bytes())
}

var errW = &msgWriter{
	Writer: os.Stderr,
	Prefix: "",
	Suffix: "\n",
}

var verbose bool

func v(format string, a ...any) {
	fmt.Fprintf(errW, format, a...)
}

func vv(format string, a ...any) {
	if verbose {
		v(format, a...)
	}
}

func reqPC(ctx context.Context, addr, hk, hv, ns, name string, cpuMilli, numWorkloads int) (*estimator.PowerConsumption, *estimator.Error, error) {
	vv("INFO: estimate power consumption addr=%s hk=%s hv=%s ns=%s name=%s cpu_milli=%d num_workloads=%d", addr, hk, hv, ns, name, cpuMilli, numWorkloads)
	opts := []estimator.ClientOption{}
	if hk != "" && hv != "" {
		opts = append(opts, estimator.ClientOptionAddRequestHeader(hk, hv))
	}
	if verbose {
		opts = append(opts, estimator.ClientOptionGetRequestAsCurl(&msgWriter{
			Writer: os.Stderr,
			Prefix: "DEBUG: ",
			Suffix: "\n",
		}))
	}
	client, err := estimator.NewClient(addr, ns, name, opts...)
	if err != nil {
		return nil, nil, err
	}
	return client.EstimatePowerConsumption(ctx, cpuMilli, numWorkloads)
}

func print(r io.Writer, jsonStructPointer any) error {
	p, err := json.Marshal(jsonStructPointer)
	if err != nil {
		return err
	}
	p = append(p, '\n')
	_, err = r.Write(p)
	return err
}

func csv2Ints(s string) ([]int, error) {
	ss := strings.Split(s, ",")
	ret := make([]int, len(ss))
	for i, e := range ss {
		v, err := strconv.Atoi(e)
		if err != nil {
			return nil, err
		}
		ret[i] = v
	}
	vv("INFO: parse params %s->%v", s, ret)
	return ret, nil
}

func main() {
	nn := flag.String("n", "default/default", "Estimator Namespace/Name")
	addr := flag.String("a", fmt.Sprintf("http://localhost:%d", estimator.ServerDefaultPort), "Estimator address")
	p := flag.String("p", "500,5", "request parameters")
	h := flag.String("H", "", "a request header e.g. 'X-API-KEY: hoge'")
	flag.BoolVar(&verbose, "v", false, "print detailed logs")

	flag.Parse()

	help := func(exitCode int) {
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: %s [option]... <command>\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "\nCommands:\n  pc\testimate power consumption; -p=<cpu_milli>,<num_workloads>\n")
			fmt.Fprintf(os.Stderr, "\nOptions:\n")
			flag.PrintDefaults()
		}
		flag.Usage()
		os.Exit(exitCode)
	}

	if len(os.Args) < 2 {
		help(1)
	}

	// Namespace/Name
	nns := strings.Split(*nn, "/")
	if len(nns) != 2 {
		help(1)
	}
	ns := nns[0]
	name := nns[1]

	// Key: Value
	var hk, hv string
	hs := strings.Split(*h, ":")
	if len(hs) >= 2 {
		hk = strings.TrimSpace(hs[0])
		hv = strings.TrimSpace(hs[1])
	}

	switch os.Args[len(os.Args)-1] {
	case "pc":
		params, err := csv2Ints(*p)
		if err != nil {
			help(1)
		}
		if len(params) != 2 {
			help(1)
		}
		ctx, cncl := context.WithTimeout(context.Background(), 3*time.Second)
		defer cncl()
		pc, apiErr, err := reqPC(ctx, *addr, hk, hv, ns, name, params[0], params[1])
		if err != nil {
			v("ERROR: %v", err)
			os.Exit(1)
		}
		if apiErr != nil {
			v("ERROR:\n  code: %v\n  message: %v", apiErr.Code, apiErr.Message)
			os.Exit(1)
		}
		if err := print(os.Stdout, &pc); err != nil {
			v("ERROR: %v", err)
			os.Exit(1)
		}
	default:
		help(1)
	}
}
