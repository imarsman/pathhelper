package paths

import (
	"fmt"
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

func reverseInt(a []int) []int {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}

	return a
}

type Ordered interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 |
		uint32 | uint64 | uintptr |
		float32 | float64 | string
}

func reverseGeneric[T Ordered](a []T) []T {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}

	return a
}

func TestReverse(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	fmt.Println("input", a)
	a = reverseGeneric(a)
	fmt.Println("reversed input", a)

	sA := []string{"a", "b", "c", "d", "e"}
	fmt.Println("input", sA)
	sA = reverseGeneric(sA)
	fmt.Println("reversed input", sA)
}

// go test -benchmem -run=Bench -bench=Reverse

func BenchmarkReverseType(b *testing.B) {
	a := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		a = reverseInt(a)
	}
}

func BenchmarkReverseGeneric(b *testing.B) {
	a := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		a = reverseGeneric(a)
	}
}
