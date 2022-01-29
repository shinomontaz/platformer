package animation

type AnimatingConfig interface {
	W() float64
	H() float64
	M() float64
	N() string
	Get() ([]string, []string, [][]int)
	GetGroups() (map[string][]string, map[string][]string, map[string][][]int)
}
