package common

import (
	"math/rand"
	"testing"
	"time"

	"github.com/shinomontaz/pixel"
)

var Qtsize = 10

func BenchmarkCreate(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	w := 1000.0
	h := 100000.0
	mainrect := pixel.R(0.0, 0.0, w, h)
	randomsqrs := make([][]pixel.Rect, b.N)
	for i := 0; i < b.N; i++ {
		randomsqrs[i] = createRandomsqares(w, h, 10000)
	}

	for i := 0; i < b.N; i++ {
		qt := New(1, Qtsize, mainrect)
		for j := 0; j < 10000; j++ {
			qt.Insert(Objecter{R: randomsqrs[i][j]})
		}
	}
}

func BenchmarkRetrieve(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	w := 1000.0
	h := 100000.0
	mainrect := pixel.R(0.0, 0.0, w, h)
	cnt := 100000
	randomsqrs := createRandomsqares(w, h, cnt)
	testsqrs := createRandomsqares(w, h, b.N)

	qt := New(1, Qtsize, mainrect)
	for j := 0; j < cnt; j++ {
		qt.Insert(Objecter{R: randomsqrs[j]})
	}

	for i := 0; i < b.N; i++ {
		qt.Retrieve(testsqrs[i])
	}
}

func TestRetrieve(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	w := 1000.0
	h := 100000.0
	mainrect := pixel.R(0.0, 0.0, w, h)
	cnt := 100000
	testcnt := 1000
	randomsqrs := createRandomsqares(w, h, cnt)
	testsqrs := createRandomsqares(w, h, testcnt)

	qt := New(1, Qtsize, mainrect)
	for j := 0; j < cnt; j++ {
		qt.Insert(Objecter{R: randomsqrs[j]})
	}

	var last time.Time
	var dt, mean float64
	times := make([]float64, testcnt)
	for i := 0; i < testcnt; i++ {
		last = time.Now()
		qt.Retrieve(testsqrs[i])
		dt = time.Since(last).Seconds()

		times[i] = dt
	}

	for _, t := range times {
		mean += t
	}

	mean /= float64(len(times))

	t.Logf("Mean retrieve %f", mean)
}

func TestCompare(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	w := 1000.0
	h := 100000.0
	mainrect := pixel.R(0.0, 0.0, w, h)
	cnt := 100000
	testcnt := 1000
	randomsqrs := createRandomsqares(w, h, cnt)
	testsqrs := createRandomsqares(w, h, testcnt)
	testsqrs2 := make([]pixel.Rect, testcnt)

	for i := 0; i < len(testsqrs); i++ {
		testsqrs2[i] = rndSqr(testsqrs[i].W(), testsqrs[i].H()).Moved(testsqrs[i].Min)
	}

	qt := New(1, Qtsize, mainrect)
	for j := 0; j < cnt; j++ {
		qt.Insert(Objecter{R: randomsqrs[j]})
	}

	diff := 0.0
	times := make([]float64, testcnt)
	var last time.Time
	var dt float64
	for i := 0; i < testcnt; i++ {
		last = time.Now()
		qt.Retrieve(testsqrs[i])
		qt.Retrieve(testsqrs2[i])
		dt = time.Since(last).Seconds()
		times[i] = dt

		last = time.Now()
		res := qt.Retrieve(testsqrs[i])
		qt2 := New(1, Qtsize, testsqrs[i])
		for _, r := range res {
			qt2.Insert(Objecter{R: r.R})
		}
		qt2.Retrieve(testsqrs2[i])
		dt = time.Since(last).Seconds()

		times[i] -= dt
	}

	for _, t := range times {
		diff += t
	}

	diff /= float64(len(times))

	t.Logf("diff mean %f", diff)
}

func createRandomsqares(w, h float64, cnt int) []pixel.Rect {
	randomsqrs := make([]pixel.Rect, cnt)
	for j := 0; j < cnt; j++ {
		randomsqrs[j] = rndSqr(w, h)
	}
	return randomsqrs
}

func rndSqr(w, h float64) pixel.Rect {
	mix := common.GetRandFloat() * w
	max := mix + common.GetRandFloat()*(w-mix)
	miy := common.GetRandFloat() * h
	may := miy + common.GetRandFloat()*(h-miy)
	return pixel.R(mix, miy, max, may)
}
