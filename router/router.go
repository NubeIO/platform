package router

import (
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	authconstants "github.com/NubeIO/nubeio-rubix-lib-auth-go/constants"
	"github.com/NubeIO/platform/config"
	"github.com/NubeIO/platform/constants"
	"github.com/NubeIO/platform/controller"
	"github.com/NubeIO/platform/logger"
	"github.com/NubeIO/platform/model"
	"github.com/NubeIO/platform/services/appstore"
	"github.com/NubeIO/platform/services/info"
	systeminfo "github.com/NubeIO/platform/services/system"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func NotFound() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, http.StatusNotFound, "rubix-platform: api not found")
		ctx.JSON(http.StatusNotFound, model.Message{Message: message})
	}
}

func Setup() *gin.Engine {
	engine := gin.New()
	// Set gin access logs
	if viper.GetBool("gin.log.store") {
		fileLocation := fmt.Sprintf("%s/rubix-platform.access.log", config.Config.GetAbsDataDir())
		f, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, constants.Permission)
		if err != nil {
			logger.Logger.Errorf("Failed to create access log file: %v", err)
		} else {
			gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		}
	}
	gin.SetMode(viper.GetString("gin.log.level"))
	engine.NoRoute(NotFound())
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"},
		AllowHeaders: []string{
			"Authorization", "Content-Type", "Upgrade", "Origin", "Connection", "Accept-Encoding", "Accept-Language",
			"Host", "Referer", "User-Agent", "Accept", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers",
			"X-Host",
		},
		ExposeHeaders:          []string{"Content-Length"},
		AllowCredentials:       true,
		AllowAllOrigins:        true,
		AllowBrowserExtensions: true,
		MaxAge:                 12 * time.Hour,
	}))

	systemCtl := systemctl.New(false, 30)
	systemInfo := systeminfo.New()
	api := controller.Controller{
		SystemCtl:  systemCtl,
		FileMode:   0755,
		Instances:  make(map[string]*controller.Instance),
		Lock:       sync.Mutex{},
		Config:     config.Config,
		SystemInfo: systemInfo,
		Networking: info.New(&info.System{}),
		Store:      appstore.New(fmt.Sprintf("/%s", config.Config.GetAbsDataDir())),
	}
	err := api.LoadFromFile("./db.yaml")
	if err != nil {
		log.Fatal(err)
	}
	engine.POST("/api/users/login", api.Login)
	systemApi := engine.Group("/api/system")
	{
		systemApi.GET("/ping", api.Ping)
		systemApi.GET("/arch", api.GetArch)
	}

	handleAuth := func(c *gin.Context) { c.Next() }
	handleUserAuth := func(c *gin.Context) { c.Next() }
	if config.Config.Auth() {
		handleAuth = api.HandleAuth(authconstants.UserRole)
		handleUserAuth = api.HandleUserAuth(authconstants.UserRole)
	}

	apiProxyROSRoutes := engine.Group("/ros")
	apiProxyROSRoutes.Any("/*proxy_path", api.ROSProxy)

	apiRoutes := engine.Group("/api", handleAuth)

	systemRoutes := apiRoutes.Group("/system")
	{
		systemRoutes.GET("/info", api.GetSystemInfo)
		systemRoutes.POST("/reboot", api.RebootHost)
	}

	appControl := apiRoutes.Group("/systemctl")
	{
		appControl.POST("/enable", api.SystemCtlEnable)
		appControl.POST("/disable", api.SystemCtlDisable)
		appControl.GET("/show", api.SystemCtlShow)
		appControl.POST("/start", api.SystemCtlStart)
		appControl.GET("/status", api.SystemCtlStatus)
		appControl.POST("/stop", api.SystemCtlStop)
		appControl.POST("/reset-failed", api.SystemCtlResetFailed)
		appControl.POST("/daemon-reload", api.SystemCtlDaemonReload)
		appControl.POST("/restart", api.SystemCtlRestart)
		appControl.POST("/mask", api.SystemCtlMask)
		appControl.POST("/unmask", api.SystemCtlUnmask)
		appControl.GET("/state", api.SystemCtlState)
		appControl.GET("/is-enabled", api.SystemCtlIsEnabled)
		appControl.GET("/is-active", api.SystemCtlIsActive)
		appControl.GET("/is-running", api.SystemCtlIsRunning)
		appControl.GET("/is-failed", api.SystemCtlIsFailed)
		appControl.GET("/is-installed", api.SystemCtlIsInstalled)
	}

	syscallControl := apiRoutes.Group("/syscall")
	{
		syscallControl.POST("/unlink", api.SyscallUnlink)
		syscallControl.POST("/link", api.SyscallLink)
	}

	files := apiRoutes.Group("/files")
	{
		files.GET("/exists", api.FileExists)            // needs to be a file
		files.GET("/walk", api.WalkFile)                // similar as find in linux command
		files.GET("/list", api.ListFiles)               // list all files and folders
		files.POST("/create", api.CreateFile)           // create file only
		files.POST("/copy", api.CopyFile)               // copy either file or folder
		files.POST("/rename", api.RenameFile)           // rename either file or folder
		files.POST("/move", api.MoveFile)               // move file only
		files.POST("/upload", api.UploadFile)           // upload single file
		files.POST("/download", api.DownloadFile)       // download single file
		files.GET("/read", api.ReadFile)                // read single file
		files.PUT("/write", api.WriteFile)              // write single file
		files.DELETE("/delete", api.DeleteFile)         // delete single file
		files.DELETE("/delete-all", api.DeleteAllFiles) // deletes file or folder
	}

	dirs := apiRoutes.Group("/dirs")
	{
		dirs.GET("/exists", api.DirExists)  // needs to be a folder
		dirs.POST("/create", api.CreateDir) // create folder
	}

	zip := apiRoutes.Group("/zip")
	{
		zip.POST("/unzip", api.Unzip)
		zip.POST("/zip", api.ZipDir)
	}

	user := engine.Group("/api/users", handleUserAuth)
	{
		user.GET("", api.GetUser)
		user.PUT("", api.UpdateUser)
	}

	token := engine.Group("/api/tokens", handleUserAuth)
	{
		token.GET("", api.GetTokens)
		token.GET("/:uuid", api.GetToken)
		token.POST("/generate", api.GenerateToken)
		token.PUT("/:uuid/block", api.BlockToken)
		token.PUT("/:uuid/regenerate", api.RegenerateToken)
		token.DELETE("/:uuid", api.DeleteToken)
	}

	restartJobRoutes := apiRoutes.Group("/restart-jobs")
	{
		restartJobRoutes.GET("", api.GetRestartJob)
		restartJobRoutes.PUT("", api.UpdateRestartJob)
		restartJobRoutes.DELETE("unit/:unit", api.DeleteRestartJob)
	}

	networkingFirewallRoutes := apiRoutes.Group("/firewall")
	{
		networkingFirewallRoutes.GET("", api.UWFStatusList)
		networkingFirewallRoutes.POST("/status", api.UWFStatus)
		networkingFirewallRoutes.POST("/active", api.UWFActive)
		networkingFirewallRoutes.POST("/enable", api.UWFEnable)
		networkingFirewallRoutes.POST("/disable", api.UWFDisable)
		networkingFirewallRoutes.POST("/port/open", api.UWFOpenPort)
		networkingFirewallRoutes.POST("/port/close", api.UWFClosePort)
	}

	networkingRoutes := apiRoutes.Group("/info")
	{

		networkingRoutes.GET("/internet", api.InternetIP)

		networkingNetworkRoutes := networkingRoutes.Group("networks")
		{
			networkingNetworkRoutes.POST("/restart", api.RestartNetworking)
		}

		networkingInterfaceRoutes := networkingRoutes.Group("interfaces")
		{
			networkingInterfaceRoutes.POST("/exists", api.DHCPPortExists)
			networkingInterfaceRoutes.POST("/auto", api.DHCPSetAsAuto)
			networkingInterfaceRoutes.POST("/static", api.DHCPSetStaticIP)
			networkingInterfaceRoutes.POST("/reset", api.InterfaceUpDown)
			networkingInterfaceRoutes.POST("/pp", api.InterfaceUp)
			networkingInterfaceRoutes.POST("/down", api.InterfaceDown)
		}

		networkingFirewallRoutes := networkingRoutes.Group("/firewall")
		{
			networkingFirewallRoutes.GET("", api.UWFStatusList)
			networkingFirewallRoutes.POST("/status", api.UWFStatus)
			networkingFirewallRoutes.POST("/active", api.UWFActive)
			networkingFirewallRoutes.POST("/enable", api.UWFEnable)
			networkingFirewallRoutes.POST("/disable", api.UWFDisable)
			networkingFirewallRoutes.POST("/port/open", api.UWFOpenPort)
			networkingFirewallRoutes.POST("/port/close", api.UWFClosePort)
		}
	}

	storeRoutes := apiRoutes.Group("/store")
	{
		appStoreRoutes := storeRoutes.Group("/apps")
		{
			appStoreRoutes.POST("", api.UploadAddOnAppStore)
			appStoreRoutes.GET("/exists", api.CheckAppExistence)
		}

		pluginStoreRoutes := storeRoutes.Group("/plugins")
		{
			pluginStoreRoutes.GET("", api.GetPluginsStorePlugins)
			pluginStoreRoutes.POST("", api.UploadPluginStorePlugin)
		}

		moduleStoreRoutes := storeRoutes.Group("/modules")
		{
			moduleStoreRoutes.GET("", api.GetModulesStoreModules)
			moduleStoreRoutes.POST("", api.UploadModuleStoreModule)
		}
	}

	apiRoutes.GET("/hosts", api.GetAllInstancesHandler)
	apiRoutes.GET("/hosts/:name", api.GetInstancesHandler)
	apiRoutes.POST("/hosts", api.CreateInstance)
	apiRoutes.GET("/hosts/start", api.StartInstanceHandler)
	apiRoutes.GET("/hosts/stop", api.StopInstanceHandler)
	apiRoutes.GET("/hosts/restart", api.RestartInstanceHandler)
	apiRoutes.DELETE("/hosts/:name", api.DeleteInstanceHandler)

	return engine
}
