package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"example/GoLangAzureBlobStorage/lib"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type NotCorrectJson struct {
	Name string
}

type DownloadMultipleFileTests struct {
	name string
	response []byte
	mockResponse []*lib.InMemoryFile
	statusCode int
}

type mockDownloadMultipleFile struct {
	response []*lib.InMemoryFile
	err error
}
 
func (c *mockDownloadMultipleFile) DownloadMultipleFilesFromCloud(ctx context.Context, containerIds lib.Containers) ([]*lib.InMemoryFile, error){
	fmt.Println("Mock Function called")
	return c.response, c.err
}

func TestDownloadFileReturn400BecauseJsonRequestNotAccurate(t *testing.T) {
	w := httptest.NewRecorder()

	ctx := getTestGinContext(w)

	notCorrectJson := new(NotCorrectJson)
	notCorrectJson.Name = "1"

	mockJsonPost(ctx, notCorrectJson)

	DownloadMultiple(ctx)

	var errorResponse errorResponse
	
	err := json.Unmarshal((w.Body.Bytes()), &errorResponse)
	assert.NoError(t, err)

	e, err := json.Marshal(errorResponse)
	assert.NoError(t, err)

	fmt.Println("Body: " + w.Body.String())
	fmt.Println("Response:", errorResponse.Message)

	errorResponseString := string(e)
	assert.EqualValues(t, 400, w.Code)

	assert.Equal(t, `{"Message":"Json Request not correct"}`, errorResponseString)
}

func TestDownloadMultipleFile(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("Hello  World")
	
	inMemoryFiles := []*lib.InMemoryFile {
		{
			FileName: "Test.txt",
			Content: buf.Bytes(),
		},
		{
			FileName: "Test1.txt",
			Content: buf.Bytes(),
		},
	}

	test := &DownloadMultipleFileTests {
		name: "Download Multiple File Successfully",
		mockResponse: inMemoryFiles,
		response: createZipFileForMockTest(inMemoryFiles),
		statusCode: 200,
	}

	DownloadMultipleRepos = &mockDownloadMultipleFile{
		response: test.mockResponse,
		err: nil,
	}

	w := httptest.NewRecorder()

	ctx := getTestGinContext(w)

	container := new(lib.Containers)
	container.ContainerIds = []string{"1","2"}

	mockJsonPost(ctx, container)

	DownloadMultiple(ctx)
	fmt.Println("Download Multiple function finished")

	var response []byte
			
	err := json.Unmarshal((w.Body.Bytes()), &response)
	assert.NoError(t, err)

	fmt.Println("Body: " + w.Body.String())
	fmt.Println("Response:", response)

	assert.EqualValues(t, test.statusCode, w.Code)

	assert.Equal(t, test.response, response)

}

func mockJsonPost(c *gin.Context, jsonStruct interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	fmt.Println(jsonStruct)

	jsonbytes, err := json.Marshal(jsonStruct)
    if err != nil {
        panic(err)
    }
    
    // the request body must be an io.ReadCloser
    // the bytes buffer though doesn't implement io.Closer,
    // so you wrap it in a no-op closer
    c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	fmt.Println(c.Request.Body)
	fmt.Println("finished mock Json Post")
}

func createZipFileForMockTest(inMemoryFiles []*lib.InMemoryFile) []byte{
	fmt.Println("we are in the Mock zipData function")
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)

	for _, file := range inMemoryFiles {
		zipFile, err := zipWriter.Create(file.FileName)
		if err != nil {
			fmt.Println(err)
		}
		_, err = zipFile.Write(file.Content)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Make sure to check the error on Close.
	err := zipWriter.Close()
	if err != nil {
		fmt.Println(err)
	}

	//write the zipped file to the disk
	return buf.Bytes()
}