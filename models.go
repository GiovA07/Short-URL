package main

import (
	"time"
)

type URL struct {
	ID string 			`json:"id"`
	OriginalUrl string 	`json:"url"`
	ShortCode string 	`json:"shortCode"`
	CreateAt time.Time  `json:"createAt"`
	UpdateAt time.Time  `json:"updateAt"`
	AccessCount int 	`json:"accessCount"`
}