package storage

import (
	"context"

	"github.com/mises-id/sns/lib/codes"
)

type OSSStorage struct{}

func (s *OSSStorage) Upload(ctx context.Context, filePath, filename string, file File) error {
	return codes.ErrUnimplemented
}
