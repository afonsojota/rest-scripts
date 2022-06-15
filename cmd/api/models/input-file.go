package models

import "bytes"

type InputFile struct {
	Id     int32  `db:"id"`
	Name   string `db:"name"`
	Owner  string `db:"owner"`
	Status string `db:"status"`
	Data   bytes.Buffer
}
