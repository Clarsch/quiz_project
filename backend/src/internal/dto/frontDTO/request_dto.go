package frontdto

import (
	"encoding/json"
	"time"
)

type Request struct {
	RequestType  string          `json:"request"`
	Data         json.RawMessage `json:"data"`
	ReceivedTime time.Time       `json:"-"`
}
