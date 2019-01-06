package fcache

import (
	"os"
	"runtime"
	"testing"
	"time"
)

func TestFileTime(t *testing.T) {
	t.Log(runtime.GOOS)
	time.Now().Unix()
	ft, err := GetFileTime(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ctime: %#v", time.Unix(0, ft.CreateTime).String())
	t.Logf("mtime: %#v", time.Unix(0, ft.ModifyTime).String())
	t.Logf("atime: %#v", time.Unix(0, ft.AccessTime).String())
}
