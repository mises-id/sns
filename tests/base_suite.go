package tests

import (
	"context"

	"github.com/mises-id/sns/lib/db"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
}

func (s *BaseTestSuite) SetupSuite() {
	db.SetupMongo(context.Background())
}

func (s *BaseTestSuite) TearDownSuite() {
}

func (s *BaseTestSuite) Clean(tables ...string) {
}

func (s *BaseTestSuite) Acquire(tables ...string) {
}
