package models

type Manifest struct {
	Contests []Contest `json:"contests"`
}

type Contest struct {
	Name           string      `json:"name"`
	SelectionLimit int         `json:"selectionLimit"`
	Selections     []Selection `json:"selections"`
}

type Selection struct {
	Name string `json:"name"`
}
