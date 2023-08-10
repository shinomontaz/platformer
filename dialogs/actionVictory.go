package dialogs

import (
	"platformer/events"
)

func actionVictory(d *Dialog) {
	// TODO!
	d.Notify(events.GAMEVENT_VICTORY)
}
