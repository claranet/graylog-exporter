
package graylog

import (
	"fmt"
	"encoding/json"
)

type Journal struct {
	Enabled				bool `json:"enabled"`
	AppendEventsPerSecond		int  `json:"append_events_per_second"`
	ReadEventsPerSecond		int  `json:"read_events_per_second"`
	UncommittedJournalEntries	int  `json:"uncommitted_journal_entries"`
	JournalSize			int  `json:"journal_size"`
	JournalSizeLimit		int  `json:"journal_size_limit"`
	NumberOfSegments		int  `json:"number_of_segments"`
	OldestSegment			string `json:"oldest_segment"`
	JournalConfig			*JournalConfig `json:"journal_config"`
}

type JournalConfig struct {
	Directory			string
	SegmentSize			int
	SegmentAge			int
	MaxSize				int
	MaxAge				int
	FlushInterval			int
	FlushAge			int
}

func (g *Graylog) GetClusterJournal(nodeId string) Journal {
	res, _ := g.makeRequest("GET", fmt.Sprintf("/cluster/%s/journal", nodeId))
	data := json.NewDecoder(res.Body)

	var d Journal
	data.Decode(&d)

	return d
}
func (g *Graylog) GetSystemJournal() Journal {
	res, _ := g.makeRequest("GET", "/system/journal")
	data := json.NewDecoder(res.Body)

	var d Journal
	data.Decode(&d)

	return d
}
