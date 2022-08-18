package sftp

import "os"

func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true

}
