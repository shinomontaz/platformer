package quadtree

import (
	"testing"

	"github.com/shinomontaz/pixel"
)

type mockObject struct {
	id   int
	rect pixel.Rect
}

func (m mockObject) Rect() pixel.Rect {
	return m.rect
}

func (m mockObject) Id() int {
	return m.id
}

func TestRetrieve(t *testing.T) {
	// Create a Quadtree with a bounding box and max items
	qt := New(0, 2, pixel.R(0, 0, 100, 100))

	// Create some mock objects
	obj1 := mockObject{id: 1, rect: pixel.R(10, 10, 20, 20)}
	obj2 := mockObject{id: 2, rect: pixel.R(30, 30, 40, 40)}
	obj3 := mockObject{id: 3, rect: pixel.R(50, 50, 60, 60)}

	// Insert objects into the Quadtree
	qt.Insert(obj1)
	qt.Insert(obj2)
	qt.Insert(obj3)

	// Test case 1: Retrieve objects within a rectangle that intersects with obj1 and obj2
	res := qt.Retrieve(pixel.R(0, 0, 49, 49))
	expected := []int{1, 2}
	if !equal(res, expected) {
		t.Errorf("Expected %v, but got %v", expected, res)
	}

	res = qt.Retrieve(pixel.R(0, 0, 60, 60))
	expected = []int{3, 1, 2}
	if !equal(res, expected) {
		t.Errorf("Expected %v, but got %v", expected, res)
	}

	// Test case 2: Retrieve objects within a rectangle that intersects with obj3
	res = qt.Retrieve(pixel.R(51, 50, 70, 70))
	expected = []int{3}
	if !equal(res, expected) {
		t.Errorf("Expected %v, but got %v", expected, res)
	}

	// Test case 3: Retrieve objects within a rectangle that does not intersect with any objects
	res = qt.Retrieve(pixel.R(51, 0, 90, 49))
	expected = []int{}
	if !equal(res, expected) {
		t.Errorf("Expected %v, but got %v", expected, res)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

/*
var Qtsize = 10

func TestCompare(t *testing.T) {
	rand.Init(1 << 10)

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

func BenchmarkCreate(b *testing.B) {
	rand.Init(1 << 10)
	w := 1000.0
	h := 100000.0
	mainrect := pixel.R(0.0, 0.0, w, h)
	randomsqrs := make([][]pixel.Rect, b.N)
	for i := 0; i < b.N; i++ {
		randomsqrs[i] = createRandomsqares(w, h, 10000)
	}

	b.Run(
		fmt.Sprintf("input_size_%d", b.N),
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				qt := New(1, Qtsize, mainrect)
				for j := 0; j < 10000; j++ {
					qt.Insert(Objecter{R: randomsqrs[i][j]})
				}
			}
		})
}

func BenchmarkRetrieve(b *testing.B) {
	rand.Init(1 << 10)
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

func createRandomsqares(w, h float64, cnt int) []pixel.Rect {
	randomsqrs := make([]pixel.Rect, cnt)
	for j := 0; j < cnt; j++ {
		randomsqrs[j] = rndSqr(w, h)
	}
	return randomsqrs
}

func rndSqr(w, h float64) pixel.Rect {
	mix := rand.Float64() * w
	max := mix + rand.Float64()*(w-mix)
	miy := rand.Float64() * h
	may := miy + rand.Float64()*(h-miy)
	return pixel.R(mix, miy, max, may)
}

*/
