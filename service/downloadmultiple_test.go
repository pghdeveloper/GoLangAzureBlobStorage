package service

import (
	"archive/zip"
	"bytes"
	"context"
	"example/GoLangAzureBlobStorage/lib"
	"fmt"
	"testing"
)

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

func TestDownloadMultipleFile(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("Hello  World")
	
	inMemoryFiles := []*lib.InMemoryFile {
		{
			FileName: "Test.pdf",
			Content: buf.Bytes(),
		},
		{
			FileName: "Test1.pdf",
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

}

func createZipFileForMockTest(inMemoryFiles []*lib.InMemoryFile) []byte{
	fmt.Println("we are in the zipData function")
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