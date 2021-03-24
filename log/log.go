/*
@Time : 2021/3/8 5:14 PM
@Author : charley
@File : log
*/
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func Init() {
	Logger = NewLogger(zapcore.ErrorLevel)
}

func InitLogger(isDebug bool) {
	if isDebug {
		Logger = NewLogger(zapcore.DebugLevel)
	} else {
		Logger = NewLogger(zapcore.ErrorLevel)
	}
}
