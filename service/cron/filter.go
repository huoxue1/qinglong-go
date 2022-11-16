package cron

import "github.com/huoxue1/qinglong-go/models"

type Filter interface {
	filter([]*models.Crontabs) []*models.Crontabs
}

type RegFilter struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

type NoRegFilter struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}
