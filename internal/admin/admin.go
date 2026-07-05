package admin

import (
	"github.com/Aquariues/flower-doro-api/internal/config"
	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	"github.com/GoAdminGroup/go-admin/engine"
	adminConfig "github.com/GoAdminGroup/go-admin/modules/config"
	adminDB "github.com/GoAdminGroup/go-admin/modules/db"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	_ "github.com/GoAdminGroup/themes/sword"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, cfg config.Config) error {
	eng := engine.Default()

	goAdminConfig := &adminConfig.Config{
		Databases: adminConfig.DatabaseList{
			"default": {
				Driver: adminDB.DriverPostgresql,
				Dsn:    cfg.AdminDatabaseURL,
			},
		},
		AppID:                      "flower-doro-api",
		Language:                   cfg.GoAdminLanguage,
		UrlPrefix:                  cfg.GoAdminPrefix,
		Theme:                      "sword",
		Title:                      cfg.GoAdminTitle,
		Logo:                       "FlowerDoro",
		MiniLogo:                   "FD",
		IndexUrl:                   "/",
		LoginUrl:                   "/login",
		Debug:                      cfg.AppEnv != "production",
		Env:                        cfg.AppEnv,
		InfoLogPath:                "./logs/info.log",
		ErrorLogPath:               "./logs/error.log",
		AccessLogPath:              "./logs/access.log",
		SessionLifeTime:            7200,
		Store:                      adminConfig.Store{Path: "./uploads", Prefix: "uploads"},
		FileUploadEngine:           adminConfig.FileUploadEngine{Name: "local"},
		BootstrapFilePath:          "./bootstrap.go",
		GoModFilePath:              "./go.mod",
		AssetRootPath:              "./public/",
		AllowDelOperationLog:       false,
		OperationLogOff:            false,
		AccessAssetsLogOff:         cfg.AppEnv == "production",
		HideAppInfoEntrance:        true,
		HideToolEntrance:           true,
		HidePluginEntrance:         true,
		ProhibitConfigModification: cfg.AppEnv == "production",
		HideConfigCenterEntrance:   cfg.AppEnv == "production",
	}

	return eng.AddConfig(goAdminConfig).
		AddGenerators(Generators()).
		Use(router)
}

func Generators() map[string]table.Generator {
	return map[string]table.Generator{
		"flowers":        GetFlowersTable,
		"users":          GetUsersTable,
		"garden_flowers": GetGardenFlowersTable,
		"focus_sessions": GetFocusSessionsTable,
	}
}
