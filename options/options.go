//usr/bin/env go run $0 "$@"; exit
package options

import (
        "fmt"

        "github.com/golang/glog"
        "github.com/spf13/pflag"
)


type AutoScalerConfig struct {
	Debug	bool
	PrometheusIPPort	string
	QueryExpression	string
	PollPeriod	string
	Step	string
	DegPolynomial	int
	PrintVer	bool
}

func NewAutoScalerConfig() *AutoScalerConfig {
	return &AutoScalerConfig {
		Debug:	false,
		PrometheusIPPort:	"192.168.99.100:30902",
		QueryExpression:	"avg(container_spec_cpu$_period{namespace=\"b2b-dev-hk\",pod_name=~\"b2b-web-.*\"})",
		PollPeriod:	"60s",
		Step:	"15s",
		DegPolynomial: 2,
		PrintVer:	true,
	}
}

func (c *AutoScalerConfig) ValidateFlags() error {
	var errorsFound bool
	if c.PrometheusIPPort == "" {
		errorsFound = true
		glog.Errorf("--prom-ip-port parameter cannot be empty")
	}
	if c.QueryExpression == "" {
		errorsFound = true
		glog.Errorf("--query-exp parameter cannot be empty")
	}
	if errorsFound {
		return fmt.Errorf("failed to validate all input parameters")
	}
	return nil
}

func (c *AutoScalerConfig) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&c.Debug, "debug", c.Debug,  "Enable debug.")
	fs.StringVar(&c.PrometheusIPPort, "prom-ip-port", c.PrometheusIPPort, "Prometheus ip/dns name and listen port number.")
	fs.StringVar(&c.QueryExpression, "query-exp", c.QueryExpression, "Expression used to query Prometheus for scaling.")
	fs.StringVar(&c.PollPeriod, "poll-period", c.PollPeriod, "The time to check Prometheus for metrics and perform autoscale.")
	fs.StringVar(&c.Step, "step", c.Step, "The time interval that Prometheus should report for metrics point for prediction.")
	fs.IntVar(&c.DegPolynomial, "deg-poly", c.DegPolynomial, "In what degree of polynomial curve used to fit historical data.")
	fs.BoolVar(&c.PrintVer, "version", c.PrintVer, "Print the version and exist.")
}
