package prometheus

import (
	"fmt"
	"strconv"

	"github.com/joostvdg/jpb-mon-go/pkg/pipelinerun"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	prometheusJob = "pipeline_run_stage"
)

func PushPipelineRunToGateway(prometheusApiEndpoint string, pipelineRun pipelinerun.PipelineRun, pipelineRunMetadata pipelinerun.PipelineRunMetadata) {
	fmt.Printf(" > Attempting to push metrics of %v stages to Prometheus at: %v\n", len(pipelineRun), prometheusApiEndpoint)

	var labelNames []string
	labelNames = []string{"jobName", "runId", "state", "result", "type", "name", "instance"}

	stageMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jenkins_pipeline_run_stage",
		Help: "Duration of Pipeline Run Stage in milliseconds",
	}, labelNames)
	for _, stage := range pipelineRun {

		stageMetric.With(prometheus.Labels{
			"state":    stage.State,
			"result":   stage.Result,
			"type":     stage.Type,
			"name":     stage.DisplayName,
			"instance": pipelineRunMetadata.Instance,
			"jobName":  pipelineRunMetadata.Job,
			"runId":    strconv.Itoa(pipelineRunMetadata.RunId),
		}).Add(float64(stage.DurationInMillis))
	}

	if err := push.New(prometheusApiEndpoint, prometheusJob).
		Collector(stageMetric).
		Push(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}

}
