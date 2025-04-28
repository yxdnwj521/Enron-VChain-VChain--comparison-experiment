package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
)

// 获取当前文件的路径
func getCurrentFilePathRA() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filename
}

// 保存实验结果到CSV文件，便于后续分析
func SaveResults(filename string, records [][]string) error {
	// 创建日志记录器
	logger := NewLogger(getCurrentFilePathRA())
	logger.LogFunctionEntry("SaveResults", fmt.Sprintf("保存实验结果到文件: %s", filename))

	file, err := os.Create(filename)
	if err != nil {
		logger.LogError("SaveResults", err)
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// 加载并打印CSV实验结果
func PrintResults(filename string) error {
	// 创建日志记录器
	logger := NewLogger(getCurrentFilePathRA())
	logger.LogFunctionEntry("PrintResults", fmt.Sprintf("加载并打印CSV文件: %s", filename))

	file, err := os.Open(filename)
	if err != nil {
		logger.LogError("PrintResults", err)
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for i, record := range records {
		fmt.Println(record)
		logger.LogDetail(fmt.Sprintf("记录 %d: %v", i, record))
	}

	logger.LogFunctionExit("PrintResults", "成功打印所有记录")
	return nil
}
