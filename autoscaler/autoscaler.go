//usr/bin/env go run $0 "$@"; exit
package autoscaler

import (
	"fmt"
        "io/ioutil"
        "net/http"
	//"os"
	"time"

	//apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/util/clock"

	"github.com/golang/glog"
        "github.com/gonum/matrix/mat64"
        "github.com/tidwall/gjson"
	"github.com/zanhsieh/custom-k8s-hpa/options"
        "github.com/zanhsieh/custom-k8s-hpa/regression"
)

type AutoScaler struct {
	prometheus	string
	queryExp	string
	step	string
	degree	int
	pollPeriod    time.Duration
	clock         clock.Clock
	stopCh        chan struct{}
	readyCh       chan<- struct{} // For testing.
}

func NewAutoScaler(c *options.AutoScalerConfig) (*AutoScaler, error) {
	thisPollPeriod, err := time.ParseDuration(c.PollPeriod)
	if err != nil {
		return nil, err
	}
	if time.Second > thisPollPeriod {
		thisPollPeriod = time.Second
	}
	return &AutoScaler{
		prometheus:	c.PrometheusIPPort,
		queryExp:	c.QueryExpression,
		step:	c.Step,
		degree:	c.DegPolynomial,
		pollPeriod:    thisPollPeriod,
		clock:         clock.RealClock{},
		stopCh:        make(chan struct{}),
		readyCh:       make(chan struct{}, 1),
	}, nil
}

func (s *AutoScaler) Run() {
	ticker := s.clock.Tick(s.pollPeriod)
	s.readyCh <- struct{}{} // For testing.

	// Don't wait for ticker and execute pollPrometheus() for the first time.
	s.pollPrometheus()

	for {
		select {
		case <-ticker:
			s.pollPrometheus()
		case <-s.stopCh:
			return
		}
	}
}

func (s *AutoScaler) pollPrometheus() {
	debug := false
        now := time.Now()
        minsAgo := now.Add(-5 * time.Minute)
	serverPath := "/api/v1/query_range?query=%v&start=%v&end=%v&step=%v"
        tmp := fmt.Sprintf(serverPath, s.queryExp, int32(minsAgo.Unix()), int32(now.Unix()), s.step)
        url := fmt.Sprintf("http://%v%v", s.prometheus, tmp)
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
        c := regression.Solve(timeSeries, metricVals, s.degree)
        if debug {
                fmt.Printf("c=>\n%.3f\n", mat64.Formatted(c))
        }
        latestTime := timeSeries[len(timeSeries)-1]
        timeOfPredict := 2.0*latestTime - timeSeries[len(timeSeries)-2]
        predictResult := regression.Round(regression.Predict(timeOfPredict, c, s.degree), .5, 2)
        fmt.Printf("predictResult=>%v\n", predictResult)
}
