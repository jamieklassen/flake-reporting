# Flake reporting

This should be a prototype for Concourse to track and address test flakes.

## Outcomes

A weekly slack message saying
> Hey! here are your 3 most significant flakey tests this week. Please fix 'em!
> * release-6.3.x > testflight-containerd > archive-pipeline when the step is removed from the parent pipeline config archives the child pipeline (10 times)
> * concourse > bosh-topgun-both > Garbage-collecting volumes A volume that belonged to a container that is now gone is removed from the database and worker [#129726011] (4 times)

## Implementation

* a thing to output test results to your... metrics store?
  * lets focus on getting test results into an opentelemetry-compatible format
  * countervec
  * dimensions include pipeline, job and spec
  * sounds like we will need to export via otlp, or use opentelemetry-go-contrib
    to export to datadog. kinda unforunate that we can't export opencensus since
    we're already using the opencensus agent to forward spans into wavefront
* configuration on your metrics store to visualize
