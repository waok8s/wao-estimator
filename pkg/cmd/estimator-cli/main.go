package main

import (
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

var verbose bool

func v(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func vv(format string, a ...any) {
	if verbose {
		v(format, a...)
	}
}

func reqPC(ctx context.Context, addr, ns, name string, cpuMilli, numWorkloads int) (*estimator.PowerConsumption, error) {
	vv("INFO: estimate power consumption addr=%s ns=%s name=%s cpu_milli=%d num_workloads=%d", addr, ns, name, cpuMilli, numWorkloads)
	client, err := estimator.NewClient(addr, ns, name)
	if err != nil {
		return nil, err
	}
	return client.EstimatePowerConsumption(ctx, cpuMilli, numWorkloads)
}

func print(r io.Writer, pc *estimator.PowerConsumption) error {
	p, err := json.Marshal(pc)
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
	addr := flag.String("a", "http://localhost:5678", "Estimator address")
	p := flag.String("p", "500,5", "Request parameters")
	flag.BoolVar(&verbose, "v", false, "Print detailed logs")
	flag.Parse()
	help := func(exitCode int) {
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: %s [option]... <command>\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "\nCommands:\n  pc\tEstimate power consumption. -p=<cpu_milli>,<num_workloads>\n")
			fmt.Fprintf(os.Stderr, "\nOptions:\n")
			flag.PrintDefaults()
		}
		flag.Usage()
		os.Exit(exitCode)
	}
	if len(os.Args) < 2 {
		help(1)
	}
	nns := strings.Split(*nn, "/")
	if len(nns) != 2 {
		help(1)
	}
	ns := nns[0]
	name := nns[1]
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
		pc, err := reqPC(ctx, *addr, ns, name, params[0], params[1])
		if err != nil {
			v("ERROR: %v", err)
			os.Exit(1)
		}
		if err := print(os.Stdout, pc); err != nil {
			v("ERROR: %v", err)
			os.Exit(1)
		}
	default:
		help(1)
	}
}
