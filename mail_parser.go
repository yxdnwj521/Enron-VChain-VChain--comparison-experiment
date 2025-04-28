package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// 获取当前文件的路径
func getCurrentFilePath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filename
}

type Mail struct {
	MessageID               string
	Date                    string
	From                    string
	To                      string
	Subject                 string
	MimeVersion             string
	ContentType             string
	ContentTransferEncoding string
	XFrom                   string
	XTo                     string
	XCc                     string
	XBcc                    string
	XFolder                 string
	XOrigin                 string
	XFileName               string
	Body                    string
	FileName                string
	FileSize                int64
	Keywords                []string
	ManualKeywords          []string // 新增字段用于手动输入关键词
}

var (
	fieldPattern  = regexp.MustCompile(`^(Message-ID|Date|From|To|Subject|Mime-Version|Content-Type|Content-Transfer-Encoding|X-From|X-To|X-cc|X-bcc|X-Folder|X-Origin|X-FileName):[ \t]*(.*)$`)
	keywordList   = []string{"California", "Summary", "Update", "Load", "Report", "Meeting", "Analysis", "Graph", "Chart"}
	wordFrequency map[string]int
)

func ParseMailFile(path string) (*Mail, error) {
	// 创建日志记录器
	logger := NewLogger(getCurrentFilePath())
	logger.LogFunctionEntry("ParseMailFile", "解析邮件文件: "+path)

	file, err := os.Open(path)
	if err != nil {
		logger.LogError("ParseMailFile", err)
		return nil, err
	}
	defer file.Close()
	info, _ := file.Stat()
	mail := &Mail{FileName: filepath.Base(path), FileSize: info.Size()}
	var bodyLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if m := fieldPattern.FindStringSubmatch(line); m != nil {
			switch m[1] {
			case "Message-ID":
				mail.MessageID = m[2]
			case "Date":
				mail.Date = m[2]
			case "From":
				mail.From = m[2]
			case "To":
				mail.To = m[2]
			case "Subject":
				mail.Subject = m[2]
			case "Mime-Version":
				mail.MimeVersion = m[2]
			case "Content-Type":
				mail.ContentType = m[2]
			case "Content-Transfer-Encoding":
				mail.ContentTransferEncoding = m[2]
			case "X-From":
				mail.XFrom = m[2]
			case "X-To":
				mail.XTo = m[2]
			case "X-cc":
				mail.XCc = m[2]
			case "X-bcc":
				mail.XBcc = m[2]
			case "X-Folder":
				mail.XFolder = m[2]
			case "X-Origin":
				mail.XOrigin = m[2]
			case "X-FileName":
				mail.XFileName = m[2]
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}
	mail.Body = strings.Join(bodyLines, "\n")
	wordFrequency = calculateWordFrequency(mail.Body)
	mail.Keywords = ExtractKeywords(mail.Body)

	// 自动提取高频词作为关键词
	threshold := 5 // 设定词频阈值
	for word, freq := range wordFrequency {
		if freq >= threshold {
			mail.Keywords = append(mail.Keywords, word)
		}
	}

	// 新增：自动从manual_keywords.txt读取手动关键词
	manualFile := "manual_keywords.txt"
	if _, err := os.Stat(manualFile); err == nil {
		f, err := os.Open(manualFile)
		if err == nil {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" && !strings.HasPrefix(line, "#") {
					// 支持一行多个关键词，逗号、空格分隔
					for _, kw := range strings.FieldsFunc(line, func(r rune) bool {
						return r == ',' || r == ' ' || r == '\t'
					}) {
						kw = strings.TrimSpace(kw)
						if kw != "" {
							mail.ManualKeywords = append(mail.ManualKeywords, kw)
						}
					}
				}
			}
			f.Close()
		}
	}
	if len(mail.ManualKeywords) > 0 {
		mail.Keywords = append(mail.Keywords, mail.ManualKeywords...)
	}
	return mail, nil
}

func ExtractKeywords(body string) []string {
	// 创建日志记录器
	logger := NewLogger(getCurrentFilePath())
	logger.LogFunctionEntry("ExtractKeywords", "从邮件正文中提取关键词")

	var found []string
	for _, kw := range keywordList {
		if strings.Contains(strings.ToLower(body), strings.ToLower(kw)) {
			found = append(found, kw)
			logger.LogDetail(fmt.Sprintf("自动提取关键词: %v", found))
		}
	}

	logger.LogFunctionExit("ExtractKeywords", found)
	return found
}

func RunMailParser() {
	// 创建日志记录器
	logger := NewLogger(getCurrentFilePath())
	logger.LogFunctionEntry("RunMailParser", "开始解析邮件文件")

	root := "./_sent_mail"
	logger.LogDetail(fmt.Sprintf("扫描目录: %s", root))

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		mail, err := ParseMailFile(path)
		if err != nil {
			logger.LogError("RunMailParser", fmt.Errorf("解析失败: %s", path))
			fmt.Printf("解析失败: %s\n", path)
			return nil
		}

		// 记录解析结果
		logger.LogDetail(fmt.Sprintf("成功解析文件: %s", mail.FileName))
		logger.LogDetail(fmt.Sprintf("文件大小: %d 字节", mail.FileSize))
		logger.LogDetail(fmt.Sprintf("主题: %s", mail.Subject))
		logger.LogDetail(fmt.Sprintf("最终关键词列表: %v", mail.Keywords))

		fmt.Printf("文件: %s\n大小: %d 字节\n主题: %s\n关键词: %v\n\n", mail.FileName, mail.FileSize, mail.Subject, mail.Keywords)
		return nil
	})
}

func calculateWordFrequency(body string) map[string]int {
	wordFrequency := make(map[string]int)
	words := strings.Fields(body)
	for _, word := range words {
		word = strings.ToLower(word)
		wordFrequency[word]++
	}
	return wordFrequency
}
