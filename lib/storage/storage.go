package storage

import (
	"context"

	"github.com/mises-id/sns/config/env"
)

var (
	storageService IStorageService
	Prefix         = "upload/"
)

func init() {
	switch env.Envs.StorageProvider {
	default:
		storageService = &FileStore{}
	case "local":
		storageService = &FileStore{}
	case "oss":
		storageService = &OSSStorage{}
	}
}

func UploadFile(ctx context.Context, path, filename string, file File) error {
	return storageService.Upload(ctx, env.Envs.RootPath+Prefix+path, filename, file)
}

type IStorageService interface {
	Upload(ctx context.Context, filePath, filename string, file File) error
}

type File interface {
	Read(p []byte) (n int, err error)
}
