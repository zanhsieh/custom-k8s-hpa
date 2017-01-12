//usr/bin/env go run $0 "$@"; exit
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/gonum/matrix/mat64"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
	"github.com/zanhsieh/custom-k8s-hpa/options"
	"github.com/zanhsieh/custom-k8s-hpa/regression"
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
	var (
		debug      = config.Debug
		prometheus = config.PrometheusIPPort
		serverPath = "/api/v1/query_range?query=%v&start=%v&end=%v&step=%v"
		queryExp   = config.QueryExpression
		step       = config.Step
		degree     = config.DegPolynomial
	)
	now := time.Now()
	minsAgo := now.Add(-5 * time.Minute)
	tmp := fmt.Sprintf(serverPath, queryExp, int32(minsAgo.Unix()), int32(now.Unix()), step)
	url := fmt.Sprintf("http://%v%v", prometheus, tmp)
	if debug {
		fmt.Println(url)
	}
	res, err := http.Get(url)
	if err != nil {
		glog.Fatal(err)
	}
	jsonString, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Fatal(err)
	}
	if debug {
		fmt.Printf("%s", jsonString)
	}

	metric := gjson.GetManyBytes(jsonString, "data.result.0.values.#.0", "data.result.0.values.#.1")
	if debug {
		fmt.Printf("time=>%v\n", metric[0].Array())
		fmt.Printf("value=>%v\n", metric[1].Array())
	}
	timeSeries := make([]float64, len(metric[0].Array()))
	metricVals := make([]float64, len(metric[1].Array()))
	firstElement := 0.0
	for i, time := range metric[0].Array() {
		// This avoids nth power of timestamp cause "matrix singular or near-singular with condition number xxx" problem
		if i == 0 {
			firstElement, timeSeries[i] = time.Float(), 0.0
		} else {
			timeSeries[i] = time.Float() - firstElement
		}
	}
	for j, metricV := range metric[1].Array() {
		metricVals[j] = metricV.Float()
	}
	if debug {
		fmt.Printf("timeSeries=>\n%v\n", timeSeries)
		fmt.Printf("metricVals=>\n%v\n", metricVals)
	}
	c := regression.Solve(timeSeries, metricVals, degree)
	if debug {
		fmt.Printf("c=>\n%.3f\n", mat64.Formatted(c))
	}
	latestTime := timeSeries[len(timeSeries)-1]
	timeOfPredict := 2.0*latestTime - timeSeries[len(timeSeries)-2]
	predictResult := regression.Round(regression.Predict(timeOfPredict, c, degree), .5, 2)
	fmt.Printf("predictResult=>%v\n", predictResult)
}

// Sample curl to prometheus
// curl -sS -g 'http://192.168.99.100:31951/api/v1/query_range?query=avg(container_spec_cpu_period{namespace="b2b-dev-hk",pod_name=~"b2b-web-.*"})&start='$(date +%s --date='5 minutes ago')'&end='$(date +%s)'&step=15s'

// Sample response json
// {"status":"success","data":{"resultType":"matrix","result":[{"metric":{},"values":[[1483953922,"100000"],[1483953937,"100000"],[1483953952,"100000"],[1483953967,"100000"],[1483953982,"100000"],[1483953997,"100000"],[1483954012,"100000"],[1483954027,"100000"],[1483954042,"100000"],[1483954057,"100000"],[1483954072,"100000"],[1483954087,"100000"],[1483954102,"100000"],[1483954117,"100000"],[1483954132,"100000"],[1483954147,"100000"],[1483954162,"100000"],[1483954177,"100000"],[1483954192,"100000"],[1483954207,"100000"],[1483954222,"100000"]]}]}}
