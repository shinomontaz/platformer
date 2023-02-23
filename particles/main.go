package particles

var particles []particle

func Init(maxnum int) {
	particles = make([]particle, maxnum)
}
