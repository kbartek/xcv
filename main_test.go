package main

import (
	"testing"
)

func Test_filesize(t *testing.T) {
	if filesize(1024) == "1.00K" && filesize(1024*1024) == "1.00M" && filesize(1024*1024*2.5) == "2.50M" {
		t.Log("filesize test passed")
	} else {
		t.Error("filesize test failed")
	}
}
