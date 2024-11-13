package models

type TaskAction struct {
	ID         int64      `json:"id" db:"id"`
	ExternalID int64      `json:"external_id" db:"external_id"`
	TaskID     int64      `json:"task_id" db:"task_id"`
	Type       ActionType `json:"type" db:"type"`
	URL        string     `json:"url" db:"url"`
}

type ActionType string

const (
	unspecifiedActionType  ActionType = "UNSPECIFIED"
	CommitActionType       ActionType = "COMMIT"
	mergeRequestActionType ActionType = "MERGE_REQUEST"
	mergeActionType        ActionType = "MERGE"
)
