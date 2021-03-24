/*
@Time : 2021/1/18 5:24 PM
@Author : charley
@File : logCore.go
*/
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(level zapcore.Level) *zap.SugaredLogger {
	core := newCore(level)
	//logger := zap.New(core, zap.AddCaller(), zap.Development()).Sugar()
	logger := zap.New(core, zap.Development()).Sugar()
	return logger
}

func newCore(level zapcore.Level) zapcore.Core {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

}
