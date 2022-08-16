package gateway

import (
	"time"

	"github.com/dawidl022/mooc-fi-kubernetes/dummy-site/controller/processor"
)

type applierJob struct {
	website string
	url     string
	term    chan bool
	status  processor.Status
}

func newApplierJob(applier *processor.KubernetesApplier, website string, url string) *applierJob {
	job := &applierJob{
		website: website,
		url:     url,
		term:    make(chan bool),
		status:  processor.StatusInitialising,
	}

	statusChan := make(chan processor.Status, 100)

	go applier.ApplyDummySite(website, url, job.term, statusChan)
	go updateStatus(job, statusChan)

	return job
}

func updateStatus(job *applierJob, statusChan chan processor.Status) {
	for {
		status, ok := <-statusChan
		if ok {
			job.status = status
		} else {
			break
		}
		time.Sleep(processor.SleepDuration)
	}
}
