package tests

import (
	"context"
	"sort"

	"github.com/khaiql/dbcleaner"
	"github.com/mises-id/sns/lib/db"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
	dbcleaner.DbCleaner
}

func (suite *BaseTestSuite) SetupSuite() {
	suite.DbCleaner = dbcleaner.New()
	db.SetupMongo(context.Background())
}

func (suite *BaseTestSuite) TearDownSuite() {
}

func (suite *BaseTestSuite) Clean(collections ...string) {
	sort.Strings(collections)
	suite.DbCleaner.Acquire(collections...)
	for _, collection := range collections {
		_ = db.DB().Collection(collection).Drop(context.Background())
	}
}

func (suite *BaseTestSuite) Acquire(collections ...string) {
	sort.Strings(collections)
	suite.DbCleaner.Acquire(collections...)
}
