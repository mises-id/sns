package storage

import (
	"context"
	"fmt"
	"io"
	"os"
)

type FileStore struct{}

func (s *FileStore) Upload(ctx context.Context, filePath, filename string, file File) error {
	var err error
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fmt.Println(filePath + filename)
	dst, err := os.Create(filePath + filename)
	if err != nil {
		fmt.Println(filePath+filename, err)
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
