package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Logger 结构体用于记录日志信息
type Logger struct {
	FilePath string // 日志文件路径
}

// NewLogger 创建一个新的日志记录器
// sourceFilePath 是源代码文件的路径，日志将保存为同名的.md文件
func NewLogger(sourceFilePath string) *Logger {
	ext := filepath.Ext(sourceFilePath)
	base := sourceFilePath[:len(sourceFilePath)-len(ext)]
	logFilePath := base + "r.md"
	return &Logger{FilePath: logFilePath}
}

// Log 记录一条日志信息
func (l *Logger) Log(format string, args ...interface{}) error {
	// 确保日志文件存在
	file, err := os.OpenFile(l.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 添加时间戳
	// 添加时间戳
	logEntry := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))

	// 写入日志
	_, err = file.WriteString(logEntry)
	return err
}

// LogFunctionEntry 记录函数开始执行
func (l *Logger) LogFunctionEntry(funcName, description string) error {
	l.Log("--- 分割线 ---")
	return l.Log("## 函数: %s\n开始执行: %s", funcName, description)
}

// LogFunctionExit 记录函数执行结束
func (l *Logger) LogFunctionExit(funcName string, result interface{}) error {
	l.Log("--- 分割线 ---")
	return l.Log("函数 %s 执行完成，结果: %v", funcName, result)
}

// LogError 记录错误信息
func (l *Logger) LogError(funcName string, err error) error {
	return l.Log("函数 %s 执行出错: %v", funcName, err)
}

// LogDetail 记录详细信息
func (l *Logger) LogDetail(detail string) error {
	return l.Log("%s", detail)
}
