package models

import (
	"context"
	"strconv"
	"time"

	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/config/env"
	"github.com/mises-id/sns/lib/db"
	"github.com/mises-id/sns/lib/storage"
)

type Attachment struct {
	ID        uint64        `bson:"_id"`
	Filename  string        `bson:"filename,omitempty"`
	FileType  enum.FileType `bson:"file_type"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
	file      storage.File
}

func (a *Attachment) BeforeCreate(ctx context.Context) error {
	var err error
	a.ID, err = getNextSeq(ctx, "attachmentid")
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return err
}

func (a *Attachment) FileUrl() string {
	return env.Envs.AssetHost + storage.Prefix + a.filePath() + a.Filename
}

func (a *Attachment) filePath() string {
	if a.ID == 0 {
		return "tmp"
	}
	return "attachment/" + a.CreatedAt.Format("2006/01/02/") + strconv.Itoa(int(a.ID)) + "/"
}

func (a *Attachment) UploadFile(ctx context.Context) error {
	return storage.UploadFile(ctx, a.filePath(), a.Filename, a.file)
}

func CreateAttachment(ctx context.Context, tp enum.FileType, filename string, file storage.File) (*Attachment, error) {
	attachment := &Attachment{
		Filename: filename,
		FileType: tp,
		file:     file,
	}
	if err := attachment.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	if err := attachment.UploadFile(ctx); err != nil {
		return nil, err
	}
	_, err := db.DB().Collection("attachments").InsertOne(ctx, attachment)
	return attachment, err
}
