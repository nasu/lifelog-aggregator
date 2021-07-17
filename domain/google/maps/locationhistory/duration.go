package locationhistory

type Duration struct {
	StartTimestampMs uint64 `json:"startTimestampMs,string"`
	EndTimestampMs   uint64 `json:"endTimestampMs,string"`
}
