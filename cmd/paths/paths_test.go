package paths

import (
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

	cp = newPathSet(pathTypePath, systemPathFile, systemPathDir, userPathDir)
	cmp = newPathSet(pathTypeManPath, systemManPathFile, systemManPathDir, userManPathDir)

	err := cp.populate()
	is.NoErr(err)

	err = cmp.populate()
	is.NoErr(err)

	var cpStr = cp.zshFormat()
	var cmpStr = cmp.zshFormat()

	t1 := time.Now()
	runs := 1000
	for i := 0; i < runs; i++ {
		cp = newPathSet(pathTypePath, systemPathFile, systemPathDir, userPathDir)
		cmp = newPathSet(pathTypeManPath, systemManPathFile, systemManPathDir, userManPathDir)

		err := cp.populate()
		is.NoErr(err)
		err = cmp.populate()
		is.NoErr(err)

		var cpLoopStr = cp.zshFormat()
		var cmpLoopStr = cmp.zshFormat()

		// Check for inconsistent results which would indicate a failure to produce identical results each time
		if cpStr != cpLoopStr || cmpStr != cmpLoopStr {
			t.Log("unequal on", i)
			t.Log("cp", cpStr)
			t.Log("cpLoop", cpLoopStr)
			t.Log("cmp", cmpStr)
			t.Log("cmpLoop", cmpLoopStr)
			t.Fail()
			break
		}
		is.Equal(cpStr, cpLoopStr)
		is.Equal(cmpStr, cmpLoopStr)

	}

	total := float64(time.Since(t1).Milliseconds())
	t.Logf("total ms to run %d times %.2f ms", runs, total)
	t.Logf("average time in ms to do a load of paths and manpaths %v", total/float64(runs))
	t.Log(cp.zshFormat())
	t.Log(cmp.zshFormat())
}

// BenchmarkPathLoad do benchmark of path load - also helps check for concurrency issues
// go test -benchmem -bench=.
// 0.432291 ms per op
func BenchmarkPathLoad(b *testing.B) {
	is := is.New(b)

	var cp *pathSet
	var cmp *pathSet

	cp = newPathSet(pathTypePath, systemPathFile, systemPathDir, userPathDir)
	cmp = newPathSet(pathTypeManPath, systemManPathFile, systemManPathDir, userManPathDir)

	for i := 0; i < b.N; i++ {
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
