package models

import (
	"time"

	"github.com/goccy/go-json"
)

type ActivityLog struct {
	Id         string          `json:"id" db:"id"`
	TableName  string          `json:"tableName" db:"table_name"`
	RecordId   string          `json:"recordId" db:"record_id"`
	RecordName string          `json:"recordName" db:"record_name"`
	Action     string          `json:"action" db:"action"`
	FieldName  string          `json:"fieldName" db:"field_name"`
	OldValue   json.RawMessage `json:"oldValue" db:"old_value"`
	NewValue   json.RawMessage `json:"newValue" db:"new_value"`
	UserId     string          `json:"userId" db:"user_id"`
	UserName   string          `json:"userName" db:"user_name"`
	CreatedAt  time.Time       `json:"createdAt" db:"created_at"`
}

type CreateActivityLogDTO struct {
	TableName  string      `json:"tableName" db:"table_name"`
	RecordId   string      `json:"recordId" db:"record_id"`
	RecordName string      `json:"recordName" db:"record_name"`
	Action     string      `json:"action" db:"action"`
	FieldName  string      `json:"fieldName" db:"field_name"`
	OldValue   interface{} `json:"oldValue" db:"old_value"`
	NewValue   interface{} `json:"newValue" db:"new_value"`
	UserId     string      `json:"userId" db:"user_id"`
	UserName   string      `json:"userName" db:"user_name"`
}

type GetActivityLogDTO struct {
	TableName string `json:"tableName" db:"table_name"`
	RecordId  string `json:"recordId" db:"record_id"`
	UserId    string `json:"userId" db:"user_id"`
	Limit     int
	Offset    int
}

type ActivityLogFilter struct {
	TableName string
	RecordId  string
	UserId    string
	Action    string
	DateFrom  time.Time
	DateTo    time.Time
	Limit     int
	Offset    int
}
