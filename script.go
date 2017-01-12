//usr/bin/env go run $0 "$@"; exit
package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"github.com/zanhsieh/custom-k8s-hpa/autoscaler"
	"github.com/zanhsieh/custom-k8s-hpa/options"
	"github.com/zanhsieh/custom-k8s-hpa/version"
)

func main() {
	config := options.NewAutoScalerConfig()
	config.AddFlags(pflag.CommandLine)
	pflag.Parse()
	if config.PrintVer {
		fmt.Printf("%v\n", version.VERSION)
		os.Exit(0)
	}
	if err := config.ValidateFlags(); err != nil {
		glog.Errorf("%v\n", err)
		os.Exit(1)
	}
	scaler, err := autoscaler.NewAutoScaler(config)
        if err != nil {
                glog.Errorf("%v", err)
                os.Exit(1)
        }
        // Begin autoscaling.
        scaler.Run()
}

// Sample curl to prometheus
// curl -sS -g 'http://192.168.99.100:31951/api/v1/query_range?query=avg(container_spec_cpu_period{namespace="b2b-dev-hk",pod_name=~"b2b-web-.*"})&start='$(date +%s --date='5 minutes ago')'&end='$(date +%s)'&step=15s'

// Sample response json
// {"status":"success","data":{"resultType":"matrix","result":[{"metric":{},"values":[[1483953922,"100000"],[1483953937,"100000"],[1483953952,"100000"],[1483953967,"100000"],[1483953982,"100000"],[1483953997,"100000"],[1483954012,"100000"],[1483954027,"100000"],[1483954042,"100000"],[1483954057,"100000"],[1483954072,"100000"],[1483954087,"100000"],[1483954102,"100000"],[1483954117,"100000"],[1483954132,"100000"],[1483954147,"100000"],[1483954162,"100000"],[1483954177,"100000"],[1483954192,"100000"],[1483954207,"100000"],[1483954222,"100000"]]}]}}
