package handlers

import (
	"reflect"
)

type ListResponse struct {
	Count   int           `json:"count"`
	Entries []interface{} `json:"entries"`
}

func NewListResponse(slice interface{}) *ListResponse {
	s := reflect.ValueOf(slice)

	var entries []interface{}
	if s.Kind() == reflect.Slice {

		entries = make([]interface{}, s.Len())

		for i := 0; i < s.Len(); i++ {
			entries[i] = s.Index(i).Interface()
		}
	}

	return &ListResponse{
		Count:   len(entries),
		Entries: entries,
	}
}
