package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/app/models/meta"
	"github.com/mises-id/sns/lib/codes"
	"github.com/mises-id/sns/lib/db"
	"github.com/mises-id/sns/lib/pagination"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ParentID      primitive.ObjectID `bson:"parent_id,omitempty"`
	OriginID      primitive.ObjectID `bson:"origin_id,omitempty"`
	UID           uint64             `bson:"uid,omitempty"`
	FromType      enum.FromType      `bson:"from_type"`
	StatusType    enum.StatusType    `bson:"status_type"`
	Meta          json.RawMessage    `bson:"meta,omitempty"`
	Content       string             `bson:"content,omitempty" validate:"min=0,max=4000"`
	CommentsCount uint64             `bson:"comments_count,omitempty"`
	LikesCount    uint64             `bson:"likes_count,omitempty"`
	ForwardsCount uint64             `bson:"forwards_count,omitempty"`
	DeletedAt     *time.Time         `bson:"deleted_at,omitempty"`
	CreatedAt     time.Time          `bson:"created_at,omitempty"`
	UpdatedAt     time.Time          `bson:"updated_at,omitempty"`
	User          *User              `bson:"-"`
	IsLiked       bool               `bson:"-"`
	ParentStatus  *Status            `bson:"-"`
	OriginStatus  *Status            `bson:"-"`
	metaData      meta.MetaData      `bson:"-"`
}

func (s *Status) validate(ctx context.Context) error {
	logrus.Info("xxxxx")
	err := Validate.Struct(s)
	if err != nil {
		return codes.ErrUnprocessableEntity
	}
	return nil
}

func (s *Status) BeforeCreate(ctx context.Context) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	var err error
	if !s.ParentID.IsZero() {
		s.ParentStatus, err = FindStatus(ctx, s.ParentID)
		if err != nil {
			return err
		}
		s.OriginID = s.ParentStatus.OriginID
		if s.OriginID.IsZero() {
			s.OriginID = s.ParentID
		}
	}
	if !s.OriginID.IsZero() {
		s.OriginStatus, err = FindStatus(ctx, s.OriginID)
		if err != nil {
			return err
		}
	}
	return s.validate(ctx)
}

func (s *Status) AfterCreate(ctx context.Context) error {
	var err error
	counterKey := s.FromType.CounterKey()
	if s.ParentStatus != nil {
		err = s.ParentStatus.IncStatusCounter(ctx, counterKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Status) IncStatusCounter(ctx context.Context, counterKey string, values ...int) error {
	if counterKey == "" {
		return nil
	}
	value := 1
	if len(values) > 0 {
		value = values[0]
	}
	return db.DB().Collection("statuses").FindOneAndUpdate(ctx, bson.M{"_id": s.ID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   counterKey,
				Value: value,
			}}},
		}).Err()
}

func (s *Status) GetMetaData() (meta.MetaData, error) {
	var err error
	if s.metaData == nil {
		s.metaData, err = meta.BuildStatusMeta(s.StatusType, s.Meta)
	}
	return s.metaData, err
}

func FindStatus(ctx context.Context, id primitive.ObjectID) (*Status, error) {
	status := &Status{}
	err := db.ODM(ctx).First(status, bson.M{"_id": id}).Error
	if err != nil {
		return nil, err
	}
	if err = preloadRelatedStatus(ctx, status); err != nil {
		return nil, err
	}
	if err = preloadAttachment(ctx, status); err != nil {
		return nil, err
	}
	return status, preloadStatusUser(ctx, status)
}

type CreateStatusParams struct {
	UID        uint64
	ParentID   primitive.ObjectID
	StatusType enum.StatusType
	FromType   enum.FromType
	Content    string
	MetaData   meta.MetaData
}

func CreateStatus(ctx context.Context, params *CreateStatusParams) (*Status, error) {
	status := &Status{
		UID:        params.UID,
		StatusType: params.StatusType,
		FromType:   params.FromType,
		ParentID:   params.ParentID,
		Content:    params.Content,
	}
	var err error
	if params.MetaData != nil {
		status.Meta, err = json.Marshal(params.MetaData)
		if err != nil {
			return nil, err
		}
	}
	if err = status.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	if err = db.ODM(ctx).Create(status).Error; err != nil {
		return nil, err
	}
	if err = status.AfterCreate(ctx); err != nil {
		return nil, err
	}
	if err = preloadRelatedStatus(ctx, status); err != nil {
		return nil, err
	}
	if err = preloadAttachment(ctx, status); err != nil {
		return nil, err
	}
	return status, preloadStatusUser(ctx, status)
}

func DeleteStatus(ctx context.Context, id primitive.ObjectID) error {
	_, err := db.DB().Collection("statuses").DeleteOne(ctx, bson.M{"_id": id})
	return err
}

type ListStatusParams struct {
	UIDs           []uint64
	ParentStatusID primitive.ObjectID
	FromTypes      []enum.FromType
	PageParams     *pagination.PageQuickParams
}

func ListStatus(ctx context.Context, params *ListStatusParams) ([]*Status, pagination.Pagination, error) {
	if params.PageParams == nil {
		params.PageParams = pagination.DefaultQuickParams()
	}
	statuses := make([]*Status, 0)
	chain := db.ODM(ctx)
	if params.UIDs != nil && len(params.UIDs) > 0 {
		chain = chain.Where(bson.M{"uid": bson.M{"$in": params.UIDs}})
	}
	if !params.ParentStatusID.IsZero() {
		chain = chain.Where(bson.M{"parent_id": params.ParentStatusID})
	}
	if params.FromTypes != nil {
		chain = chain.Where(bson.M{"from_type": bson.M{"$in": params.FromTypes}})
	}
	paginator := pagination.NewQuickPaginator(params.PageParams.Limit, params.PageParams.NextID, chain)
	page, err := paginator.Paginate(&statuses)
	if err != nil {
		return nil, nil, err
	}
	if err = preloadRelatedStatus(ctx, statuses...); err != nil {
		return nil, nil, err
	}
	if err = preloadAttachment(ctx, statuses...); err != nil {
		return nil, nil, err
	}
	return statuses, page, preloadStatusUser(ctx, statuses...)
}

func ListCommentStatus(ctx context.Context, statusID primitive.ObjectID, pageParams *pagination.PageQuickParams) ([]*Status, pagination.Pagination, error) {
	if pageParams == nil {
		pageParams = pagination.DefaultQuickParams()
	}
	statuses := make([]*Status, 0)
	chain := db.ODM(ctx).Where(bson.M{"parent_id": statusID, "from_type": enum.FromComment})
	paginator := pagination.NewQuickPaginator(pageParams.Limit, pageParams.NextID, chain)
	page, err := paginator.Paginate(&statuses)
	if err != nil {
		return nil, nil, err
	}
	if err = preloadRelatedStatus(ctx, statuses...); err != nil {
		return nil, nil, err
	}
	if err = preloadAttachment(ctx, statuses...); err != nil {
		return nil, nil, err
	}
	return statuses, page, preloadStatusUser(ctx, statuses...)
}

func preloadStatusUser(ctx context.Context, statuses ...*Status) error {
	userIds := make([]uint64, 0)
	for _, status := range statuses {
		userIds = append(userIds, status.UID)
	}
	users := make([]*User, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": userIds}}).Find(&users).Error
	if err != nil {
		return err
	}
	if err = PreloadUserAvatar(ctx, users...); err != nil {
		return err
	}
	if err = BatchSetFolloweState(ctx, users...); err != nil {
		return err
	}
	userMap := make(map[uint64]*User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	for _, status := range statuses {
		status.User = userMap[status.UID]
	}
	return nil
}

func preloadRelatedStatus(ctx context.Context, statuses ...*Status) error {
	statusIds := make([]primitive.ObjectID, 0)
	for _, status := range statuses {
		if !status.ParentID.IsZero() {
			statusIds = append(statusIds, status.ParentID)
		}
		if !status.OriginID.IsZero() {
			statusIds = append(statusIds, status.OriginID)
		}
	}
	relatedStatuses := make([]*Status, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": statusIds}}).Find(&relatedStatuses).Error
	if err != nil {
		return err
	}
	if err = preloadStatusUser(ctx, relatedStatuses...); err != nil {
		return err
	}
	if err = preloadAttachment(ctx, relatedStatuses...); err != nil {
		return err
	}
	statusMap := make(map[primitive.ObjectID]*Status)
	for _, status := range relatedStatuses {
		statusMap[status.ID] = status
	}
	for _, status := range statuses {
		status.ParentStatus = statusMap[status.ParentID]
		status.OriginStatus = statusMap[status.OriginID]
	}
	return nil
}

func preloadAttachment(ctx context.Context, statuses ...*Status) error {
	attachmentIDs := make([]uint64, 0)
	linkMetas := make([]*meta.LinkMeta, 0)
	for _, status := range statuses {
		if status.StatusType != enum.LinkStatus {
			continue
		}
		metaData, err := status.GetMetaData()
		if err != nil {
			return err
		}
		linkMeta := metaData.(*meta.LinkMeta)
		attachmentIDs = append(attachmentIDs, linkMeta.AttachmentID)
		linkMetas = append(linkMetas, linkMeta)
	}
	attachments := make([]*Attachment, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": attachmentIDs}}).Find(&attachments).Error
	if err != nil {
		return err
	}
	attachmentMap := make(map[uint64]*Attachment)
	for _, attachment := range attachments {
		attachmentMap[attachment.ID] = attachment
	}
	for _, linkMeta := range linkMetas {
		if attachmentMap[linkMeta.AttachmentID] != nil {
			linkMeta.AttachmentURL = attachmentMap[linkMeta.AttachmentID].FileUrl()
		}
	}
	return nil
}
