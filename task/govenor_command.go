// Copyright 2020
// limitations under the License.

package task

import (
	"context"
	"errors"
	"fmt"
	"html/template"

	"github.com/beego/beego/v2/core/admin"
)

type listTaskCommand struct{}

func (l *listTaskCommand) Execute(params ...interface{}) *admin.Result {
	resultList := make([][]string, 0, len(globalTaskManager.adminTaskList))
	for tname, tk := range globalTaskManager.adminTaskList {
		result := []string{
			template.HTMLEscapeString(tname),
			template.HTMLEscapeString(tk.GetSpec(nil)),
			template.HTMLEscapeString(tk.GetStatus(nil)),
			template.HTMLEscapeString(tk.GetPrev(context.Background()).String()),
		}
		resultList = append(resultList, result)
	}

	return &admin.Result{
		Status:  200,
		Content: resultList,
	}
}

type runTaskCommand struct{}

func (r *runTaskCommand) Execute(params ...interface{}) *admin.Result {
	if len(params) == 0 {
		return &admin.Result{
			Status: 400,
			Error:  errors.New("task name not passed"),
		}
	}

	tn, ok := params[0].(string)

	if !ok {
		return &admin.Result{
			Status: 400,
			Error:  errors.New("parameter is invalid"),
		}
	}

	if t, ok := globalTaskManager.adminTaskList[tn]; ok {
		err := t.Run(context.Background())
		if err != nil {
			return &admin.Result{
				Status: 500,
				Error:  err,
			}
		}
		return &admin.Result{
			Status:  200,
			Content: t.GetStatus(context.Background()),
		}
	} else {
		return &admin.Result{
			Status: 400,
			Error:  fmt.Errorf("task with name %s not found", tn),
		}
	}
}

func registerCommands() {
	admin.RegisterCommand("task", "list", &listTaskCommand{})
	admin.RegisterCommand("task", "run", &runTaskCommand{})
}
