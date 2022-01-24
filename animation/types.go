package animation

type AnimatingConfig interface {
	W() float64
	H() float64
	M() float64
	N() string
	Get() ([]string, []string, [][]int)
}
