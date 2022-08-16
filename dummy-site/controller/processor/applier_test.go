package processor

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type applierAction string

const (
	actionApply   applierAction = "APPLY"
	actionCleanup applierAction = "CLEANUP"
)

type applierStub struct{}

func (a *applierStub) apply(*manifests) error            { return nil }
func (a *applierStub) cleanupResources(*manifests) error { return nil }

func (a *applierStub) sleepDuration() time.Duration {
	return 5 * time.Millisecond
}

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

	assertStatus(t, []status{
		StatusWorking, StatusWorking, StatusWorking, StatusDone,
	}, statusChan)
}

type applierErroringStub struct {
	failNext bool
	applierStub
}

func (a *applierErroringStub) apply(*manifests) error {
	var res error
	if a.failNext {
		res = errors.New("stub ordered to fail")
	}
	a.failNext = !a.failNext
	return res
}

func TestApplyUntilDestroyed_ReportsErrors(t *testing.T) {
	a := applierErroringStub{}
	term := make(chan bool)
	statusChan := make(chan status, 5)
	go applyUntilDestroyed(&a, &manifests{}, term, statusChan)

	time.Sleep(4*a.sleepDuration() - time.Millisecond)
	term <- true

	assertStatus(t, []status{
		StatusWorking, StatusError, StatusWorking, StatusError, StatusDone,
	}, statusChan)
}

type applierErrorTermStub struct {
	failCount int
	applierStub
}

func (a *applierErrorTermStub) cleanupResources(*manifests) error {
	var res error
	if a.failCount < 3 {
		res = errors.New("stub ordered to fail")
	}
	a.failCount++
	return res
}

func TestApplyUntilDestroyed_RetriesTermination(t *testing.T) {
	a := applierErrorTermStub{}
	term := make(chan bool)
	statusChan := make(chan status, 5)
	go applyUntilDestroyed(&a, &manifests{}, term, statusChan)

	time.Sleep(2*a.sleepDuration() - time.Millisecond)
	term <- true

	assertStatus(t, []status{
		StatusWorking, StatusWorking, StatusError, StatusError, StatusError, StatusDone,
	}, statusChan)
}

func assertStatus(t *testing.T, expected []status, statusChan chan status) {
	var read []status
	for {
		status, ok := <-statusChan
		if ok {
			read = append(read, status)
		} else {
			break
		}
	}

	assert.Equal(t, expected, read)
}
