package dialogs

func actionAllow(name string) {
	w.AllowPortal(name) // w *world.World — глобаль в dialogs/main.go:18, SetWorld уже вызывается
}
