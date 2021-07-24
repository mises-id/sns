package attachment

import (
	"net/http"
	"testing"

	"github.com/mises-id/sns/tests/rest"
	"github.com/stretchr/testify/suite"
)

type AttachmentServerSuite struct {
	rest.RestBaseTestSuite
	tables []string
}

func (suite *AttachmentServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
}

func (suite *AttachmentServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *AttachmentServerSuite) SetupTest() {
	suite.Clean(suite.tables...)
	suite.Acquire(suite.tables...)
}

func (suite *AttachmentServerSuite) TearDownTest() {
	suite.Clean(suite.tables...)
}

func TestAttachmentServer(t *testing.T) {
	suite.Run(t, &AttachmentServerSuite{})
}

func (suite *AttachmentServerSuite) TestUpload() {
	suite.T().Run("update image success", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/attachment").WithMultipart().
			WithFile("file", "../../test.jpg").WithFormField("file_type", "image").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})
}
