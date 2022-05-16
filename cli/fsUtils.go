package cli

import (
	"os"
)

func CheckIfFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}  
		
	return false
	
}