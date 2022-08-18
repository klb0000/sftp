package sftp

import (
	"fmt"
	"testing"
)

func TestPathExists(t *testing.T) {
	path := "files.og"
	fmt.Println(PathExists(path))
}
