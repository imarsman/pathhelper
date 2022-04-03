package paths

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

// TestPaths do a simple test of fetching paths
// go clean -cache && go test -v .
func TestPaths(t *testing.T) {
	is := is.New(t)

	var cp *pathSet
	var cmp *pathSet

	cp = newPathSet(pathPath, "/etc/paths", "/etc/paths.d", "~/.config/pathhelper/paths.d")
	cmp = newPathSet(manPath, "/etc/manpaths", "/etc/manpaths.d", "~/.config/pathhelper/manpaths.d")

	err := cp.populate()
	is.NoErr(err)

	err = cmp.populate()
	is.NoErr(err)

	var cpStr = strings.Join(cp.paths, ":")
	var cmpStr = strings.Join(cmp.paths, ":")

	t1 := time.Now()
	runs := 1000
	for i := 0; i < runs; i++ {
		cp = newPathSet(pathPath, "/etc/paths", "/etc/paths.d", "~/.config/pathhelper/paths.d")
		cmp = newPathSet(manPath, "/etc/manpaths", "/etc/manpaths.d", "~/.config/pathhelper/manpaths.d")

		err := cp.populate()
		is.NoErr(err)
		err = cmp.populate()
		is.NoErr(err)

		var cpLoopStr = strings.Join(cp.paths, ":")
		var cmpLoopStr = strings.Join(cmp.paths, ":")

		if cpStr != cpLoopStr || cmpStr != cmpLoopStr {
			t.Log("unequal on", i)
			t.Log("cp", cpStr)
			t.Log("cpLoop", cpLoopStr)
			t.Log("cmp", cmpStr)
			t.Log("cmpLoop", cmpLoopStr)
			t.Fail()
			break
		}
		// is.Equal(cpStr, cpLoopStr)
		// is.Equal(cmpStr, cmpLoopStr)

	}
	total := float64(time.Since(t1).Milliseconds())
	t.Logf("total ms to run %d times %.2f ms", runs, total)
	t.Logf("average time in ms to do a load of paths and manpaths %v", total/float64(runs))
	t.Log(cp.paths)
	t.Log(cmp.paths)
}

// go test -benchmem -bench=.
// 0.465517 ms per op
func BenchmarkPathLoad(b *testing.B) {
	is := is.New(b)

	var cp *pathSet
	var cmp *pathSet

	for i := 0; i < b.N; i++ {
		cp = newPathSet(pathPath, "/etc/paths", "/etc/paths.d", "~/.config/pathhelper/paths.d")
		cmp = newPathSet(manPath, "/etc/manpaths", "/etc/manpaths.d", "~/.config/pathhelper/manpaths.d")
		err := cp.populate()
		is.NoErr(err)

		err = cmp.populate()
		is.NoErr(err)
	}

	if b.N == 1 {
		b.Log(cp.paths)
		b.Log(cmp.paths)
	}
}
