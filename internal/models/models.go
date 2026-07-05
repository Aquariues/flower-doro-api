package models

import "time"

type User struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DeviceID    string    `gorm:"uniqueIndex;size:128" json:"device_id"`
	DisplayName string    `gorm:"size:120" json:"display_name"`
	Locale      string    `gorm:"size:12;default:vi" json:"locale"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserSettings struct {
	UserID       string    `gorm:"type:uuid;primaryKey" json:"user_id"`
	WorkMinutes  int       `gorm:"default:30;not null" json:"work_minutes"`
	BreakMinutes int       `gorm:"default:5;not null" json:"break_minutes"`
	ClockStyle   string    `gorm:"size:40;default:gardenBed;not null" json:"clock_style"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

type Garden struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	UserName  string    `gorm:"size:120;default:You;not null" json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

type Flower struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	Kind                  string    `gorm:"uniqueIndex;size:80;not null" json:"kind"`
	SortOrder             int       `gorm:"uniqueIndex;not null" json:"sort_order"`
	EnglishName           string    `gorm:"size:120;not null" json:"english_name"`
	VietnameseName        string    `gorm:"size:120;not null" json:"vietnamese_name"`
	EnglishDescription    string    `gorm:"type:text" json:"english_description"`
	VietnameseDescription string    `gorm:"type:text" json:"vietnamese_description"`
	EnglishFact1          string    `gorm:"type:text" json:"english_fact_1"`
	EnglishFact2          string    `gorm:"type:text" json:"english_fact_2"`
	EnglishFact3          string    `gorm:"type:text" json:"english_fact_3"`
	VietnameseFact1       string    `gorm:"type:text" json:"vietnamese_fact_1"`
	VietnameseFact2       string    `gorm:"type:text" json:"vietnamese_fact_2"`
	VietnameseFact3       string    `gorm:"type:text" json:"vietnamese_fact_3"`
	Rarity                string    `gorm:"size:40;default:common" json:"rarity"`
	AssetName             string    `gorm:"size:120" json:"asset_name"`
	ImageURL              string    `gorm:"type:text" json:"image_url"`
	ThumbnailURL          string    `gorm:"type:text" json:"thumbnail_url"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type GardenFlower struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	GardenID       *string   `gorm:"type:uuid;index" json:"garden_id,omitempty"`
	UserID         string    `gorm:"type:uuid;index;not null" json:"user_id"`
	FlowerKind     string    `gorm:"size:80;index;not null" json:"flower_kind"`
	FocusSessionID *string   `gorm:"type:uuid;index" json:"focus_session_id,omitempty"`
	FocusMinutes   int       `gorm:"not null" json:"focus_minutes"`
	EarnedAt       time.Time `gorm:"index;not null" json:"earned_at"`
	CreatedAt      time.Time `json:"created_at"`
	Garden         Garden    `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	User           User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

type FocusSession struct {
	ID           string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       string     `gorm:"type:uuid;index;not null" json:"user_id"`
	StartedAt    time.Time  `gorm:"not null" json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	FocusMinutes int        `gorm:"not null" json:"focus_minutes"`
	Rewarded     bool       `gorm:"default:false" json:"rewarded"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	User         User       `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}
