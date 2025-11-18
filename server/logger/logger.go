package logger

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"portfolio/config"
	"portfolio/helpers"
	"sort"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	*log.Logger
	logPath     string
	cfg         *config.LoggingConfig
	file        *os.File
	environment string
	currentSize float32
	maxSize     float32
	maxBackups  int
	maxAge      int
	compress    bool
	mu          sync.Mutex
	currentDate time.Time
	rotateDaily bool
	dailyTimer  *time.Timer
}

type LogConfig struct {
	File        string  `yaml:"file"`
	Level       string  `yaml:"level"`
	MaxSize     float32 `yaml:"max_size"`
	MaxBackups  int     `yaml:"max_backups"`
	MaxAge      int     `yaml:"max_age"`
	Compress    bool    `yaml:"compress"`
	RotateDaily bool    `yaml:"rotate_daily"`
}

func NewLogger(cfg *config.LoggingConfig, environment string) *Logger {

	err := helpers.MkdirIfNotExists(cfg.File)

	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logConfig := getLogConfig(cfg)

	logger := &Logger{
		cfg:         cfg,
		maxSize:     logConfig.MaxSize * 1024 * 1024,
		maxBackups:  logConfig.MaxBackups,
		maxAge:      logConfig.MaxAge,
		compress:    logConfig.Compress,
		logPath:     logConfig.File,
		environment: environment,
		currentDate: time.Now().Truncate(24 * time.Hour),
		rotateDaily: logConfig.RotateDaily,
	}

	if err := logger.openLogFile(); err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	go logger.cleanOldLogs()

	if logger.rotateDaily {
		go logger.startDailyRotationTimer()
	}

	return logger
}

func getLogConfig(cfg *config.LoggingConfig) LogConfig {

	compress := true
	if cfg.Compress != nil {
		compress = *cfg.Compress
	}

	rotateDaily := false
	if cfg.RotateDaily != nil {
		rotateDaily = *cfg.RotateDaily
	}

	return LogConfig{
		File:        getStringOrDefault(cfg.File, "logs/portfolio.log"),
		MaxSize:     getFloat32OrDefault(cfg.MaxSize, 10.0),
		MaxBackups:  getIntOrDefault(cfg.MaxBackups, 7),
		MaxAge:      getIntOrDefault(cfg.MaxAge, 30),
		Compress:    compress,
		Level:       getStringOrDefault(cfg.Level, "info"),
		RotateDaily: rotateDaily,
	}
}

func (l *Logger) openLogFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.openLogFileUnsafe()
}

func (l *Logger) openLogFileUnsafe() error {

	if l.file != nil {
		if err := l.file.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
	}

	file, err := os.OpenFile(l.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.file = file

	if info, err := file.Stat(); err == nil {
		l.currentSize = float32(info.Size())

		if info.Size() > 0 {
			l.currentDate = info.ModTime().Truncate(24 * time.Hour)
		} else {
			l.currentDate = time.Now().Truncate(24 * time.Hour)
		}
	}

	var output io.Writer
	if l.environment == "development" {
		output = io.MultiWriter(file, os.Stdout)
	} else {
		output = file
	}

	l.Logger = log.New(output, "", log.LstdFlags|log.Lshortfile)

	return nil
}

func (l *Logger) write(level string, format string, v ...interface{}) {
	message := fmt.Sprintf("[%s] %s", level, fmt.Sprintf(format, v...))

	if l.needsRotation(len(message)) {
		if err := l.rotate(); err != nil {
			log.Printf("log rotate: %v", err)
		}
	}

	l.mu.Lock()
	l.Print(message)
	l.currentSize += float32(len(message))
	l.mu.Unlock()
}

func (l *Logger) needsRotation(messageSize int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.needsRotationUnsafe(messageSize)
}

func (l *Logger) rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		if err := l.file.Sync(); err != nil {
			log.Printf("Error syncing log file: %v", err)
		}
		if err := l.file.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
		l.file = nil
	}

	timestamp := time.Now().Format("2006-01-02")

	if l.rotateDaily {
		if contentTimestamp := l.analyzeFileContentForTimestamp(); contentTimestamp != "" {
			timestamp = contentTimestamp
		}
	}

	backupPath := l.getBackupPath(timestamp)

	if err := os.Rename(l.logPath, backupPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	if l.compress {
		go l.compressFile(backupPath)
	}

	l.currentSize = 0
	if l.rotateDaily {
		l.currentDate = time.Now().Truncate(24 * time.Hour)
	}

	if err := l.openLogFileUnsafe(); err != nil {
		return fmt.Errorf("failed to reopen log file after rotation: %w", err)
	}

	go l.cleanOldLogs()

	return nil
}

func (l *Logger) getBackupPath(timestamp string) string {
	dir := filepath.Dir(l.logPath)
	filename := filepath.Base(l.logPath)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	maxAttempts := l.maxBackups
	if maxAttempts <= 0 {
		maxAttempts = 10
	}

	baseName := fmt.Sprintf("%s-%s", name, timestamp)
	proposedPath := filepath.Join(dir, baseName+ext)

	if !l.fileExists(proposedPath) {
		return proposedPath
	}

	for i := 1; i <= maxAttempts; i++ {
		numberedName := fmt.Sprintf("%s-%s-%03d", name, timestamp, i)
		numberedPath := filepath.Join(dir, numberedName+ext)
		if !l.fileExists(numberedPath) {
			return numberedPath
		}
	}

	timeNow := time.Now()
	preciseTimestamp := timeNow.Format("2006-01-02-150405")
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", name, preciseTimestamp, ext))
}

func (l *Logger) fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}

	if l.compress {
		if _, err := os.Stat(filePath + ".gz"); err == nil {
			return true
		}
	}

	return false
}

func (l *Logger) analyzeFileContentForTimestamp() string {
	if _, err := os.Stat(l.logPath); err != nil {
		return ""
	}

	file, err := os.Open(l.logPath)
	if err != nil {
		return ""
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file during timestamp analysis: %v", err)
		}
	}()

	var firstTimestamp, lastTimestamp time.Time
	var found bool

	scanner := bufio.NewScanner(file)
	lineCount := 0
	maxLinesToCheck := 1000

	for scanner.Scan() && lineCount < maxLinesToCheck {
		line := scanner.Text()
		if timestamp := extractTimestampFromLogLine(line); !timestamp.IsZero() {
			if !found {
				firstTimestamp = timestamp
				found = true
			}
			lastTimestamp = timestamp
		}
		lineCount++
	}

	if !found {
		return ""
	}

	firstDate := firstTimestamp.Truncate(24 * time.Hour)
	lastDate := lastTimestamp.Truncate(24 * time.Hour)

	if firstDate.Equal(lastDate) {
		return firstDate.Format("2006-01-02")
	}

	return lastDate.Format("2006-01-02")
}

func extractTimestampFromLogLine(line string) time.Time {
	patterns := []string{
		"2006/01/02 15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"Jan 02 15:04:05",
	}

	for _, pattern := range patterns {
		if timestamp, err := time.Parse(pattern, extractTimeString(line, pattern)); err == nil {
			if pattern == "Jan 02 15:04:05" {
				year := time.Now().Year()
				timestamp = timestamp.AddDate(year-timestamp.Year(), 0, 0)
			}
			return timestamp
		}
	}

	return time.Time{}
}

func extractTimeString(line, pattern string) string {
	if len(line) >= len(pattern) {
		candidate := line[:len(pattern)]
		return candidate
	}
	return line
}

func (l *Logger) compressFile(filePath string) {
	time.Sleep(100 * time.Millisecond)

	input, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file for compression %s: %v", filePath, err)
		return
	}
	defer func() {
		if err := input.Close(); err != nil {
			log.Printf("input close: %v", err)
		}
	}()

	info, err := input.Stat()
	if err != nil {
		log.Printf("Error getting file info %s: %v", filePath, err)
		return
	}
	if info.Size() == 0 {
		log.Printf("Skipping compression of empty file: %s", filePath)
		return
	}

	output, err := os.Create(filePath + ".gz")
	if err != nil {
		log.Printf("Error creating compressed file %s.gz: %v", filePath, err)
		return
	}

	// defer func() {
	// 	if err := output.Close(); err != nil {
	// 		log.Printf("output close: %v", err)
	// 	}
	// }()

	gzWriter := gzip.NewWriter(output)
	defer func() {
		if err := gzWriter.Close(); err != nil {
			log.Printf("gzip close: %v", err)
		}
	}()

	_, err = io.Copy(gzWriter, input)
	if err != nil {
		log.Printf("Error compressing file %s: %v", filePath, err)
		return
	}

	if err := gzWriter.Close(); err != nil {
		log.Printf("Error finalizing compression %s: %v", filePath, err)
		return
	}

	if err := output.Close(); err != nil {
		log.Printf("Error closing compressed file %s: %v", filePath, err)
		return
	}

	if err := os.Remove(filePath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Printf("remove %s: %v", filePath, err)
	}
}

func (l *Logger) cleanOldLogs() {
	dir := filepath.Dir(l.logPath)
	filename := filepath.Base(l.logPath)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	pattern := filepath.Join(dir, name+"-*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("Error finding old log files: %v", err)
		return
	}

	var logFiles []string
	for _, match := range matches {
		if strings.HasSuffix(match, ".gz") {
			if info, err := os.Stat(match); err == nil && info.Size() == 0 {
				continue
			}
		}
		logFiles = append(logFiles, match)
	}

	sort.Slice(logFiles, func(i, j int) bool {
		info1, err1 := os.Stat(logFiles[i])
		info2, err2 := os.Stat(logFiles[j])
		if err1 != nil || err2 != nil {
			return false
		}
		return info1.ModTime().After(info2.ModTime())
	})

	if l.maxBackups > 0 && len(logFiles) > l.maxBackups {
		for _, file := range logFiles[l.maxBackups:] {
			if err := os.Remove(file); err != nil {
				log.Printf("Error removing old log file %s: %v", file, err)
			} else {
				log.Printf("Removed old log file due to maxBackups limit: %s", filepath.Base(file))
			}
		}
		logFiles = logFiles[:l.maxBackups]
	}

	if l.maxAge > 0 {
		cutoff := time.Now().AddDate(0, 0, -l.maxAge)
		var filesToKeep []string
		for _, file := range logFiles {
			if info, err := os.Stat(file); err == nil {
				if info.ModTime().Before(cutoff) {
					if err := os.Remove(file); err != nil {
						log.Printf("Error removing old log file %s: %v", file, err)
					} else {
						log.Printf("Removed old log file due to maxAge limit: %s", filepath.Base(file))
					}
				} else {
					filesToKeep = append(filesToKeep, file)
				}
			}
		}
		logFiles = filesToKeep
	}

	if l.cfg.Level == "debug" {
		log.Printf("Log cleanup completed. Kept %d backup files", len(logFiles))
	}
}

func (l *Logger) GetBackupFiles() []map[string]interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	dir := filepath.Dir(l.logPath)
	filename := filepath.Base(l.logPath)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	pattern := filepath.Join(dir, name+"-*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}

	var files []map[string]interface{}
	for _, match := range matches {
		if info, err := os.Stat(match); err == nil {
			files = append(files, map[string]interface{}{
				"name":      filepath.Base(match),
				"path":      match,
				"size":      info.Size(),
				"mod_time":  info.ModTime().Format("2006-01-02 15:04:05"),
				"age_days":  int(time.Since(info.ModTime()).Hours() / 24),
				"will_keep": l.shouldKeepFile(info),
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		time1, _ := time.Parse("2006-01-02 15:04:05", files[i]["mod_time"].(string))
		time2, _ := time.Parse("2006-01-02 15:04:05", files[j]["mod_time"].(string))
		return time1.After(time2)
	})

	return files
}

func (l *Logger) shouldKeepFile(info os.FileInfo) bool {
	if l.maxAge > 0 {
		cutoff := time.Now().AddDate(0, 0, -l.maxAge)
		if info.ModTime().Before(cutoff) {
			return false
		}
	}
	return true
}

func (l *Logger) startDailyRotationTimer() {
	go func() {
		for {
			now := time.Now()
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			duration := nextMidnight.Sub(now)

			l.mu.Lock()
			if l.dailyTimer != nil {
				l.dailyTimer.Stop()
			}
			l.dailyTimer = time.NewTimer(duration)
			timer := l.dailyTimer
			l.mu.Unlock()

			<-timer.C

			l.mu.Lock()
			if l.dailyTimer == nil {
				l.mu.Unlock()
				return
			}
			l.mu.Unlock()

			if err := l.rotate(); err != nil {
				log.Printf("Error rotating log file: %v", err)
			}
		}
	}()
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.write("INFO", format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.write("ERROR", format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.write("FATAL", format, v...)
	os.Exit(1)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.write("DEBUG", format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.write("WARN", format, v...)
}

func (l *Logger) HTTP(method, path string, statusCode int, duration, clientIP string) {
	l.write("HTTP", "%s %s %d %s %s", method, path, statusCode, duration, clientIP)
}

func (l *Logger) Rotate() error {
	return l.rotate()
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.dailyTimer != nil {
		l.dailyTimer.Stop()
		l.dailyTimer = nil
	}

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) GetRotationInfo() map[string]interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	return map[string]interface{}{
		"max_size_mb":  l.maxSize / (1024 * 1024),
		"current_size": l.currentSize,
		"max_backups":  l.maxBackups,
		"max_age_days": l.maxAge,
		"compress":     l.compress,
		"rotate_daily": l.rotateDaily,
		"current_date": l.currentDate.Format("2006-01-02"),
		"log_path":     l.logPath,
	}
}

func (l *Logger) ForceRotation() error {
	return l.rotate()
}

func (l *Logger) CheckRotationStatus() map[string]interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	currentDate := now.Truncate(24 * time.Hour)

	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	timeToMidnight := nextMidnight.Sub(now)

	return map[string]interface{}{
		"current_time":      now.Format("2006-01-02 15:04:05"),
		"current_date":      l.currentDate.Format("2006-01-02"),
		"today_date":        currentDate.Format("2006-01-02"),
		"date_changed":      !currentDate.Equal(l.currentDate),
		"rotate_daily":      l.rotateDaily,
		"timer_active":      l.dailyTimer != nil,
		"current_size":      l.currentSize,
		"max_size":          l.maxSize,
		"needs_rotation":    l.needsRotationUnsafe(0),
		"next_midnight":     nextMidnight.Format("2006-01-02 15:04:05"),
		"time_to_midnight":  timeToMidnight.String(),
		"hours_to_midnight": timeToMidnight.Hours(),
	}
}

func (l *Logger) needsRotationUnsafe(messageSize int) bool {
	sizeRotation := l.maxSize > 0 && l.currentSize+float32(messageSize) > l.maxSize

	var dateRotation bool
	if l.rotateDaily {
		currentDate := time.Now().Truncate(24 * time.Hour)
		dateRotation = !currentDate.Equal(l.currentDate)

		if !dateRotation && l.file != nil {
			if stat, err := l.file.Stat(); err == nil {
				fileModDate := stat.ModTime().Truncate(24 * time.Hour)
				if !fileModDate.Equal(currentDate) {
					dateRotation = true
				}
			}
		}
	}

	return sizeRotation || dateRotation
}

func getStringOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getFloat32OrDefault(value, defaultValue float32) float32 {
	if value == 0 {
		return defaultValue
	}
	return value
}

func getIntOrDefault(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}
