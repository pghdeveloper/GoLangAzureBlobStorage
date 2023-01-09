package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Tests struct {
	name string
	response []string
	mockResponse []string
}

type mockGetFiles struct {
	response []string
}
 
func (c *mockGetFiles) GetFilesFromCloud(ctx context.Context, containerId string) []string{
	fmt.Println("Mock Function called")
	return c.response
}

func TestGetFileNames(t *testing.T) {
	tests := []Tests {
		{
			name: "Get Files",
			response: []string{"1","2"},
			mockResponse: []string{"1","2"},
		},
		// {
		// 	name: "Get Files",
		// 	response: []string{},
		// 	mockResponse: []string{},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){	
			Repos = &mockGetFiles{
				response: test.mockResponse,
			}
			w := httptest.NewRecorder()

			ctx := GetTestGinContext(w)

			params := []gin.Param{
				{
					Key: "containerId",
					Value: "1",
				},
			}

			MockJsonGet(ctx, params)

			GetFileNames(ctx)
			var response []string
			
			err := json.Unmarshal((w.Body.Bytes()), &response)
			assert.NoError(t, err)

			fmt.Println("Body: " + w.Body.String())
			fmt.Println("Response:", response)

			assert.EqualValues(t, http.StatusOK, w.Code)

			assert.Equal(t, test.response, response)
		})
	}
}

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx
}

func MockJsonGet(c *gin.Context, params gin.Params) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("strArray", "['1','2']")

	c.Params = params
}