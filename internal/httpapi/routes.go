package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/Aquariues/flower-doro-api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type API struct {
	db *gorm.DB
}

func Register(router *gin.Engine, db *gorm.DB) {
	api := API{db: db}

	router.GET("/healthz", api.health)

	v1 := router.Group("/api/v1")
	v1.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"name": "FlowerDoro API", "version": "v1"})
	})
	v1.GET("/flowers", api.listFlowers)
	v1.GET("/flowers/:kind", api.getFlower)
	v1.POST("/flowers", api.createFlower)
	v1.POST("/users", api.createUser)
	v1.GET("/users/:id/settings", api.getSettings)
	v1.PUT("/users/:id/settings", api.updateSettings)
	v1.GET("/users/:id/garden", api.getGarden)
	v1.GET("/users/:id/flower-book", api.getFlowerBook)
	v1.POST("/users/:id/flowers", api.addGardenFlower)
	v1.POST("/users/:id/focus-sessions/complete", api.completeFocusSession)
}

func (api API) health(c *gin.Context) {
	sqlDB, err := api.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if err := sqlDB.PingContext(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (api API) listFlowers(c *gin.Context) {
	var flowers []models.Flower
	if err := api.db.Order("sort_order asc").Find(&flowers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"flowers": flowers})
}

func (api API) getFlower(c *gin.Context) {
	var flower models.Flower
	err := api.db.Where("kind = ?", c.Param("kind")).First(&flower).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "flower not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flower)
}

type createFlowerRequest struct {
	Kind                  string   `json:"kind" binding:"required"`
	SortOrder             int      `json:"sort_order"`
	EnglishName           string   `json:"english_name" binding:"required"`
	VietnameseName        string   `json:"vietnamese_name" binding:"required"`
	EnglishDescription    string   `json:"english_description"`
	VietnameseDescription string   `json:"vietnamese_description"`
	EnglishFacts          []string `json:"english_facts"`
	VietnameseFacts       []string `json:"vietnamese_facts"`
	Rarity                string   `json:"rarity"`
	AssetName             string   `json:"asset_name"`
}

func (api API) createFlower(c *gin.Context) {
	var req createFlowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	flower := models.Flower{
		Kind:                  req.Kind,
		SortOrder:             req.SortOrder,
		EnglishName:           req.EnglishName,
		VietnameseName:        req.VietnameseName,
		EnglishDescription:    req.EnglishDescription,
		VietnameseDescription: req.VietnameseDescription,
		EnglishFact1:          sliceValue(req.EnglishFacts, 0),
		EnglishFact2:          sliceValue(req.EnglishFacts, 1),
		EnglishFact3:          sliceValue(req.EnglishFacts, 2),
		VietnameseFact1:       sliceValue(req.VietnameseFacts, 0),
		VietnameseFact2:       sliceValue(req.VietnameseFacts, 1),
		VietnameseFact3:       sliceValue(req.VietnameseFacts, 2),
		Rarity:                defaultString(req.Rarity, "common"),
		AssetName:             req.AssetName,
	}
	if err := api.db.Create(&flower).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, flower)
}

type createUserRequest struct {
	DeviceID    string `json:"device_id" binding:"required"`
	DisplayName string `json:"display_name"`
	Locale      string `json:"locale"`
}

func (api API) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		DeviceID:    req.DeviceID,
		DisplayName: req.DisplayName,
		Locale:      defaultString(req.Locale, "vi"),
	}
	if err := api.db.Where(models.User{DeviceID: req.DeviceID}).FirstOrCreate(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := api.ensureGardenAndSettings(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (api API) getGarden(c *gin.Context) {
	user, garden, ok := api.userGarden(c)
	if !ok {
		return
	}

	var flowers []models.GardenFlower
	if err := api.db.Where("user_id = ?", user.ID).Order("earned_at desc").Find(&flowers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"garden":  garden,
		"flowers": flowers,
		"summary": api.gardenSummary(user.ID),
		"stats":   focusStats(flowers, time.Now()),
	})
}

func (api API) getSettings(c *gin.Context) {
	var user models.User
	if !api.findUser(c, &user) {
		return
	}

	var settings models.UserSettings
	if err := api.db.Where(models.UserSettings{UserID: user.ID}).FirstOrCreate(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

type updateSettingsRequest struct {
	WorkMinutes  int    `json:"work_minutes" binding:"required,min=1"`
	BreakMinutes int    `json:"break_minutes" binding:"required,min=1"`
	ClockStyle   string `json:"clock_style" binding:"required"`
}

func (api API) updateSettings(c *gin.Context) {
	var user models.User
	if !api.findUser(c, &user) {
		return
	}

	var req updateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings := models.UserSettings{
		UserID:       user.ID,
		WorkMinutes:  req.WorkMinutes,
		BreakMinutes: req.BreakMinutes,
		ClockStyle:   req.ClockStyle,
	}
	if err := api.db.Where(models.UserSettings{UserID: user.ID}).FirstOrCreate(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	settings.WorkMinutes = req.WorkMinutes
	settings.BreakMinutes = req.BreakMinutes
	settings.ClockStyle = req.ClockStyle
	if err := api.db.Model(&settings).Updates(models.UserSettings{
		WorkMinutes:  req.WorkMinutes,
		BreakMinutes: req.BreakMinutes,
		ClockStyle:   req.ClockStyle,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := api.db.First(&settings, "user_id = ?", user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

type flowerBookEntry struct {
	Flower         models.Flower `json:"flower"`
	Unlocked       bool          `json:"unlocked"`
	CollectedCount int64         `json:"collected_count"`
}

func (api API) getFlowerBook(c *gin.Context) {
	var user models.User
	if !api.findUser(c, &user) {
		return
	}

	var flowers []models.Flower
	if err := api.db.Order("sort_order asc").Find(&flowers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	counts := api.flowerCounts(user.ID)
	entries := make([]flowerBookEntry, 0, len(flowers))
	unlocked := 0
	for _, flower := range flowers {
		count := counts[flower.Kind]
		if count > 0 {
			unlocked++
		}
		entries = append(entries, flowerBookEntry{
			Flower:         flower,
			Unlocked:       count > 0,
			CollectedCount: count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"flowers":        entries,
		"total":          len(entries),
		"unlocked_total": unlocked,
	})
}

type addGardenFlowerRequest struct {
	FlowerKind     string     `json:"flower_kind" binding:"required"`
	FocusSessionID *string    `json:"focus_session_id"`
	FocusMinutes   int        `json:"focus_minutes" binding:"required,min=30"`
	EarnedAt       *time.Time `json:"earned_at"`
}

func (api API) addGardenFlower(c *gin.Context) {
	var req addGardenFlowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	earnedAt := time.Now().UTC()
	if req.EarnedAt != nil {
		earnedAt = *req.EarnedAt
	}
	user, garden, ok := api.userGarden(c)
	if !ok {
		return
	}
	if !api.flowerExists(c, req.FlowerKind) {
		return
	}

	flower := models.GardenFlower{
		GardenID:       &garden.ID,
		UserID:         user.ID,
		FlowerKind:     req.FlowerKind,
		FocusSessionID: req.FocusSessionID,
		FocusMinutes:   req.FocusMinutes,
		EarnedAt:       earnedAt,
	}
	if err := api.db.Create(&flower).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, flower)
}

type completeFocusSessionRequest struct {
	FocusMinutes int        `json:"focus_minutes" binding:"required,min=1"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	FlowerKind   string     `json:"flower_kind"`
}

func (api API) completeFocusSession(c *gin.Context) {
	var req completeFocusSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, garden, ok := api.userGarden(c)
	if !ok {
		return
	}

	completedAt := time.Now().UTC()
	if req.CompletedAt != nil {
		completedAt = *req.CompletedAt
	}
	startedAt := completedAt.Add(-time.Duration(req.FocusMinutes) * time.Minute)
	if req.StartedAt != nil {
		startedAt = *req.StartedAt
	}
	eligible := req.FocusMinutes >= 30

	var session models.FocusSession
	var reward *models.GardenFlower
	err := api.db.Transaction(func(tx *gorm.DB) error {
		session = models.FocusSession{
			UserID:       user.ID,
			StartedAt:    startedAt,
			CompletedAt:  &completedAt,
			FocusMinutes: req.FocusMinutes,
			Rewarded:     eligible,
		}
		if err := tx.Create(&session).Error; err != nil {
			return err
		}
		if !eligible {
			return nil
		}

		kind := req.FlowerKind
		if kind == "" {
			if err := tx.Model(&models.Flower{}).Select("kind").Order("random()").Limit(1).Scan(&kind).Error; err != nil {
				return err
			}
		}
		if kind == "" {
			return errors.New("flower catalog is empty")
		}

		flower := models.GardenFlower{
			GardenID:       &garden.ID,
			UserID:         user.ID,
			FlowerKind:     kind,
			FocusSessionID: &session.ID,
			FocusMinutes:   req.FocusMinutes,
			EarnedAt:       completedAt,
		}
		if err := tx.Create(&flower).Error; err != nil {
			return err
		}
		reward = &flower
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session": session,
		"reward":  reward,
	})
}

type gardenSummaryRow struct {
	FlowerKind   string `json:"flower_kind"`
	Count        int64  `json:"count"`
	FocusMinutes int64  `json:"focus_minutes"`
}

func (api API) gardenSummary(userID string) []gardenSummaryRow {
	var rows []gardenSummaryRow
	_ = api.db.Model(&models.GardenFlower{}).
		Select("flower_kind, count(*) as count, coalesce(sum(focus_minutes), 0) as focus_minutes").
		Where("user_id = ?", userID).
		Group("flower_kind").
		Order("count desc, flower_kind asc").
		Scan(&rows).Error
	return rows
}

func (api API) flowerCounts(userID string) map[string]int64 {
	var rows []struct {
		FlowerKind string
		Count      int64
	}
	_ = api.db.Model(&models.GardenFlower{}).
		Select("flower_kind, count(*) as count").
		Where("user_id = ?", userID).
		Group("flower_kind").
		Scan(&rows).Error

	counts := make(map[string]int64, len(rows))
	for _, row := range rows {
		counts[row.FlowerKind] = row.Count
	}
	return counts
}

func (api API) findUser(c *gin.Context, user *models.User) bool {
	if err := api.db.First(user, "id = ?", c.Param("id")).Error; err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func (api API) userGarden(c *gin.Context) (models.User, models.Garden, bool) {
	var user models.User
	if !api.findUser(c, &user) {
		return models.User{}, models.Garden{}, false
	}
	if err := api.ensureGardenAndSettings(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return models.User{}, models.Garden{}, false
	}

	var garden models.Garden
	if err := api.db.Where(models.Garden{UserID: user.ID}).First(&garden).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return models.User{}, models.Garden{}, false
	}
	return user, garden, true
}

func (api API) ensureGardenAndSettings(user *models.User) error {
	garden := models.Garden{UserID: user.ID, UserName: defaultString(user.DisplayName, "You")}
	if err := api.db.Where(models.Garden{UserID: user.ID}).FirstOrCreate(&garden).Error; err != nil {
		return err
	}
	settings := models.UserSettings{UserID: user.ID, WorkMinutes: 30, BreakMinutes: 5, ClockStyle: "gardenBed"}
	return api.db.Where(models.UserSettings{UserID: user.ID}).FirstOrCreate(&settings).Error
}

func (api API) flowerExists(c *gin.Context, kind string) bool {
	var count int64
	if err := api.db.Model(&models.Flower{}).Where("kind = ?", kind).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "flower kind not found"})
		return false
	}
	return true
}

type focusStatsResponse struct {
	Today     int `json:"today"`
	ThisWeek  int `json:"this_week"`
	ThisMonth int `json:"this_month"`
	Total     int `json:"total"`
}

func focusStats(flowers []models.GardenFlower, now time.Time) focusStatsResponse {
	stats := focusStatsResponse{Total: len(flowers)}
	for _, flower := range flowers {
		if sameDay(flower.EarnedAt, now) {
			stats.Today++
		}
		if sameWeek(flower.EarnedAt, now) {
			stats.ThisWeek++
		}
		if flower.EarnedAt.Year() == now.Year() && flower.EarnedAt.Month() == now.Month() {
			stats.ThisMonth++
		}
	}
	return stats
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func sameWeek(a, b time.Time) bool {
	ay, aw := a.ISOWeek()
	by, bw := b.ISOWeek()
	return ay == by && aw == bw
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func sliceValue(values []string, index int) string {
	if index >= len(values) {
		return ""
	}
	return values[index]
}
