# Enron VChain/VChain+ 对比实验

## 项目简介

本项目用于在Enron电子邮件数据集上对比vchain与vchain+两种区块链索引结构的性能，包括gas消耗、文件大小、数据验证时间、VO大小等指标。

## 文件结构与功能

### 主要功能文件

1. **experiment.go** - 主实验流程
   - `RunExperiment()`: 执行完整的对比实验流程
   - `comparePerformance()`: 比较两种索引结构的性能指标
   - `generateReport()`: 生成实验结果报告

2. **enron_loader.go** - 数据集加载
   - `LoadEnronDataset(path string)`: 加载指定路径下的Enron邮件数据集
   - `getCurrentFilePathEL()`: 获取当前文件路径用于日志记录

3. **logger.go** - 日志记录
   - `NewLogger()`: 创建新的日志记录器实例
   - `LogFunctionEntry/Exit()`: 记录函数进入/退出
   - `LogError()`: 记录错误信息
   - `LogDetail()`: 记录详细过程信息

4. **mail_parser.go** - 邮件解析
   - `ParseKeywords()`: 解析manual_keywords.txt中的关键词
   - `ProcessMailContent()`: 处理邮件内容提取关键词

5. **result_analyzer.go** - 结果分析
   - `SaveResults()`: 保存实验结果到CSV
   - `AnalyzeResults()`: 分析实验结果数据

## 使用方法

### 安装指南

1. **系统要求**
   - Go 1.18+ 版本
   - 至少4GB可用内存
   - 10GB可用磁盘空间

2. **安装依赖**
   ```bash
   # 安装Go依赖
   go mod tidy
   
   # 可选: 安装测试依赖
   go get -t ./...
   ```

3. **依赖项列表**
   - github.com/stretchr/testify (测试)
   - github.com/gin-gonic/gin (Web框架)
   - github.com/spf13/viper (配置管理)

### 配置实验

1. 将Enron数据集放在`./data`目录下
2. 在`manual_keywords.txt`中添加关键词(每行一个或多个，用逗号/空格/Tab分隔)

### 运行实验
```bash
go run .
```

### 查看结果
- 控制台输出
- `experiment_result.csv`文件
- 日志文件

## 贡献指南

1. **报告问题**
   - 在GitHub Issues中报告问题
   - 提供重现步骤和预期/实际结果

2. **提交PR**
   - Fork仓库并创建新分支
   - 遵循现有代码风格
   - 添加适当的测试用例
   - 更新相关文档

3. **联系方式**
   - 邮箱: 1917413192@qq.com
   - GitHub Discussions: 用于讨论新功能


## 自定义代码

1. **修改实验参数**
   - 在`experiment.go`中调整`experimentConfig`结构体参数

2. **添加新指标**
   - 在`result_analyzer.go`中添加新的分析函数
   - 更新`SaveResults()`函数以包含新指标

3. **扩展关键词处理**
   - 修改`mail_parser.go`中的`ParseKeywords()`函数
   - 添加新的关键词处理逻辑

## 实验指标说明

| 指标 | 说明 |
|------|------|
| Gas消耗 | 区块链操作的计算成本 |
| 文件大小 | 处理的数据文件大小 |
| 验证时间 | 数据验证所需时间 |
| VO大小 | 验证对象(Verification Object)的大小 |

## 注意事项

- 确保有足够的磁盘空间存储实验结果
- 实验过程中不要修改输入文件
- 删除旧实验结果: `rm *r.md`
