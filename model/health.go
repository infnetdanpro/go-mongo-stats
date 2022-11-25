package model

type HealthCheck struct {
	Mongo  bool `json:"mongo"`
	Rabbit bool `json:"rabbit"`
}
