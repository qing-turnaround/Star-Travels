package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"web_app/controller"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	r   *gin.Engine
	url string
)

func TestCreatePostHandler(t *testing.T) {
	url := "http://120.79.17.230:9999/posts"
	// 定义多个测试用例（定义Body时，采用json格式）
	tests := []struct {
		testName   string             // 测试用例的名字
		testBody   string             // 测试主体
		testExpect controller.ResCode // 预期结果(只比较返回 状态码)
	}{
		{"test1", `{"community_id": 2, "title": "Test1", "content": "Just Test1"}`, controller.CodeSuccess},
		{"test2", `{"community_id": 2, "title": "", "content": "Just Test2"}`, controller.CodeInvalidParam},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// 构建一个http请求
			req := httptest.NewRequest("POST", url, strings.NewReader(tt.testBody)) //这里使用httptest
			// 设置Token
			req.Header = map[string][]string{"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0MjE4NDkyOTEzMTM2ODQ0OCwidXNlcm5hbWUiOiLor7jokZvpnZIiLCJleHAiOjE2NTE5MjAwMTcsImlzcyI6IuivuOiRm-mdkiJ9.MYjHmiEkd6j1o0zVBm4rwJMMK2VSDisc5vDFqoEDkks"}}
			// mock 一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			// 校验状态码是否符合预期
			if !assert.Equal(t, http.StatusOK, w.Code) {
				t.Errorf("test Failed, the result is %v", w.Code)
			}

			// 解析并检验响应内容是否复合预期
			resp := new(controller.ResponseData)
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert.Nil(t, err)
			assert.Equal(t, tt.testExpect, resp.Code)
		})
	}
}
