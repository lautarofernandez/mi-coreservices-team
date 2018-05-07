package models

import (
	"encoding/json"
	"time"
)

//ExportStatus are the possible states of the export
type ExportStatus string

const (
	// ExportError when the export has a error
	ExportError ExportStatus = "error"

	// ExportPending when the export is pending
	ExportPending ExportStatus = "pending"

	// ExportProcessing when the export is currently processing
	ExportProcessing ExportStatus = "processing"

	// ExportDone when the export is ready to be extracted
	ExportDone ExportStatus = "done"
)

// ExportItem defines the struct of a generic export request
type ExportItem struct {
	ID               string
	Request          json.RawMessage
	NotifyCompletion ExportNotifyCompletion
	ResourceName     string
	FileFormat       ExportFileFormat
	ExportStatus     ExportStatus
	ErrorDescription string
	NotificationSent bool
	RequestedAt      *time.Time
}

// ExportFileFormat are the possible format of the output file
type ExportFileFormat string

const (
	// JSONFormat specified json format to export activities
	JSONFormat ExportFileFormat = "application/json"

	// CSVFormat when the export is pending
	CSVFormat ExportFileFormat = "application/csv"
)

//ExportNotifyCompletion defines struct of the export request
type ExportNotifyCompletion struct {
	BqCluster   string `json:"bq_cluster"`
	BqTopicName string `json:"bq_topic_name" binding:"required"`
}
