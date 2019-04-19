package bro

import (
	"strings"
	"testing"
)

func Test_JoinPath(t *testing.T) {
	testCases := []struct {
		absolutePath string
		relativePath string
		finalPath    string
	}{
		{"/path1/path2", "path3/path4", "/path1/path2/path3/path4"},
		{"path1//////path2////", "/////path3////path4", "path1/path2/path3/path4"},
		{"path1/path2////", "path3/path4///", "path1/path2/path3/path4/"},
	}
	for _, tC := range testCases {
		t.Run(strings.Join([]string{" Test path: ", tC.absolutePath, tC.relativePath}, " "), func(t *testing.T) {
			finalPath := JoinPath(tC.absolutePath, tC.relativePath)
			if finalPath != tC.finalPath {
				t.Fatalf("want %s, get %s", tC.finalPath, finalPath)
			}
		})
	}
}
