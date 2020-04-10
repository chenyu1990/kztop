package kreedz

import (
	"context"
	"net/http"
	"time"
)

type Organization uint

const (
	_ Organization = iota
	XtremeJumps
	CosyClimbing
	WorldSurf
	DebugWorldRecord = 0
)

type WorldRecord struct {
	Context          context.Context
	Organization     Organization
	recordFileHeader http.Header
	NewRecords       []*RecordInfo
	News             map[string]map[string][]*RecordInfo
	NewsDate         time.Time
}

type RecordInfo struct {
	MapName string
	Holder  string
	Country string
	Time    float64
}
