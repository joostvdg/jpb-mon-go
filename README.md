# jpb-mon-go

Jenkins Pipeline Binary for Monitoring (Go edition)

## Prometheus Queries

```PromQL
sum(jenkins_pipeline_run_stage) by (jobName, runId) / 1000
```
