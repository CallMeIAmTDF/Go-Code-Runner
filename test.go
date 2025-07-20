//package main
//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/redis/go-redis/v9"
//	"net/http"
//	"os"
//	"os/exec"
//	"path/filepath"
//	"sort"
//	"strings"
//	"sync"
//	"time"
//)
//
//var ctx = context.Background()
//
//var rdb = redis.NewClient(&redis.Options{
//	Addr:     "localhost:6379",
//	Password: "",
//	DB:       0,
//})
//
//type StatusEnum string
//
//const (
//	ACCEPTED              StatusEnum = "ACCEPTED"
//	WRONG_ANSWER          StatusEnum = "WRONG_ANSWER"
//	TIME_LIMIT_EXCEEDED   StatusEnum = "TIME_LIMIT_EXCEEDED"
//	MEMORY_LIMIT_EXCEEDED StatusEnum = "MEMORY_LIMIT_EXCEEDED"
//	COMPILATION_ERROR     StatusEnum = "COMPILATION_ERROR"
//	OUTPUT_LIMIT_EXCEEDED StatusEnum = "OUTPUT_LIMIT_EXCEEDED"
//	PRESENTATION_ERROR    StatusEnum = "PRESENTATION_ERROR"
//	IN_QUEUE              StatusEnum = "IN_QUEUE"
//	CANT_BE_JUDGED        StatusEnum = "CANT_BE_JUDGED"
//)
//
//type TestCase struct {
//	Input          string `json:"input"`
//	ExpectedOutput string `json:"expected_output"`
//}
//
//type RunRequest struct {
//	Language     string     `json:"language"`
//	Code         string     `json:"code"`
//	TimeLimit    int        `json:"time_limit"`
//	SubmissionID string     `json:"submission_id"`
//	MemoryLimit  int        `json:"memory_limit"`
//	ProblemID    int        `json:"problem_id"`
//	TestCases    []TestCase `json:"testcases"`
//}
//
//type TestCaseResult struct {
//	CaseID  int        `json:"case_id"`
//	Status  StatusEnum `json:"status"`
//	Message string     `json:"message"`
//	Time    int64      `json:"time"`
//}
//
//func SaveTestCaseResult(submissionID string, result TestCaseResult) error {
//	key := fmt.Sprintf("submission:%s:testcase:%d", submissionID, result.CaseID)
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return err
//	}
//	return rdb.Set(ctx, key, jsonData, 0).Err()
//}
//
//func GetTestCaseResults(submissionID string) ([]TestCaseResult, error) {
//	pattern := fmt.Sprintf("submission:%s:testcase:*", submissionID)
//	keys, err := rdb.Keys(ctx, pattern).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	var results []TestCaseResult
//	for _, key := range keys {
//		val, err := rdb.Get(ctx, key).Result()
//		if err != nil {
//			continue
//		}
//		var r TestCaseResult
//		if err := json.Unmarshal([]byte(val), &r); err == nil {
//			results = append(results, r)
//		}
//	}
//
//	sort.Slice(results, func(i, j int) bool {
//		return results[i].CaseID < results[j].CaseID
//	})
//
//	return results, nil
//}
//
//func callCallbackAPI(status StatusEnum, total, passed int, message string, maxTime int64, submissionId string) {
//	go func() {
//		body := map[string]interface{}{
//			"status":         status,
//			"totalTestCases": total,
//			"passedTestCase": passed,
//			"message":        message,
//			"maxTime":        maxTime,
//			"id":             submissionId,
//		}
//		jsonBody, _ := json.Marshal(body)
//		req, err := http.NewRequest("POST",
//			"http://localhost:8080/submission/handleCallBackApiHeheheheasd",
//			bytes.NewBuffer(jsonBody),
//		)
//		if err != nil {
//			return
//		}
//		req.Header.Set("Content-Type", "application/json")
//		req.Header.Set("X-Code-Runner-Api-HeHeheHEhe", "req-"+time.Now().Format("20060102150405"))
//		client := &http.Client{Timeout: 3 * time.Second}
//		resp, err := client.Do(req)
//		if err == nil {
//			resp.Body.Close()
//		}
//	}()
//}
//
//func compileJava(code string) (string, StatusEnum, string) {
//	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("java_%d", time.Now().UnixNano()))
//	os.MkdirAll(tempDir, 0755)
//	sourceFile := filepath.Join(tempDir, "Main.java")
//	if err := os.WriteFile(sourceFile, []byte(code), 0644); err != nil {
//		return "", COMPILATION_ERROR, "Failed to write Java file"
//	}
//	compile := exec.Command("javac", sourceFile)
//	compile.Dir = tempDir
//	if err := compile.Run(); err != nil {
//		return "", COMPILATION_ERROR, err.Error()
//	}
//	return tempDir, ACCEPTED, ""
//}
//
//func compileCpp(code string) (string, StatusEnum, string) {
//	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("cpp_%d", time.Now().UnixNano()))
//	os.MkdirAll(tempDir, 0755)
//	sourceFile := filepath.Join(tempDir, "main.cpp")
//	binary := filepath.Join(tempDir, "main.out")
//	os.WriteFile(sourceFile, []byte(code), 0644)
//	compile := exec.Command("g++", sourceFile, "-o", binary)
//	if err := compile.Run(); err != nil {
//		return "", COMPILATION_ERROR, err.Error()
//	}
//	return binary, ACCEPTED, ""
//}
//
//func runCompiledJava(classDir string, input string, timeLimit int) (string, StatusEnum, string, int64) {
//	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeLimit)*time.Second)
//	defer cancel()
//	cmd := exec.CommandContext(ctxTimeout, "java", "-cp", classDir, "Main")
//	stdin, _ := cmd.StdinPipe()
//	stdout := &bytes.Buffer{}
//	cmd.Stdout = stdout
//	cmd.Stderr = stdout
//	stdin.Write([]byte(input))
//	stdin.Close()
//	start := time.Now()
//	err := cmd.Run()
//	elapsed := time.Since(start).Milliseconds()
//	if ctxTimeout.Err() == context.DeadlineExceeded {
//		return "", TIME_LIMIT_EXCEEDED, "Time limit exceeded", elapsed
//	}
//	if err != nil {
//		return "", CANT_BE_JUDGED, err.Error(), elapsed
//	}
//	return strings.TrimSpace(stdout.String()), ACCEPTED, "", elapsed
//}
//
//func runSingleTest(language, code, input string, timeLimit int, compiledPath string) (string, StatusEnum, string, int64) {
//	tempDir := os.TempDir()
//	cmd := &exec.Cmd{}
//	switch language {
//	case "python":
//		sourceFile := filepath.Join(tempDir, fmt.Sprintf("main_%d.py", time.Now().UnixNano()))
//		os.WriteFile(sourceFile, []byte(code), 0644)
//		cmd = exec.Command("python", sourceFile)
//	case "cpp":
//		cmd = exec.Command(compiledPath)
//	case "java":
//		return runCompiledJava(compiledPath, input, timeLimit)
//	default:
//		return "", CANT_BE_JUDGED, "Unsupported language", 0
//	}
//	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeLimit)*time.Second)
//	defer cancel()
//	cmd = exec.CommandContext(ctxTimeout, cmd.Path, cmd.Args[1:]...)
//	stdin, _ := cmd.StdinPipe()
//	stdout := &bytes.Buffer{}
//	cmd.Stdout = stdout
//	cmd.Stderr = stdout
//	stdin.Write([]byte(input))
//	stdin.Close()
//	start := time.Now()
//	err := cmd.Run()
//	elapsed := time.Since(start).Milliseconds()
//	if ctxTimeout.Err() == context.DeadlineExceeded {
//		return "", TIME_LIMIT_EXCEEDED, "Time limit exceeded", elapsed
//	}
//	if err != nil {
//		return "", CANT_BE_JUDGED, err.Error(), elapsed
//	}
//	return strings.TrimSpace(stdout.String()), ACCEPTED, "", elapsed
//}
//
//func validateCode(code, language string) (bool, string) {
//	allowedImports := map[string][]string{
//		"python": {"math", "sys", "collections", "itertools", "functools"},
//		"cpp":    {"iostream", "vector", "string", "algorithm", "map", "set", "queue"},
//		"java":   {"java.util", "java.math"},
//	}
//	notAllowed := []string{}
//	lines := strings.Split(code, "\n")
//	for _, line := range lines {
//		trim := strings.TrimSpace(line)
//		switch language {
//		case "python":
//			if strings.HasPrefix(trim, "import") || strings.HasPrefix(trim, "from") {
//				ok := false
//				for _, lib := range allowedImports["python"] {
//					if strings.Contains(trim, lib) {
//						ok = true
//						break
//					}
//				}
//				if !ok {
//					notAllowed = append(notAllowed, line)
//				}
//			}
//		case "cpp":
//			if strings.HasPrefix(trim, "#include") {
//				ok := false
//				for _, lib := range allowedImports["cpp"] {
//					if strings.Contains(trim, lib) {
//						ok = true
//						break
//					}
//				}
//				if !ok {
//					notAllowed = append(notAllowed, line)
//				}
//			}
//		case "java":
//			if strings.HasPrefix(trim, "import") {
//				ok := false
//				for _, lib := range allowedImports["java"] {
//					if strings.Contains(trim, lib) {
//						ok = true
//						break
//					}
//				}
//				if !ok {
//					notAllowed = append(notAllowed, line)
//				}
//			}
//		}
//	}
//	if len(notAllowed) > 0 {
//		return false, strings.Join(notAllowed, ", ")
//	}
//	return true, ""
//}
//
//func main() {
//	r := gin.Default()
//
//	r.POST("/runcode", func(c *gin.Context) {
//		var req RunRequest
//		if err := c.BindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		if ok, reason := validateCode(req.Code, req.Language); !ok {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid imports", "reason": reason})
//			return
//		}
//
//		total := len(req.TestCases)
//		var passedCount int
//		var maxTime int64
//		compiledPath := ""
//		if req.Language == "java" {
//			path, status, msg := compileJava(req.Code)
//			if status != ACCEPTED {
//				callCallbackAPI(status, total, 0, msg, 0, req.SubmissionID)
//				c.JSON(http.StatusOK, gin.H{"status": status, "message": msg})
//				return
//			}
//			compiledPath = path
//		} else if req.Language == "cpp" {
//			path, status, msg := compileCpp(req.Code)
//			if status != ACCEPTED {
//				callCallbackAPI(status, total, 0, msg, 0, req.SubmissionID)
//				c.JSON(http.StatusOK, gin.H{"status": status, "message": msg})
//				return
//			}
//			compiledPath = path
//		}
//
//		var wg sync.WaitGroup
//		var mu sync.Mutex
//		stop := false
//		results := make([]TestCaseResult, total)
//		for i, t := range req.TestCases {
//			wg.Add(1)
//			go func(i int, t TestCase) {
//				defer wg.Done()
//				if stop {
//					return
//				}
//				output, status, msg, elapsed := runSingleTest(req.Language, req.Code, t.Input, req.TimeLimit, compiledPath)
//				result := TestCaseResult{CaseID: i, Status: status, Message: msg, Time: elapsed}
//				if status == ACCEPTED && output != strings.TrimSpace(t.ExpectedOutput) {
//					result.Status = WRONG_ANSWER
//					result.Message = fmt.Sprintf("Expected: %s, Got: %s", t.ExpectedOutput, output)
//				}
//				mu.Lock()
//				results[i] = result
//				if result.Time > maxTime {
//					maxTime = result.Time
//				}
//				if result.Status == ACCEPTED {
//					passedCount++
//				} else if !stop {
//					stop = true
//					callCallbackAPI(result.Status, total, passedCount, result.Message, maxTime, req.SubmissionID)
//				}
//				SaveTestCaseResult(req.SubmissionID, result)
//				mu.Unlock()
//			}(i, t)
//		}
//		wg.Wait()
//		if !stop {
//			callCallbackAPI(ACCEPTED, total, passedCount, "All testcases passed", maxTime, req.SubmissionID)
//			c.JSON(http.StatusOK, gin.H{"status": ACCEPTED, "passed": passedCount, "total": total})
//		}
//	})
//
//	r.GET("/submission/:id/results", func(c *gin.Context) {
//		submissionID := c.Param("id")
//		results, err := GetTestCaseResults(submissionID)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(http.StatusOK, results)
//	})
//
//	r.Run(":1337")
//}
