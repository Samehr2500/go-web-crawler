package models

import (
	"github.com/google/uuid"
	"time"
)

// Stock schema of the stock table
type Stock struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	PaperName   string    `json:"paper_name"`
	CompanyName string    `json:"company_name"`
	DailyRate   string    `json:"daily_rate"`
	MarketValue float32   `json:"market_value"`
}
