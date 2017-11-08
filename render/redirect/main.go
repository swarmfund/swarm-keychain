package redirect

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type ActionType string

var (
	ActionRecoveryRequest ActionType = "recovery_request"
	ActionSignup          ActionType = "signup"

	// generic redirects

	SignupPayload = Payload{
		Status: http.StatusOK,
		Action: ActionSignup,
	}

	ServerError = Payload{
		Status: http.StatusInternalServerError,
	}

	NotFound = Payload{
		Status: http.StatusNotFound,
	}

	// recovery request

	RecoveryRequestAlreadyUploaded = Payload{
		Status: http.StatusBadGateway,
		Action: ActionRecoveryRequest,
		Data: map[string]interface{}{
			"reason": "already uploaded",
		},
	}

	RecoveryRequestShowCode = func(username, code string) *Payload {
		return &Payload{
			Status: http.StatusOK,
			Action: ActionRecoveryRequest,
			Data: map[string]interface{}{
				"username":     username,
				"code":         code,
				"is_fulfilled": false,
			},
		}
	}

	RecoveryRequestCreateWallet = func(username string) *Payload {
		return &Payload{
			Status: http.StatusOK,
			Action: ActionRecoveryRequest,
			Data: map[string]interface{}{
				"username":     username,
				"is_fulfilled": true,
			},
		}
	}
)

type Payload struct {
	Status int                    `json:"status"`
	Action ActionType             `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

func (p *Payload) Encode() (string, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	if err != nil {
		return "", err
	}
	return encoded, nil
}
