package models

import "time"

type Column struct {
	ID          string    `db:"id" json:"id"`
	SectionID   string    `db:"section_id" json:"sectionId"`
	Name        string    `db:"name" json:"name"`
	Field       string    `db:"field" json:"field"`
	Position    int       `db:"position" json:"position"`
	Type        string    `db:"type" json:"type"`
	Width       int       `db:"width" json:"width"`
	ParentID    string    `db:"parent_id" json:"parentId"`
	AllowSort   bool      `db:"allow_sort" json:"allowSort"`
	AllowFilter bool      `db:"allow_filter" json:"allowFilter"`
	CreatedAt   time.Time `db:"created_at" json:"created"`
	Children    []*Column `json:"children"`
}

type GetColumnsDTO struct {
	SectionID string `json:"sectionId"`
}

type ColumnsDTO struct {
	ID          string `db:"id" json:"id"`
	SectionID   string `db:"section_id" json:"sectionId" binding:"required"`
	Name        string `db:"name" json:"name" binding:"required"`
	Field       string `db:"field" json:"field" binding:"required"`
	Position    int    `db:"position" json:"position"`
	Type        string `db:"type" json:"type" binding:"required"`
	Width       int    `db:"width" json:"width" default:"200"`
	ParentID    string `db:"parent_id" json:"parentId" binding:"omitempty"`
	AllowSort   bool   `db:"allow_sort" json:"allowSort" default:"true"`
	AllowFilter bool   `db:"allow_filter" json:"allowFilter" default:"true"`
}

type DeleteColumnDTO struct {
	ID string `json:"id" db:"id" binding:"required"`
}

type DeleteColumnsDTO struct {
	SectionID string `db:"section_id" json:"sectionId" binding:"required"`
}
