package db

type ReportLog struct {
	SourceIp        string `json:"source_ip"`
	ReportTimestamp uint64 `json:"report_timestamp"`
	ReportContent   string `json:"report_content"`
}
