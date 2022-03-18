package world

import (
	"strconv"

	"github.com/faiface/pixel"

	"platformer/actor"
	"platformer/ai"
	"platformer/common"
	"platformer/config"
	"platformer/events"
	"platformer/factories"

	"github.com/faiface/pixel/pixelgl"
	"github.com/salviati/go-tmx/tmx"
)

type World struct {
	Height float64
	Width  float64
	qtTile *common.Quadtree
	qtObjs *common.Quadtree
	qtPhys *common.Quadtree

	gravity float64

	tm           *tmx.Map
	geom         tmx.ObjectGroup
	scenery      tmx.ObjectGroup
	meta         tmx.Object
	batches      []*pixel.Batch
	batchIndices map[string]int
	sprites      map[string]*pixel.Sprite

	objects     map[int]tmx.Object
	objectTiles map[int]*tmx.DecodedTile
	phys        map[int]tmx.Object
	tiles       map[int]*tmx.DecodedTile

	viewport pixel.Rect
	enmeta   []tmx.Object
	enemies  []*actor.Actor
	hero     *actor.Actor

	visibleTiles []common.Objecter
	visibleObjs  []common.Objecter

	alerts []Alert
}

func New(source string) *World {
	tm, err := tmx.ReadFile(source)
	if err != nil {
		panic(err)
	}

	w := World{
		tm:           tm,
		batches:      make([]*pixel.Batch, 0),
		batchIndices: make(map[string]int),
		sprites:      make(map[string]*pixel.Sprite),
		tiles:        make(map[int]*tmx.DecodedTile),
		phys:         make(map[int]tmx.Object),
		objects:      make(map[int]tmx.Object),

		objectTiles: make(map[int]*tmx.DecodedTile),
		enmeta:      make([]tmx.Object, 0),
		enemies:     make([]*actor.Actor, 0),
	}

	w.init()

	return &w
}

func (w *World) init() {
	for _, og := range w.tm.ObjectGroups {
		if og.Name == "geom" {
			w.geom = og
		}
		if og.Name == "meta" {
			for _, o := range og.Objects {
				if o.Type == "scene" {
					w.meta = o
				}
				if o.Type == "enemy" {
					w.enmeta = append(w.enmeta, o)
				}
			}
		}
		if og.Name == "scenery" {
			w.scenery = og
		}
	}
	w.Height = float64(w.tm.TileHeight * w.tm.Height)
	w.Width = float64(w.tm.TileWidth * w.tm.Width)

	r := pixel.R(0.0, 0.0, w.Width, w.Height)
	w.qtTile = common.New(1, r)
	w.qtPhys = common.New(1, r)
	w.qtObjs = common.New(1, r)

	w.initProps()
	w.initSets()
	w.initTiles()
	w.initPhys()
	w.initObjs()
}

func (w *World) initProps() {
	for _, p := range w.tm.Properties {
		if p.Name == "gravity" {
			if g, err := strconv.ParseFloat(p.Value, 64); err == nil {
				w.gravity = g
			}
		}
	}
}

func (w *World) InitEnemies() {
	for _, o := range w.enmeta {
		w.AddEnemy(o)
	}
}

func (w *World) initSets() {
	batchCounter := 0
	for _, tileset := range w.tm.Tilesets {
		if len(tileset.Tiles) > 0 && tileset.Image.Source == "" {
			for _, tile := range tileset.Tiles {
				if _, alreadyLoaded := w.sprites[tile.Image.Source]; !alreadyLoaded {
					sprite, pictureData := loadSprite(tile.Image.Source)
					w.sprites[tile.Image.Source] = sprite
					w.batches = append(w.batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
					w.batchIndices[tile.Image.Source] = batchCounter
					batchCounter++
				}
			}
		} else {
			if _, alreadyLoaded := w.sprites[tileset.Image.Source]; !alreadyLoaded {
				sprite, pictureData := loadSprite(tileset.Image.Source)
				w.sprites[tileset.Image.Source] = sprite
				w.batches = append(w.batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
				w.batchIndices[tileset.Image.Source] = batchCounter
				batchCounter++
			}
		}
	}
}

func (w *World) initTiles() {
	for _, layer := range w.tm.Layers {
		for tileIndex, tile := range layer.DecodedTiles {
			if tile.Nil {
				continue
			}
			ts := tile.Tileset
			// Calculate the framing for the tile within its tileset's source image
			gamePos := indexToGamePos(tileIndex, w.tm.Width, w.tm.Height)
			pos := gamePos.ScaledXY(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight)))
			w.tiles[tileIndex] = tile
			res := w.qtTile.Insert(common.Objecter{ID: tileIndex, R: pixel.R(pos.X, pos.Y, pos.X+float64(ts.TileWidth), pos.Y+float64(ts.TileHeight))})
			if !res {
				panic("canot insert tile!")
			}
		}
	}
}

func (w *World) initPhys() {
	for _, o := range w.geom.Objects {
		min := pixel.V(
			float64(o.X),
			float64(w.Height)-float64(o.Y)-float64(o.Height),
		)
		max := pixel.Vec{
			X: min.X + float64(o.Width),
			Y: min.Y + float64(o.Height),
		}

		rc := pixel.Rect{
			Min: min,
			Max: max,
		}

		w.phys[o.GID] = o
		w.qtPhys.Insert(common.Objecter{ID: o.GID, R: rc})
	}
}

func (w *World) initObjs() {
	for _, o := range w.scenery.Objects {
		min := pixel.V(
			float64(o.X),
			float64(w.Height)-float64(o.Y),
		)
		max := pixel.Vec{
			X: min.X + float64(o.Width),
			Y: min.Y + float64(o.Height),
		}

		rc := pixel.Rect{
			Min: min,
			Max: max,
		}

		dTile, err := w.tm.DecodeGID(tmx.GID(o.GID))
		if err != nil {
			panic(err) // TODO!
		}
		w.objectTiles[o.GID] = dTile

		w.objects[o.GID] = o
		w.qtObjs.Insert(common.Objecter{ID: o.GID, R: rc})
	}
}

func (w *World) Update(rect pixel.Rect, dt float64) {
	w.viewport = rect
	w.visibleTiles = w.qtTile.Retrieve(rect)
	w.visibleObjs = w.qtObjs.Retrieve(rect)

	if w.hero != nil {
		w.hero.Update(dt)
	}

	ai.Update(dt)
	for _, en := range w.enemies {
		en.Update(dt)
	}

	updateStrikes(dt, w.enemies, w.hero)
	updateSpells(dt, w.enemies, w.hero)
	updateAlerts(dt)
}

func (w *World) GetQt() *common.Quadtree {
	return w.qtPhys
}

func (w *World) GetGravity() float64 {
	return w.gravity
}

func (w *World) AddHero(h *actor.Actor) {
	w.hero = h
}

func (w *World) GetHeroPos() pixel.Vec {
	return w.hero.GetPos()
}

func (w *World) GetHero() *actor.Actor {
	return w.hero
}

func (w *World) GetHeroHp() int {
	return w.hero.GetHp()
}

func (w *World) IsSee(from, to pixel.Vec) bool {
	if !w.viewport.Contains(from) { // check offscreen conditions 1
		return false
	}
	line := pixel.L(from, to)
	if line.Len() > w.viewport.W() { // check offscreen conditions 2
		return false
	}

	box := Box(from, to) // create broadbox
	objs := w.qtPhys.Retrieve(box)

	for _, o := range objs { // check collision for line from -> to against physic layer
		if len(o.R.IntersectionPoints(line)) > 0 {
			return false
		}
	}

	return true
}

func (w *World) AddEnemy(meta tmx.Object) {
	enemy := factories.NewActor(config.Profiles[meta.Name], w)
	enemy.Move(pixel.V(meta.X, w.Height-meta.Y))
	factories.NewAi(config.Profiles[meta.Name].Type, enemy, w)
	w.enemies = append(w.enemies, enemy)
}

func (w *World) AddAlert(pos pixel.Vec, force float64) {
	al := addAlert(pos, force)
	for _, en := range w.enemies {
		alrect := al.GetRect()
		if alrect.Contains(en.GetPos()) {
			a := ai.GetByObj(en)
			if a != nil {
				a.Listen(events.ALERT, alrect.Center())
			}
		}
	}
}

func (w *World) AddStrike(owner *actor.Actor, r pixel.Rect, power int, speed pixel.Vec) {
	AddStrike(owner, r, power, speed)
}

func (w *World) AddSpell(owner *actor.Actor, t pixel.Vec, spell string) {
	AddSpell(owner, t, spell)
}

func (w *World) GetCenter() pixel.Vec {
	rect := pixel.Rect{
		Min: pixel.V(
			float64(w.meta.X),
			w.Height-float64(w.meta.Y)-float64(w.meta.Height),
		),
		Max: pixel.V(
			float64(w.meta.X)+float64(w.meta.Width),
			w.Height-float64(w.meta.Y),
		),
	}
	return rect.Center()
	//	return pixel.V(float64(w.meta.X), w.Height-float64(w.meta.Y)-float64(w.meta.Height))

}

func (w *World) Draw(win *pixelgl.Window) {
	for _, batch := range w.batches {
		batch.Clear()
	}

	for _, t := range w.visibleTiles {
		tile := w.tiles[t.ID]
		ts := tile.Tileset
		tID := int(tile.ID)

		numRows := ts.Tilecount / ts.Columns
		x, y := tileIDToCoord(tID, ts.Columns, numRows)

		iX := float64(x) * float64(ts.TileWidth)
		fX := iX + float64(ts.TileWidth)
		iY := float64(y) * float64(ts.TileHeight)
		fY := iY + float64(ts.TileHeight)

		sprite := w.sprites[ts.Image.Source]
		sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
		pos := t.R.Center()
		sprite.Draw(w.batches[w.batchIndices[ts.Image.Source]], pixel.IM.Moved(pos))
	}

	for _, obj := range w.visibleObjs {
		o := w.objects[obj.ID]
		dTile := w.objectTiles[o.GID]

		ts := dTile.Tileset
		tID := o.GID - int(ts.FirstGID)

		tile := ts.Tiles[tID]

		iX := 0.0
		fX := float64(tile.Image.Width)
		iY := 0.0
		fY := float64(tile.Image.Height)

		sprite := w.sprites[tile.Image.Source]
		sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
		sprite.Draw(w.batches[w.batchIndices[tile.Image.Source]], pixel.IM.Moved(obj.R.Center()))
	}

	for _, batch := range w.batches {
		batch.Draw(win)
	}
	for _, e := range w.enemies {
		e.Draw(win)
	}
	drawStrikes(win)

	drawAlerts(win)
	if w.hero != nil {
		w.hero.Draw(win)
	}

	drawSpells(win)
}
