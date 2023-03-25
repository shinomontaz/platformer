package activities

type HitBoxOption func(hb *HitBox)

func WithTTL(ttl float64) HitBoxOption {
	return func(hb *HitBox) {
		hb.ttl = ttl
	}
}
