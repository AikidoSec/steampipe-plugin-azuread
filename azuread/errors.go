package azuread

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type RequestError struct {
	StatusCode      int
	Code            string
	Message         string
	RequestID       string `json:",omitempty"`
	ClientRequestID string `json:",omitempty"`
}

func (m *RequestError) Error() string {
	errStr, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(errStr)
}

func getErrorObject(err error) *RequestError {
	var requestErr *RequestError
	if errors.As(err, &requestErr) {
		return requestErr
	}

	return &RequestError{Message: err.Error()}
}

func isIgnorableErrorPredicate(ignoreErrorCodes []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if err != nil {
			if terr, ok := err.(*RequestError); ok {
				for _, item := range ignoreErrorCodes {
					if terr != nil && (terr.Code == item || strings.Contains(terr.Message, item)) {
						return true
					}
				}
			}
		}

		return false
	}
}
