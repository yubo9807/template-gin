package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileItem struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Ext      string     `json:"ext"`
	IsDir    bool       `json:"isDir"`
	Size     int64      `json:"size"`
	Time     int64      `json:"time"`
	Children []FileItem `json:"children"`
}

func FileInfo(filename string, deep bool) (FileItem, error) {
	info := FileItem{}
	file, err := os.Open(filename)
	if err != nil {
		return FileItem{}, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	info.Name = fileInfo.Name()
	info.IsDir = fileInfo.IsDir()
	info.Path = filename
	if !info.IsDir {
		info.Ext = filepath.Ext(info.Name)
	}
	info.Size, _ = file.Seek(0, io.SeekEnd)
	info.Time = fileInfo.ModTime().Unix()
	info.Children = []FileItem{}
	if deep && info.IsDir {
		info.Children, _ = FileCatalog(filename, true)
	}
	return info, nil
}

func FileCatalog(filename string, deep bool) ([]FileItem, error) {
	files, err := os.ReadDir(filename)
	if err != nil {
		return nil, err
	}
	list := []FileItem{}
	for _, file := range files {
		info, _ := file.Info()
		name := info.Name()
		path := filename + string(filepath.Separator) + name
		isDir := info.IsDir()
		ext := ""
		if !isDir {
			ext = filepath.Ext(name)
		}
		children := []FileItem{}
		if deep && isDir {
			child, _ := FileCatalog(path, deep)
			children = child
		}
		list = append(list, FileItem{
			Name:     name,
			Path:     path,
			Ext:      ext,
			IsDir:    isDir,
			Size:     info.Size(),
			Time:     info.ModTime().Unix(),
			Children: children,
		})
	}
	return list, nil
}

func FileCreateWithDirs(path string) error {
	// 获取文件的目录路径
	dir := filepath.Dir(path)
	// 检查目录是否存在，如果不存在则创建它
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("无法创建目录 %s: %v", dir, err)
		}
	}

	// 创建文件
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("无法创建文件 %s: %v", path, err)
	}
	defer file.Close()

	return nil
}

// isPathWithinBase 检查给定路径是否在基础路径内，防止路径遍历
func isPathWithinBase(path, base string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && rel != ".."
}

// unzip 将 zip 文件解压到目标目录
func FileUnzip(zipPath, destDir string) error {
	// 1. 打开 ZIP 文件
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("打开 ZIP 文件失败: %v", err)
	}
	defer reader.Close()

	// 2. 遍历 ZIP 包中的每个文件/目录
	for _, file := range reader.File {
		// 构建目标路径 (防止路径遍历攻击)
		filePath := filepath.Join(destDir, file.Name)

		// 安全检查：确保目标路径在 destDir 内
		if !isPathWithinBase(filePath, destDir) {
			return fmt.Errorf("文件路径非法，可能涉及路径遍历: %s", file.Name)
		}

		// 3. 处理目录
		if file.FileInfo().IsDir() {
			fmt.Printf("创建目录: %s\n", filePath)
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return fmt.Errorf("创建目录失败 %s: %v", filePath, err)
			}
			continue
		}

		// 4. 处理文件
		// 确保文件所在目录存在
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("创建文件目录失败 %s: %v", filepath.Dir(filePath), err)
		}

		// 打开 ZIP 中的文件
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("打开 ZIP 中的文件失败 %s: %v", file.Name, err)
		}

		// 创建本地目标文件
		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			return fmt.Errorf("创建本地文件失败 %s: %v", filePath, err)
		}

		// 将数据从 ZIP 文件复制到本地文件
		_, err = io.Copy(targetFile, fileReader)
		// 关闭文件 (注意：先关闭源，再关闭目标)
		fileReader.Close()
		closeErr := targetFile.Close()
		if err != nil {
			return fmt.Errorf("复制文件内容失败 %s: %v", file.Name, err)
		}
		if closeErr != nil {
			return fmt.Errorf("关闭目标文件失败 %s: %v", filePath, closeErr)
		}
	}

	return nil
}
