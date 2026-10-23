package world

func (w *World) UpdateCreatures(dt float64) {
	if w.crtrs == nil {
		return
	}
	w.crtrs.Update(dt, w.visiblePhys, w.visibleSpec)
}
