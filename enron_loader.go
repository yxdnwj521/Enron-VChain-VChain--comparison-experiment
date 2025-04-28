package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

// 获取当前文件的路径
func getCurrentFilePathEL() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filename
}

// 加载Enron数据集的辅助函数，后续可扩展为解析邮件内容
func LoadEnronDataset(path string) ([][]byte, error) {
	// 创建日志记录器
	logger := NewLogger(getCurrentFilePathEL())
	logger.LogFunctionEntry("LoadEnronDataset", fmt.Sprintf("加载数据集: %s", path))

	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.LogError("LoadEnronDataset", err)
		return nil, err
	}

	// 检查空文件夹
	if len(files) == 0 {
		logger.LogDetail("警告: 数据集文件夹为空")
		return [][]byte{}, nil
	}

	var dataset [][]byte
	logger.LogDetail(fmt.Sprintf("找到%d个文件", len(files)))

	for _, file := range files {
		if !file.IsDir() {
			logger.LogDetail(fmt.Sprintf("读取文件: %s", file.Name()))
			data, err := os.ReadFile(path + "/" + file.Name())
			if err != nil {
				logger.LogError("LoadEnronDataset", fmt.Errorf("读取文件失败: %s", file.Name()))
				fmt.Println("读取文件失败:", file.Name())
				continue
			}
			dataset = append(dataset, data)
		}
	}

	logger.LogFunctionExit("LoadEnronDataset", fmt.Sprintf("成功加载%d个文件", len(dataset)))
	return dataset, nil
}

// 生成数据集统计报告
func GenerateDatasetReport(dataset [][]byte) string {
	if len(dataset) == 0 {
		return "数据集为空，无统计信息"
	}

	totalSize := 0
	for _, data := range dataset {
		totalSize += len(data)
	}
	avgSize := totalSize / len(dataset)

	return fmt.Sprintf("数据集统计:\n- 文件数量: %d\n- 总大小: %d bytes\n- 平均文件大小: %d bytes",
		len(dataset), totalSize, avgSize)
}
