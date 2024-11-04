package models

type TaskAction struct {
	ID     int64
	TaskID int64
	Type   ActionType
	URL    string
}

type ActionType string

const (
	commit       ActionType = "COMMIT"
	mergeRequest ActionType = "MERGE_REQUEST"
	merge        ActionType = "MERGE"
)
