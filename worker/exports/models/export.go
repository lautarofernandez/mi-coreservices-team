package models

import (
	"encoding/json"
)

//ExportStatus are the possible states of the export
type ExportStatus string

var (
	//ExportError when the export has a error
	ExportError = ExportStatus("error")
	//ExportPending when the export is pending
	ExportPending = ExportStatus("pending")
	//ExportDone when the export is ready to be extracted
	ExportDone = ExportStatus("done")
)

//ActivityExportItem defines the struct of activities exports
type ExportItem struct {
	ID               string
	Request          json.RawMessage
	NotifyCompletion ExportNotifyCompletion
	ResourceName     string
	FileFormat       ExportFileFormat
	ExportStatus     ExportStatus
	ErrorDescription string
	SendNotification bool
}

//ExportFileFormat are the possible format of the output file
type ExportFileFormat string

var (
	//JSONFormat specified json format to export activities
	JSONFormat = ExportFileFormat("application/json")
	//CSVFormat when the export is pending
	CSVFormat = ExportFileFormat("application/csv")
)

//ExportNotifyCompletion defines struct of the export request
type ExportNotifyCompletion struct {
	BqTopicName string `json:"bq_topic_name" binding:"required"`
}
