package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type DownloadFileTests struct {
	name string
	response []byte
	mockResponse *bytes.Buffer
	statusCode int
}

type NoDownloadFileTests struct {
	name string
	response string
	mockResponse *bytes.Buffer
	statusCode int
}

type errorResponse struct {
	Message string
}

type mockDownloadFile struct {
	response *bytes.Buffer
	err error
}
 
func (c *mockDownloadFile) DownloadFileFromCloud(ctx context.Context, containerId string, fileName string) (*bytes.Buffer, error){
	fmt.Println("Mock Function called")
	return c.response, c.err
}

func TestDownloadFileReturns500(t *testing.T) {
	buf := &bytes.Buffer{}
	
	test := &NoDownloadFileTests  {
		name: "Download File Not Successful",
		mockResponse: buf,
		response: `{"Message":"Issue Downloading the file"}`,
		statusCode: 500,
	}

	mockError := errors.New("Error")
	DownloadRepos = &mockDownloadFile{
		response: test.mockResponse,
		err: mockError,
	}

	w := httptest.NewRecorder()

	ctx := getTestGinContext(w)

	params := []gin.Param{
		{
			Key: "fileName",
			Value: "123.pdf",
		},
	}

	mockJsonGet(ctx, params)

	DownloadFile(ctx)
	fmt.Println("DownloadFiles function finished")

	var errorResponse errorResponse
	
	err := json.Unmarshal((w.Body.Bytes()), &errorResponse)
	assert.NoError(t, err)

	e, err := json.Marshal(errorResponse)
	assert.NoError(t, err)

	fmt.Println("Body: " + w.Body.String())
	fmt.Println("Response:", errorResponse.Message)

	errorResponseString := string(e)
	assert.EqualValues(t, test.statusCode, w.Code)

	assert.Equal(t, test.response, errorResponseString)
}

func TestDownloadFile(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("Hello  World")

    fmt.Println("Begin Testing Download File")
    // Write strings to the Buffer.
	fmt.Println("Begin Testing Download File 2")
	
	test := &DownloadFileTests {
		name: "Download File Successfully",
		mockResponse: buf,
		response: buf.Bytes(),
		statusCode: 200,
	}

	DownloadRepos = &mockDownloadFile{
		response: test.mockResponse,
		err: nil,
	}

	w := httptest.NewRecorder()

	ctx := getTestGinContext(w)

	params := []gin.Param{
		{
			Key: "fileName",
			Value: "123.pdf",
		},
	}

	mockJsonGet(ctx, params)

	DownloadFile(ctx)
	fmt.Println("DownloadFiles function finished")

	var response []byte
			
	err := json.Unmarshal((w.Body.Bytes()), &response)
	assert.NoError(t, err)

	fmt.Println("Body: " + w.Body.String())
	fmt.Println("Response:", response)

	assert.EqualValues(t, test.statusCode, w.Code)

	assert.Equal(t, test.response, response)
}

func getTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx
}

func mockJsonGet(c *gin.Context, params gin.Params) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")

	c.Params = params
}