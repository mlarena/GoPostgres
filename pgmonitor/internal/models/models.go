package models

import "time"

type DatabaseStat struct {
	CheckTime        time.Time
	DBName          string
	SizePretty      string
	SizeBytes       int64
	Collation       string
	ConnectionLimit int
	ConnectionsAllowed bool
}

type LongRunningQuery struct {
	CheckTime       time.Time
	DBName         string
	PID            int
	Username       string
	Application    string
	ClientAddr     string
	BackendStart   time.Time
	QueryStart     time.Time
	Duration       time.Duration
	Query          string
	State          string
}

type Lock struct {
	CheckTime      time.Time
	DBName        string
	BlockedPID    int
	BlockedUser   string
	BlockedQuery  string
	BlockingPID   int
	BlockingUser  string
	BlockingQuery string
	LockType      string
	Mode          string
	Duration      time.Duration
}

type ResourceUsage struct {
	CheckTime            time.Time
	DBName              string
	ActiveConnections   int
	MaxConnections      int
	ConnectionUsagePct  float64
	CacheHitRatio       float64
	TransactionsPerSec  float64
	TuplesFetchedPerSec float64
	TuplesInsertedPerSec float64
	TuplesUpdatedPerSec float64
	TuplesDeletedPerSec float64
}

type ExecutionHistory struct {
	ID             int64
	ExecutionTime  time.Time
	Duration       time.Duration
	DatabasesScanned int
	TablesScanned  int
	Success        bool
	ErrorMessage   string
}