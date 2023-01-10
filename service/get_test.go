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
	statusCode int
}

type TestsNoFileNames struct {
	name string
	response string
	mockResponse []string
	statusCode int
}

type ErrorResponse struct {
	Message string
}

type mockGetFiles struct {
	response []string
}
 
func (c *mockGetFiles) GetFilesFromCloud(ctx context.Context, containerId string) []string{
	fmt.Println("Mock Function called")
	return c.response
}

func TestGetFileNamesReturns404(t * testing.T) {
	test := &TestsNoFileNames {
			name: "Get Files But Returns Empty String Array",
			mockResponse: []string{},
			response: `{"Message":"Files not found"}`,
			statusCode: 404,
	}

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
	fmt.Println("GetFileNames function finished")

	var errorResponse ErrorResponse
	
	err := json.Unmarshal((w.Body.Bytes()), &errorResponse)
	assert.NoError(t, err)

	e, err := json.Marshal(errorResponse)
	assert.NoError(t, err)

	fmt.Println("Body: " + w.Body.String())
	fmt.Println("Response: " + errorResponse.Message)
	fmt.Println("Response after marshalling: " + string(e))

	errorResponseString := string(e)
	assert.EqualValues(t, test.statusCode, w.Code)

	assert.Equal(t, test.response, errorResponseString)
}

func TestGetFileNames(t *testing.T) {
	tests := []Tests {
		{
			name: "Get Files",
			response: []string{"1","2"},
			mockResponse: []string{"1","2"},
			statusCode: 200,
		},
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
			fmt.Println("GetFileNames function finished")
			var response []string
			
			err := json.Unmarshal((w.Body.Bytes()), &response)
			assert.NoError(t, err)

			fmt.Println("Body: " + w.Body.String())
			fmt.Println("Response:", response)

			assert.EqualValues(t, test.statusCode, w.Code)

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
	//c.Set("strArray", "['1','2']")

	c.Params = params
}