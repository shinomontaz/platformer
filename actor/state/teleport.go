package state

import (
	"math"
	"platformer/common/rand"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

// Drawer — состояние рисует себя само (вместо одного спрайта)
type Drawer interface {
	Draw(t pixel.Target, rect pixel.Rect, flip float64)
}

type shard struct {
	spr  *pixel.Sprite
	off  pixel.Vec // центр куска относительно центра кадра
	dir  pixel.Vec // откуда прилетает
	dist float64
	rot  float64
}

type Teleport struct {
	Common
	time   float64
	limit  float64
	sprite *pixel.Sprite
	frame  pixel.Rect
	shards []shard
}

func NewTeleport(a Actor, an Animater) *Teleport {
	return &Teleport{
		Common: Common{
			id:    TELEPORT,
			a:     a,
			anims: an,
			trs:   a.GetTransition(TELEPORT),
		},
		limit: 0.5,
	}
}

func (s *Teleport) Start() {
	s.time = 0.0
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	// собираемся в позу stand — дальше переход в STAND без скачка
	pic, frame := s.anims.GetSprite("stand", 0)
	s.sprite.Set(pic, frame)
	s.frame = frame

	fc := frame.Center()
	s.shards = s.shards[:0]
	for _, r := range splitRect(frame, 6+intn(10)) {
		off := r.Center().Sub(fc)
		dir := pixel.V(rand.Float64()-0.5, rand.Float64()-0.5).Unit()
		if off.Len() > 0.01 {
			dir = off.Unit().Add(dir.Scaled(0.6)).Unit()
		}
		s.shards = append(s.shards, shard{
			spr:  pixel.NewSprite(pic, r),
			off:  off,
			dir:  dir,
			dist: 15 + rand.Float64()*35,
			rot:  (rand.Float64() - 0.5) * math.Pi,
		})
	}
}

func (s *Teleport) Busy() bool { return true }

func (s *Teleport) Listen(e int, v *pixel.Vec) {
	s.a.Inform(events.BUSY)
	// команды контроллера игнорируем — поза не меняется
}

func (s *Teleport) Update(dt float64) {
	if s.time >= s.limit {
		s.a.SetState(STAND)
		return
	}
	s.time += dt
}

func (s *Teleport) GetSprite() *pixel.Sprite {
	return s.sprite
}

func (s *Teleport) Draw(t pixel.Target, rect pixel.Rect, flip float64) {
	q := math.Min(s.time/s.limit, 1)
	f := math.Pow(1-q, 3) // 1 -> 0, тормозит у цели
	scale := pixel.V(rect.W()/s.frame.W()*flip, rect.H()/s.frame.H())
	world := pixel.IM.ScaledXY(pixel.ZV, scale).Moved(rect.Center())
	for _, sh := range s.shards {
		m := pixel.IM.
			Rotated(pixel.ZV, sh.rot*f).
			Moved(sh.off.Add(sh.dir.Scaled(sh.dist * f))).
			Chained(world)
		sh.spr.DrawColorMask(t, m, pixel.Alpha(q))
	}
}

// intn: 0..n-1 через Float64() из [0, 1)
func intn(n int) int {
	return int(rand.Float64() * float64(n))
}

func splitRect(b pixel.Rect, n int) []pixel.Rect {
	if n < 4 {
		n = 4
	}
	minPiece := math.Max(math.Min(b.W(), b.H())/8, 3)
	c := b.Center().Add(pixel.V(
		(0.5-rand.Float64())*b.W()/2, // не picCenter.X: кадр лежит не в (0,0) листа
		(0.5-rand.Float64())*b.H()/2,
	))
	quads := [4]pixel.Rect{
		pixel.R(b.Min.X, c.Y, c.X, b.Max.Y), // top-left
		pixel.R(c.X, c.Y, b.Max.X, b.Max.Y), // top-right
		pixel.R(b.Min.X, b.Min.Y, c.X, c.Y), // bottom-left
		pixel.R(c.X, b.Min.Y, b.Max.X, c.Y), // bottom-right
	}
	counts := [4]int{1, 1, 1, 1}
	for rest := n - 4; rest > 0; rest-- {
		counts[intn(4)]++
	}
	res := make([]pixel.Rect, 0, n)
	for i, q := range quads {
		res = append(res, split(q, counts[i], minPiece)...)
	}
	return res
}

func split(r pixel.Rect, n int, minPiece float64) []pixel.Rect {
	canCutX := r.W() >= 2*minPiece
	canCutY := r.H() >= 2*minPiece
	if n <= 1 || (!canCutX && !canCutY) {
		return []pixel.Rect{r}
	}

	cutX := canCutX && (!canCutY || intn(2) == 0)
	var a, b pixel.Rect
	if cutX {
		c := r.Min.X + minPiece + rand.Float64()*(r.W()-2*minPiece)
		a = pixel.R(r.Min.X, r.Min.Y, c, r.Max.Y)
		b = pixel.R(c, r.Min.Y, r.Max.X, r.Max.Y)
	} else {
		c := r.Min.Y + minPiece + rand.Float64()*(r.H()-2*minPiece)
		a = pixel.R(r.Min.X, c, r.Max.X, r.Max.Y)
		b = pixel.R(r.Min.X, r.Min.Y, r.Max.X, c)
	}

	na := 1 + intn(n-1) // 1..n-1, хотя бы по куску в каждой половине
	return append(split(a, na, minPiece), split(b, n-na, minPiece)...)
}
