package models

import "time"

type Column struct {
	ID          string    `db:"id" json:"id"`
	SectionID   string    `db:"section_id" json:"section_id"`
	Name        string    `db:"name" json:"name"`
	Field       string    `db:"field" json:"field"`
	Position    int       `db:"position" json:"position"`
	Type        string    `db:"type" json:"type"`
	Width       int       `db:"width" json:"width"`
	ParentID    string    `db:"parent_id" json:"parent_id"`
	AllowSort   bool      `db:"allow_sort" json:"allow_sort"`
	AllowFilter bool      `db:"allow_filter" json:"allow_filter"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Children    []*Column `json:"children"`
}

type GetColumnsDTO struct {
	SectionID string `json:"section_id"`
}

type ColumnsDTO struct {
	ID          string `db:"id" json:"id" binding:"required"`
	SectionID   string `db:"section_id" json:"section_id" binding:"required"`
	Name        string `db:"name" json:"name" binding:"required"`
	Field       string `db:"field" json:"field" binding:"required"`
	Position    int    `db:"position" json:"position" binding:"required"`
	Type        string `db:"type" json:"type" binding:"required"`
	Width       int    `db:"width" json:"width" default:"200"`
	ParentID    string `db:"parent_id" json:"parent_id" binding:"omitempty"`
	AllowSort   bool   `db:"allow_sort" json:"allow_sort" default:"true"`
	AllowFilter bool   `db:"allow_filter" json:"allow_filter" default:"true"`
}

type DeleteColumnDTO struct {
	ID string `json:"id" db:"id" binding:"required"`
}
