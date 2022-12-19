package model

import "go.mongodb.org/mongo-driver/bson"

type Filters struct {
	Filters bson.M `json:"filters"`
}

type DistinctFilters struct {
	Field   string `json:"field"`
	Filters bson.M `json:"filters,omitempty"`
}
