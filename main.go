package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// 获取当前文件的路径
func getCallerFilePath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filename
}

// 主函数，用于运行所有实验
func main() {
	// 创建日志记录器
	mainLogger := NewLogger(getCallerFilePath())
	mainLogger.Log("# 实验开始运行")
	mainLogger.Log("## 主程序初始化")

	// 获取所有Go文件的路径
	files, err := filepath.Glob(filepath.Join(filepath.Dir(getCallerFilePath()), "*.go"))
	if err != nil {
		mainLogger.LogError("main", err)
		fmt.Println("获取文件列表失败:", err)
		return
	}

	mainLogger.Log("找到以下Go文件:")
	for _, file := range files {
		mainLogger.Log("- %s", filepath.Base(file))
	}

	// 运行邮件解析器
	mainLogger.Log("## 开始运行邮件解析器")
	RunMailParser()

	// 运行实验
	mainLogger.Log("## 开始运行实验")
	// 这里调用的是experiment.go中的实验函数，而不是main函数
	RunExperiment()

	// 打印实验结果
	mainLogger.Log("## 打印实验结果")
	err = PrintResults("experiment_result.csv")
	if err != nil {
		mainLogger.LogError("PrintResults", err)
		fmt.Println("打印结果失败:", err)
	}

	mainLogger.Log("# 实验结束")
	fmt.Println("所有操作已完成，日志已保存到相应的.md文件中")
}

// RunExperiment 运行实验的函数，替代experiment.go中的main函数
func RunExperiment() {
	// 获取当前文件路径
	experimentPath, err := filepath.Abs("experiment.go")
	if err != nil {
		fmt.Println("获取experiment.go路径失败:", err)
		return
	}

	// 创建日志记录器
	logger := NewLogger(experimentPath)
	logger.LogFunctionEntry("RunExperiment", "开始运行vchain/vchain+对比实验")

	fmt.Println("Enron数据集vchain/vchain+对比实验")
	logger.LogDetail("加载邮件数据集")

	// 加载邮件数据集
	dataset, err := LoadEnronDataset("./_sent_mail")
	if err != nil {
		logger.LogError("LoadEnronDataset", err)
		fmt.Println("数据集加载失败:", err)
		return
	}
	logger.LogDetail(fmt.Sprintf("成功加载数据集，共%d个文件", len(dataset)))

	// 构建索引
	logger.LogDetail("构建索引结构")
	vchain := NewVChain()
	vchainplus := NewVChainPlus()
	logger.LogDetail("索引结构构建完成")

	// 关键词查询列表
	keywordsList := [][]string{
		{"California"},
		{"Summary"},
		{"Update"},
		{"Load"},
		{"Report"},
		{"Meeting"},
		{"Analysis"},
		{"Graph"},
		{"Chart"},
		{"California", "Summary"},
		{"Update", "Load"},
		{"Report", "Meeting", "Analysis"},
	}
	logger.LogDetail(fmt.Sprintf("准备测试%d组关键词", len(keywordsList)))

	// 记录实验结果
	var records [][]string
	records = append(records, []string{"方法", "关键词数量", "文件大小(字节)", "验证时间(毫秒)", "VO大小(字节)", "Gas消耗"})

	fileInfo, ferr := os.Stat("./_sent_mail")
	fileSize := int64(0)
	if ferr == nil {
		fileSize = fileInfo.Size()
	}
	logger.LogDetail(fmt.Sprintf("数据集大小: %d字节", fileSize))

	for i, keywords := range keywordsList {
		logger.LogDetail(fmt.Sprintf("测试第%d组关键词: %v", i+1, keywords))

		// vchain测试
		logger.LogDetail("执行vchain查询")
		vo, verifyTime, gasUsed := vchain.Query(keywords)
		logger.LogDetail(fmt.Sprintf("vchain查询完成: 验证时间=%v, VO大小=%d, Gas消耗=%d",
			verifyTime, len(vo), gasUsed))
		records = append(records, []string{"vchain", fmt.Sprintf("%d", len(keywords)),
			fmt.Sprintf("%d", fileSize), fmt.Sprintf("%d", verifyTime.Milliseconds()),
			fmt.Sprintf("%d", len(vo)), fmt.Sprintf("%d", gasUsed)})

		// vchain+测试
		logger.LogDetail("执行vchain+查询")
		vo2, verifyTime2, gasUsed2 := vchainplus.Query(keywords)
		logger.LogDetail(fmt.Sprintf("vchain+查询完成: 验证时间=%v, VO大小=%d, Gas消耗=%d",
			verifyTime2, len(vo2), gasUsed2))
		records = append(records, []string{"vchain+", fmt.Sprintf("%d", len(keywords)),
			fmt.Sprintf("%d", fileSize), fmt.Sprintf("%d", verifyTime2.Milliseconds()),
			fmt.Sprintf("%d", len(vo2)), fmt.Sprintf("%d", gasUsed2)})
	}

	// 保存实验结果
	resultFile := "experiment_result.csv"
	logger.LogDetail(fmt.Sprintf("保存实验结果到%s", resultFile))

	err = SaveResults(resultFile, records)
	if err != nil {
		logger.LogError("SaveResults", err)
		fmt.Println("结果文件创建失败:", err)
		return
	}

	logger.LogFunctionExit("RunExperiment", "实验完成，结果已保存")
	fmt.Println("实验完成，结果已保存到", resultFile)
}
