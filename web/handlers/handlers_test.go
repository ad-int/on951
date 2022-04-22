package handlers

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

//ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
//ctx.Set("article_id1", 1)
//handlers.GetComments(ctx)
//image_links_parser.Process(`ikljghi  <img src="data:application/json;base64," /> 232refwdsf <img src="data:image/png;Es!a, QQQ" />grd`)
//return

type handlersTestSuite struct {
	suite.Suite
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}
