package ai

import (
	"github.com/shinomontaz/pixel"
)

// type Actor interface {
// 	AddListener(s common.Subscriber)
// 	GetPos() pixel.Rect
// 	GetDir() int
// 	SetTarget(pv pixel.Vec)
// 	SetSkill(s skills.Skill)
// 	GetSkills() []skills.Skill
// 	Listen(e int, v pixel.Vec)
// }

type Stater interface {
	Update(dt float64)
	Start(poi pixel.Vec)
	IsAlerted() bool
	Listen(e int, v pixel.Vec)
}

type Worlder interface {
	IsSee(from, to pixel.Vec) bool
	//	AddAlert(place pixel.Vec, raduis float64)
}

type Alerter interface {
	GetRect() pixel.Rect
}
