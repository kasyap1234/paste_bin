package models

import "time"

type PasteFilters struct {
	Languages []string   `json:"languages,omitempty" query:"languages"`
	DateFrom  *time.Time `json:"date_from,omitempty" query:"date_from"`
	DateTo    *time.Time `json:"date_to,omitempty" query:"date_to"`
	SortBy    string     `json:"sort_by,omitempty" query:"sort_by"`
	SortOrder string     `json:"sort_order,omitempty" query:"sort_order"`
}
