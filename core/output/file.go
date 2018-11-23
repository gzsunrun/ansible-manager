package output

import (
	"bufio"
	"os"
)

type FileLog struct {
	fd       *os.File
	FilePath string
}

func NewFileLog(path string) (*FileLog, error) {
	fg := new(FileLog)
	fg.FilePath = path

	return fg, nil
}

func (fg *FileLog) Write(msg string) error {
	if fg.fd == nil {
		fd, err := os.OpenFile(fg.FilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		fg.fd = fd
	}
	fg.fd.WriteString(msg + "\n")
	return nil
}

func (fg *FileLog) Read() ([]string, error) {
	ret := make([]string, 0)
	file, err := os.Open(fg.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (fg *FileLog) Close() error {
	return fg.fd.Close()
}

func (fg *FileLog) Clean() error {
	return os.Remove(fg.FilePath)
}
