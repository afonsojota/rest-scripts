package utils

import (
	"net/http"
)

type MessageType string

const (
	PENDING  MessageType = "PENDING"
	RUNNING  MessageType = "RUNNING"
	FINISHED MessageType = "FINISHED"
	ERROR    MessageType = "ERROR"
)

var DefaultClientId = "CX_SYSTEM_APPLICATION"

func GetDefaultHeaders() http.Header {
	var header = http.Header{}

	header.Add("X-Admin-Id", DefaultClientId)
	header.Add("Content-Type", "application/json")
	//header.Add("X-Caller-Id", DefaultClientId)
	//header.Add("X-Caller-Scopes", "admin")

	return header
}
