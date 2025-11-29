package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func LogWrite(catalog, msg string) {
	filename := time.Now().String()[0:10] + ".log"
	url := catalog + string(filepath.Separator) + filename
	file, err := os.OpenFile(url, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		FileCreateWithDirs(url)
		file, _ = os.OpenFile(url, os.O_CREATE, 0666)
	}
	defer file.Close()
	file.WriteString(msg + "\n\n")
}

// 清理日志
func LogClear(catalog string, t int) {
	flag := time.Now().AddDate(0, 0, -t).Unix()
	filepath.Walk("logs", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".log") {
			return nil
		}

		name := strings.Split(info.Name(), ".")[0]
		t, _ := time.Parse("2006-01-02", name)
		if t.Unix() < flag {
			os.Remove(path)
		}
		return nil
	})
}
