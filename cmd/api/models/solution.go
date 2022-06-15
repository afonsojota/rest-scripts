package models

type Solution struct {
	Id       int64     `json:"id"`
	Versions []Version `json:"versions"`
}

type Version struct {
	Sites    []string  `json:"sites"`
	Contents []Content `json:"contents"`
}

type Content struct {
	Channel string `json:"channel"`
	Body    Body   `json:"content"`
}

type Body struct {
	Template string `json:"body"`
}
