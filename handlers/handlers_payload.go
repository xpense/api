package handlers

import (
	"reflect"
)

// ListResponse is a 'generic' json response whenever the response is a list of something
type ListResponse struct {
	Count   int           `json:"count"`
	Entries []interface{} `json:"entries"`
}

// NewListResponse creates a list response out of a slice of interfaces
func NewListResponse(slice interface{}) *ListResponse {
	entries := []interface{}{}
	if s := reflect.ValueOf(slice); s.Kind() == reflect.Slice {
		for i := 0; i < s.Len(); i++ {
			entry := s.Index(i).Interface()
			entries = append(entries, entry)
		}
	}

	return &ListResponse{
		Count:   len(entries),
		Entries: entries,
	}
}
