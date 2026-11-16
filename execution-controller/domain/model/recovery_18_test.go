package model

import (
	pkgmodel "github.com/carolsimone/continuo/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRecoveryInitOOMKeepsResourceFailure(t *testing.T) {
	for _, reason := range []string{"Error", "OOMKilled"} {
		t.Run(reason, func(t *testing.T) {
			obs := JobObservation{Pods: []PodObservation{{Phase: PodFailed, Inits: []ContainerObservation{{Name: "init-work", Terminated: true, ExitCode: 137, TerminatedReason: reason}}}}}
			facts := obs.Facts(time.Now())
			got := Classify(facts)
			want := pkgmodel.FailureReasonExitCode
			if reason == "OOMKilled" {
				want = pkgmodel.FailureReasonOOMKilled
			}
			assert.Equal(t, want, got.Reason)
			assert.Equal(t, "init-work", facts.FailedContainer)
		})
	}
}
