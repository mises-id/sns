package attachment

import (
	"net/http"
	"testing"

	"github.com/mises-id/sns/tests/rest"
	"github.com/stretchr/testify/suite"
)

type AttachmentServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
}

func (suite *AttachmentServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "attachments"}
}

func (suite *AttachmentServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *AttachmentServerSuite) SetupTest() {
	suite.Clean(suite.collections...)
	suite.Acquire(suite.collections...)
}

func (suite *AttachmentServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
}

func TestAttachmentServer(t *testing.T) {
	suite.Run(t, &AttachmentServerSuite{})
}

func (suite *AttachmentServerSuite) TestUpload() {
	suite.T().Run("upload image success", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/attachment").WithMultipart().
			WithFile("file", "../../test.jpg").WithFormField("file_type", "image").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("id").Equal(1)
	})

	suite.T().Run("upload video success", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/attachment").WithMultipart().
			WithFile("file", "../../test.mp4").WithFormField("file_type", "video").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("id").Equal(2)
		resp.Value("data").Object().Value("url").Equal("http://localhost/upload/attachment/2021/07/24/2/test.mp4")
	})
}
