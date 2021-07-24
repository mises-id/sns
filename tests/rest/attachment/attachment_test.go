package attachment

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/mises-id/sns/config/env"
	"github.com/mises-id/sns/tests/rest"
	"github.com/stretchr/testify/suite"
)

type AttachmentServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
	files       []string
}

func (suite *AttachmentServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "attachments"}
	suite.files = []string{path.Join(env.Envs.RootPath, "upload", "attachment")}
}

func (suite *AttachmentServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *AttachmentServerSuite) SetupTest() {
	suite.Clean(suite.collections...)
	suite.Acquire(suite.collections...)
}

func (suite *AttachmentServerSuite) CleanFiles() {
	for _, file := range suite.files {
		_ = os.RemoveAll(file)
	}
}

func (suite *AttachmentServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
	suite.CleanFiles()
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
		url := fmt.Sprintf("http://localhost/upload/attachment/%s/2/test.mp4", time.Now().Format("2006/01/02"))
		resp.Value("data").Object().Value("url").Equal(url)
	})
}
