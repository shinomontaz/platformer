package state

func New(code int, a Actor, an Animater) Stater {
	var st Stater
	switch code {
	case STAND:
		st = NewStand(a, an)
	case IDLE:
		st = NewIdle(a, an)
	case WALK:
		st = NewWalk(a, an)
	case INTERACT:
		st = NewInteract(a, an)
	case RUN:
		st = NewRun(a, an)
	case JUMP:
		st = NewJump(a, an)
	case FALL:
		st = NewFall(a, an)
	case HIT:
		st = NewHit(a, an)
	case DEAD:
		st = NewDead(a, an)
	case DEADSUNK:
		st = NewDeadSunk(a, an)
	case ATTACK:
		st = NewAttack(a, an)
	case MELEE:
		st = NewMelee(a, an)
	case MELEEMOVE:
		st = NewMeleemove(a, an)
	case CAST:
		st = NewCast(a, an)
	case RANGED:
		st = NewRanged(a, an)
	}
	return st
}
