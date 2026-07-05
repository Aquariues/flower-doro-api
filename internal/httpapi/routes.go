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
	v1.GET("/users/:id/garden", api.getGarden)
	v1.POST("/users/:id/flowers", api.addGardenFlower)
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
	if err := api.db.Order("id asc").Find(&flowers).Error; err != nil {
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
	Kind                  string `json:"kind" binding:"required"`
	EnglishName           string `json:"english_name" binding:"required"`
	VietnameseName        string `json:"vietnamese_name" binding:"required"`
	EnglishDescription    string `json:"english_description"`
	VietnameseDescription string `json:"vietnamese_description"`
	Rarity                string `json:"rarity"`
	AssetName             string `json:"asset_name"`
}

func (api API) createFlower(c *gin.Context) {
	var req createFlowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	flower := models.Flower{
		Kind:                  req.Kind,
		EnglishName:           req.EnglishName,
		VietnameseName:        req.VietnameseName,
		EnglishDescription:    req.EnglishDescription,
		VietnameseDescription: req.VietnameseDescription,
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
	c.JSON(http.StatusOK, user)
}

func (api API) getGarden(c *gin.Context) {
	var user models.User
	if err := api.db.First(&user, "id = ?", c.Param("id")).Error; err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	var flowers []models.GardenFlower
	if err := api.db.Where("user_id = ?", user.ID).Order("earned_at desc").Find(&flowers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user, "flowers": flowers})
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

	flower := models.GardenFlower{
		UserID:         c.Param("id"),
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

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
