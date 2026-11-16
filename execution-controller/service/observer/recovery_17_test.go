package observer

import (
	"context"
	"testing"
	"time"

	"github.com/carolsimone/continuo/execution-controller/domain/command"
	"github.com/carolsimone/continuo/execution-controller/domain/model"
	pkgmodel "github.com/carolsimone/continuo/pkg/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryPythonLogPhraseDoesNotFailSuccess(t *testing.T) {
	for _, nodeType := range []pkgmodel.NodeType{pkgmodel.NodeTypePythonNode, pkgmodel.NodeTypePythonCsv, pkgmodel.NodeTypePythonApi, pkgmodel.NodeTypeDbtModel} {
		for _, candidate := range []bool{false, true} {
			mode := "production"
			if candidate {
				mode = "validation"
			}
			t.Run(string(nodeType)+"/"+mode, func(t *testing.T) {
				o := observerWithLogs(t, &fakeLogs{full: "Nothing to do with this diagnostic: processed 4 rows"})
				d := model.NewDeployment(command.DeployTask{NodeType: string(nodeType), ServiceName: "analytics", SchemaName: "public", TableName: "orders"}, nil, time.Now())
				if candidate {
					d = model.NewValidationDeployment(command.ValidationDeployTask{NodeType: string(nodeType), ServiceName: "analytics", SchemaName: "public", TableName: "orders"}, nil, time.Now(), false)
				}
				f := o.readLogs(context.Background(), d)
				code := int32(0)
				f.ExitCode = &code
				verdict := model.Classify(f)
				if nodeType.IsPython() {
					assert.True(t, verdict.Succeeded, "successful Python task must not become a dbt selection failure")
				} else {
					assert.Equal(t, pkgmodel.FailureReasonDbtNothingToDo, verdict.Reason, "dbt control retains no-selection handling")
				}
			})
		}
	}
}
