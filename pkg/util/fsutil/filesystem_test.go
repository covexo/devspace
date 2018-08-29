package fsutil

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/covexo/devspace/pkg/util/randutil"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestWriteToFileAndReadFile(t *testing.T) {

	//Let's create a new file and check if the content is correct.

	randomString, e := randutil.GenerateRandomString(10)

	assert.Nil(t, e)

	writeData := []byte("Content " + randomString)

	fileName := os.TempDir() + "/" + randomString

	e = WriteToFile(writeData, fileName)

	assert.Nilf(t, e, errors.Details(e))

	//There should be 18 bytes in the file. We'll only read 17 to test out whether this method reads the correct amount of bytes.
	readedData, e := ReadFile(fileName, 17)

	assert.Nil(t, e)
	assert.Len(t, readedData, 17)

	for n, byte := range readedData {
		assert.Equal(t, byte, writeData[n])
	}

	//Now let's overwrite the content

	newData := []byte("New Content " + randomString)

	WriteToFile(newData, fileName)

	//Read everything
	readedData, e = ReadFile(fileName, -1)

	assert.Nil(t, e)
	assert.Len(t, readedData, 22)

	for n, byte := range readedData {
		assert.Equal(t, byte, newData[n])
	}

}

func TestCopy(t *testing.T) {

	randomString, e := randutil.GenerateRandomString(10)
	sourceFile, e := ioutil.TempFile("", randomString)
	assert.Nil(t, e)
	defer os.Remove(sourceFile.Name())

	randomString, e = randutil.GenerateRandomString(10)
	assert.Nil(t, e)
	destPath := os.TempDir() + "/" + randomString

	randomString, e = randutil.GenerateRandomString(10)
	WriteToFile([]byte(randomString), sourceFile.Name())
	assert.Nil(t, e)

	Copy(sourceFile.Name(), destPath)
	defer os.Remove(destPath)

	sourceContent, e1 := ReadFile(sourceFile.Name(), -1)
	destContent, e2 := ReadFile(destPath, -1)

	assert.Nil(t, e1)
	assert.Nil(t, e2)

	assert.Equal(t, randomString, string(sourceContent))
	assert.Equal(t, randomString, string(destContent))

}

func TestGetHomeDir(t *testing.T) {

	homeDirByMethod := GetHomeDir()
	homeDirByOS := os.Getenv("HOME")
	if homeDirByOS == "" {
		homeDirByOS = os.Getenv("USERPROFILE")
	}

	assert.Equal(t, homeDirByOS, homeDirByMethod)
}

func TestGetCurrentGofileDir(t *testing.T) {

	currentGofileDirByMethod := GetCurrentGofileDir()
	expectedPath := os.Getenv("GOPATH") + "\\src\\github.com\\covexo\\devspace\\pkg\\util\\fsutil"

	expectedPath = strings.Replace(expectedPath, "\\", "/", -1)
	currentGofileDirByMethod = strings.Replace(currentGofileDirByMethod, "\\", "/", -1)

	assert.Equal(t, expectedPath, currentGofileDirByMethod)
}

func TestGetCurrentGofile(t *testing.T) {

	currentGofileByMethod := GetCurrentGofile()
	expectedPath := os.Getenv("GOPATH") + "\\src\\github.com\\covexo\\devspace\\pkg\\util\\fsutil\\filesystem_test.go"

	expectedPath = strings.Replace(expectedPath, "\\", "/", -1)
	currentGofileByMethod = strings.Replace(currentGofileByMethod, "\\", "/", -1)

	assert.Equal(t, expectedPath, currentGofileByMethod)
}
