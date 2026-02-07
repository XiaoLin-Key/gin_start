package logger

import (
	"fmt"
	"gin_start/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Lg *zap.Logger
	S  *zap.SugaredLogger
)

func InitLogger() (err error) {
	writeSyncer := getLogWriter()
	encoder := getEncoder()

	levelStr := viper.GetString("log.level")
	var level zapcore.Level

	//  将字符串反序列化为 zapcore.Level 类型
	// 如果解析失败（比如填错了），默认使用 DebugLevel
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		level = zapcore.DebugLevel
	}

	// 使用解析出的级别创建 Core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	Lg = zap.New(core, zap.AddCaller())
	S = Lg.Sugar()
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	now := time.Now().Format("2006-01-02")

	conf := settings.Conf.Log
	if conf.Filename == "" {
		conf.Filename = "./logs" // 保底路径
	}
	fileName := fmt.Sprintf("%s/%s.log", conf.Filename, now)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackups,
		LocalTime:  true,
	}
	fileWriter := zapcore.AddSync(lumberJackLogger)

	if settings.Conf.Mode == "dev" {
		// 如果是开发模式，使用 MultiWriteSyncer 同时输出到控制台和文件
		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			fileWriter,
		)
	}

	// 3. 其他模式（如 prod）只返回文件输出流
	return zapcore.NewMultiWriteSyncer(fileWriter)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 执行后续处理逻辑
		c.Next()

		cost := time.Since(start)

		// 统一使用你昨天配置好的 Lg (zap.Logger)
		// 这样它会自动触发你写的 getLogWriter (按天生成文件) 和 getEncoder (控制台彩色)
		Lg.Info(path,
			zap.Int("状态码", c.Writer.Status()),
			zap.String("方法", c.Request.Method),
			zap.String("路径", path),
			zap.String("查询参数", query),
			zap.String("客户端IP", c.ClientIP()),
			zap.String("用户代理", c.Request.UserAgent()),
			zap.String("错误信息", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("耗时", cost),
		)
	}
}

// GinRecovery recover 项目可能出现的 panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否是由于连接断开引起的 panic
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求原始数据用于记录
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					// 这种情况下，你昨天的红色 Error 配置会非常显眼
					Lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					Lg.Error("[Panic 崩溃恢复]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())), // 打印完整堆栈
					)
				} else {
					Lg.Error("[Panic 崩溃恢复]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 发生 panic 时返回 500 状态码
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
