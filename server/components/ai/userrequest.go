package ai

import "laatoo.io/sdk/utils"

type UserRequest struct {
	Agent          string          `json:"agent"`
	Message        string          `json:"message"`
	Context        utils.StringMap `json:"context"`
	AdditionalInfo utils.StringMap `json:"additionalinfo"`
	SessionId      string          `json:"sessionId"`
}
