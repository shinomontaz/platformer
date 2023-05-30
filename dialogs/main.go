package dialogs

import (
	"encoding/json"
	"platformer/actor"
	"platformer/common"
	"platformer/world"

	"github.com/shinomontaz/pixel/text"
)

var (
	loader    *common.Loader
	dlgs      map[int]Dialog
	atlas     *text.Atlas
	atlasbig  *text.Atlas
	activeDlg *Dialog
	w         *world.World
)

func Init(l *common.Loader) {
	loader = l
	atlas = text.NewAtlas(common.GetFont("regular16"), text.ASCII)
	atlasbig = text.NewAtlas(common.GetFont("regular20"), text.ASCII)

	alldialogs, err := loader.LoadFileToString("dialogs.json")
	if err != nil {
		panic(err)
	}
	dlgst := make([]Dialog, 0)
	dlgs = make(map[int]Dialog)
	err = json.Unmarshal([]byte(alldialogs), &dlgst)
	for _, d := range dlgst {
		dlgs[d.ID] = d
	}
	if err != nil {
		panic(err)
	}
}

func SetWorld(wo *world.World) {
	w = wo
}

func SetActive(id int, a *actor.Actor) {
	dlg, ok := dlgs[id]
	if !ok {
		panic("no such dialog exists!")
	}
	activeDlg = &dlg
	activeDlg.a = a
}

func UnsetActive() {
	activeDlg = nil
}

func GetActive() *Dialog {
	return activeDlg
}

func Get(id int) *Dialog {
	res, ok := dlgs[id]
	if !ok {
		return nil
	}
	return &res
}
