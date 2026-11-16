package observer

import (
	"testing"
	"time"

	"github.com/carolsimone/continuo/execution-controller/domain/model"
	pkgmodel "github.com/carolsimone/continuo/pkg/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMissingJobPreservesKnownOutcome(t *testing.T) {
	now := time.Now()
	for name, tc := range map[string]struct {
		facts   model.AttemptFacts
		want    pkgmodel.FailureReason
		stopped bool
	}{
		"unexplained control": {want: pkgmodel.FailureReasonJobMissing},
		"observed deadline":   {facts: model.AttemptFacts{JobFailureReason: "DeadlineExceeded"}, want: pkgmodel.FailureReasonDeadlineExceeded},
		"observed eviction":   {facts: model.AttemptFacts{PodReason: "Evicted"}, want: pkgmodel.FailureReasonEvicted},
		"requested stop":      {facts: model.AttemptFacts{StopRequestedAt: &now}, stopped: true},
	} {
		t.Run(name, func(t *testing.T) {
			f := tc.facts
			endedFacts(func(next model.AttemptFacts) { f, _ = f.Merge(next) }, func() model.AttemptFacts { return f }, model.JobObservation{}, model.JobObservation{}, model.AttemptFacts{}, model.AttemptFacts{Reason: pkgmodel.FailureReasonJobMissing}, now)
			verdict := model.Classify(f)
			assert.Equal(t, tc.want, verdict.Reason)
			assert.Equal(t, tc.stopped, verdict.Stopped)
		})
	}
}
