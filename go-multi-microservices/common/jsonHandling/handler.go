package jsonhandling

import (
	"encoding/json"
	"time"

	gh "github.com/go-multi-microservices/common/generalhandling"
)

type Tweet struct {
	Creator   string    `json:"creator,omitempty" binding:"required"`
	Body      string    `json:"body,omitempty" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func JsonMarshal(data interface{}) []byte {
	_json, err := json.Marshal(data)
	gh.HandleError(err, "Error encoding JSON")
	return _json
}
