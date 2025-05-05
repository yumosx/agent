package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/yumosx/got/pkg/suitex"
	"testing"
)

type HandlerSuite struct {
	server *gin.Engine
	suite.Suite
}

func (h *HandlerSuite) SetupSuite() {
}

func (h *HandlerSuite) TestChat() {
	t := h.T()

	type Msg struct {
		Message string `json:"message"`
	}

	testCases := []struct {
		Name     string
		Input    Msg
		Expect   string
		WantCode int64
	}{
		{
			Name:  "请求",
			Input: Msg{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := json.Marshal(tc.Input)
			require.NoError(t, err)
			response, err := suitex.MockPostResponse(h.server, "/chat", req)
			require.NoError(t, err)
			assert.Equal(t, response.Code, tc.WantCode)
		})
	}
}
