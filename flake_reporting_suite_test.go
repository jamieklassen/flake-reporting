package flake_reporting_test

import (
	"context"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/stdout"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"google.golang.org/grpc"
)

func TestFlakeReporting(t *testing.T) {
	RegisterFailHandler(Fail)
	reporter, err := NewOpenTelemetrySummaryReporter()
	if err != nil {
		Fail(err.Error())
	}
	defer reporter.Stop()
	RunSpecsWithDefaultAndCustomReporters(
		t,
		"FlakeReporting Suite",
		[]Reporter{reporter},
	)
}

func NewOpenTelemetrySummaryReporter() (*OpenTelemetrySummaryReporter, error) {
	exp, err := otlpExporter()
	if err != nil {
		return nil, err
	}
	pusher := push.New(
		basic.New(
			simple.NewWithInexpensiveDistribution(),
			exp,
		),
		exp,
	)
	pusher.Start()
	return &OpenTelemetrySummaryReporter{
		pusher: pusher,
		counter: metric.Must(
			pusher.Provider().Meter("ginkgo"),
		).NewInt64Counter("failures"),
		labels: []kv.KeyValue{
			kv.String("pipeline", "some-pipeline"),
			kv.String("job", "some-test-suite"),
		},
	}, nil
}

func otlpExporter() (export.Exporter, error) {
	return otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:30080"),
		otlp.WithGRPCDialOption(
			grpc.WithBlock(),
			grpc.WithTimeout(1*time.Second),
		),
	)
}

func stdoutExporter() (export.Exporter, error) {
	return stdout.NewExporter(stdout.WithPrettyPrint())
}

type OpenTelemetrySummaryReporter struct {
	pusher  *push.Controller
	counter metric.Int64Counter
	labels  []kv.KeyValue
}

func (otsr *OpenTelemetrySummaryReporter) Stop() {
	otsr.pusher.Stop()
}

func (otsr *OpenTelemetrySummaryReporter) SpecSuiteWillBegin(
	config config.GinkgoConfigType,
	summary *types.SuiteSummary,
) {
}
func (otsr *OpenTelemetrySummaryReporter) BeforeSuiteDidRun(
	setupSummary *types.SetupSummary,
) {
}
func (otsr *OpenTelemetrySummaryReporter) SpecWillRun(
	specSummary *types.SpecSummary,
) {
}
func (otsr *OpenTelemetrySummaryReporter) SpecDidComplete(
	specSummary *types.SpecSummary,
) {
	if specSummary.HasFailureState() {
		specName := strings.Join(specSummary.ComponentTexts[1:], " ")
		otsr.counter.Add(
			context.TODO(),
			1,
			append(otsr.labels, kv.String("spec", specName))...,
		)
	}
}
func (otsr *OpenTelemetrySummaryReporter) AfterSuiteDidRun(
	setupSummary *types.SetupSummary,
) {
}
func (otsr *OpenTelemetrySummaryReporter) SpecSuiteDidEnd(
	summary *types.SuiteSummary,
) {
}
