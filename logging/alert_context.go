package logging

import (
	"log/slog"

	. "github.com/atropos112/gocore/types"
	"github.com/atropos112/gocore/utils"
)

type PubSubError struct {
	Message string         `json:"message"`
	Args    map[string]any `json:"args"`
}

type AlertContext struct {
	Logger  *slog.Logger
	Publish func(pub_obj PublishableObject, source string)
	Source  string
}

func (a *AlertContext) Error(msg string, args ...any) {
	a.Logger.Error(msg, args...)
	argsMap := make(map[string]any)
	for i := 0; i < len(args); i += 2 {
		argsMap[args[i].(string)] = args[i+1]
	}

	a.Publish(PubSubError{
		Message: msg,
		Args:    argsMap,
	}, a.Source)
}

func (a *AlertContext) ErrorAndDie(msg string, args ...any) {
	a.Error(msg, args...)
	panic("Failed with message: " + msg)
}

func (a *AlertContext) GetCredOrAlertAndDie(value string) string {
	cred, err := utils.GetCred(value)
	if err != nil {
		a.ErrorAndDie("Failed to get credential", "error", err)
	}

	return cred
}
