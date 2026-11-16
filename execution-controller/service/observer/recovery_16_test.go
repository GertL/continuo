package observer

import (
	"context"
	"errors"
	"testing"

	pkgmodel "github.com/carolsimone/continuo/pkg/domain/model"
	"github.com/carolsimone/continuo/pkg/seedreload"
	"github.com/stretchr/testify/require"
)

type recoveryLogs struct {
	block     string
	mainError bool
}

func (l recoveryLogs) GetPodLogs(context.Context, string, int64) (string, string, error) {
	if l.mainError {
		return "", "", errors.New("transient main-log lookup failure")
	}
	return "", "", nil
}
func (l recoveryLogs) GetContainerLog(context.Context, string, string) (string, error) {
	return l.block, nil
}

func TestRecoveryPreflightSurvivesMainLogFailure(t *testing.T) {
	for _, mainError := range []bool{false, true} {
		t.Run(map[bool]string{false: "successful_main_log_control", true: "main_log_failed"}[mainError], func(t *testing.T) {
			o := observerWithLogs(t, &fakeLogs{})
			o.deps.Logs = recoveryLogs{block: preflightBlock("error", "unowned_dependents", `"checked":true`), mainError: mainError}
			f := o.readLogs(context.Background(), seedDeployment(seedreload.PreflightCheck))
			require.Equal(t, pkgmodel.FailureReasonUnownedDependents, f.Reason)
			require.NotNil(t, f.SeedPreflight)
		})
	}
}
