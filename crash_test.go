package msgpack_test

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/shamaton/msgpack/v2"
)

var crashDir = filepath.Join("testdata", "crashers")

func TestCrashBinary(t *testing.T) {
	entries, err := os.ReadDir(crashDir)
	if err != nil {
		t.Fatalf("os.ReadDir error. err: %+v", err)
	}

	ch := make(chan string, len(entries))

	// worker
	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go check(t, &wg, ch)
	}

	for _, entry := range entries {
		ch <- filepath.Join(crashDir, entry.Name())
	}
	close(ch)
	wg.Wait()
}

func check(t *testing.T, wg *sync.WaitGroup, ch <-chan string) {
	t.Helper()
	defer wg.Done()

	for {
		path, ok := <-ch // closeされると ok が false になる
		if !ok {
			return
		}

		file, err := os.Open(path)
		if err != nil {
			t.Logf("%s file open error. err: %+v", path, err)
			t.Fail()
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			t.Logf("%s io.ReadAll error. err: %+v", path, err)
			t.Fail()
			return
		}

		var r interface{}
		err = msgpack.Unmarshal(data, &r)
		if err == nil {
			t.Logf("err should be occurred.\nname: %s\nlen: %d\nbin: % x",
				file.Name(), len(data), data,
			)
			t.Fail()
			return
		}
	}
}
