package processor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type applierAction string

const (
	actionApply   applierAction = "APPLY"
	actionCleanup applierAction = "CLEANUP"
)

type applierSpy struct {
	actions []applierAction
	applierStub
}

func (a *applierSpy) apply(*manifests) error {
	a.actions = append(a.actions, actionApply)
	return nil
}

func (a *applierSpy) cleanupResources(*manifests) error {
	a.actions = append(a.actions, actionCleanup)
	return nil
}

type applierStub struct{}

func (a *applierStub) apply(*manifests) error            { return nil }
func (a *applierStub) cleanupResources(*manifests) error { return nil }

func (a *applierStub) sleepDuration() time.Duration {
	return 5 * time.Millisecond
}

func TestApplyUntilDestroyed_RunsUtilTerminated(t *testing.T) {
	a := applierSpy{}
	term := make(chan bool)
	status := make(chan status, 4)
	go applyUntilDestroyed(&a, &manifests{}, term, status)

	time.Sleep(3*a.sleepDuration() - time.Millisecond)
	term <- true

	assert.Equal(t, []applierAction{
		actionApply, actionApply, actionApply, actionCleanup,
	}, a.actions)
}

func TestApplyUntilDestroyed_ReportsStatus(t *testing.T) {
	a := applierStub{}
	term := make(chan bool)
	statusChan := make(chan status, 4)
	go applyUntilDestroyed(&a, &manifests{}, term, statusChan)

	time.Sleep(3*a.sleepDuration() - time.Millisecond)
	term <- true

	var read []status
	for {
		status, ok := <-statusChan
		if ok {
			read = append(read, status)
		} else {
			break
		}
	}

	assert.Equal(t, []status{
		StatusWorking, StatusWorking, StatusWorking, StatusDone,
	}, read)
}
