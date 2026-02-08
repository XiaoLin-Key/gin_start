package main

import (
	"context"
	"fmt"
	"gin_start/dao/mysql"
	"gin_start/dao/redis"
	"gin_start/logger"
	"gin_start/pkg/snowflake"
	"gin_start/pkg/validator"
	"gin_start/routes"
	"gin_start/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title           Gin-Start 项目接口文档
// @version         1.0
// @description     这是一个基于 Gin 框架开发的社区后端项目。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9090
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description 输入格式: Bearer <your_token>  (注意中间有空格)
func main() {
	// 1. 初始化各组件 (Viper, Logger, MySQL, Redis,Snowflake)
	if err := initAll(); err != nil {
		panic(err)
	}
	defer logger.Lg.Sync()

	// 2. 注册路由 (在这里解决你说的无关日志)
	r := routes.SetUpRouter(settings.Conf.Mode)

	// 3. 开启服务（优雅启动与关机）
	if err := runServer(r); err != nil {
		logger.Lg.Fatal("服务异常退出", zap.Error(err))
	}
}

// 进一步封装初始化逻辑，让 main 更短
func initAll() error {
	if err := settings.InitViper(); err != nil {
		return fmt.Errorf("init viper failed: %w", err)
	}
	if err := logger.InitLogger(); err != nil {
		return fmt.Errorf("init logger failed: %w", err)
	}
	if err := mysql.InitMySQL(); err != nil {
		return fmt.Errorf("init mysql failed: %w", err)
	}
	if err := redis.InitRedis(); err != nil {
		return fmt.Errorf("init redis failed: %w", err)
	}
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		return fmt.Errorf("init snowflake failed: %w", err)
	}
	if err := validator.InitTrans("zh"); err != nil {
		return fmt.Errorf("init validator trans failed: %w", err)
	}
	return nil
}

func runServer(r *gin.Engine) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Lg.Error("监听启动失败", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关机流程
	fmt.Print("\r正在关闭服务...") // \r 可以覆盖掉 Ctrl+C 产生的 ^C 字符，更美观

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP服务关闭异常: %v", err)
	}

	closeResources()
	fmt.Println(" [完成]") // 最终输出：正在关闭服务... [完成]
	return nil
}

func closeResources() {
	if sqlDB, err := mysql.GetDB().DB(); err == nil {
		_ = sqlDB.Close()
	}
	_ = redis.Rdb.Close()
}
