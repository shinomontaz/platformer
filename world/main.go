package world

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/shinomontaz/pixel"

	"platformer/actor"
	"platformer/background"
	"platformer/common"
	"platformer/creatures"
	"platformer/loot"
	"platformer/particles"
	"platformer/projectiles"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/shinomontaz/pixel/imdraw"
	"github.com/shinomontaz/pixel/pixelgl"

	tmx "github.com/lafriks/go-tiled"
)

type World struct {
	b      *background.Back
	loader *common.Loader

	cnv    *pixelgl.Canvas
	cnv2   *pixelgl.Canvas
	Height float64
	Width  float64
	qtTile *common.Quadtree
	qtObjs *common.Quadtree
	qtPhys *common.Quadtree
	qtSpec *common.Quadtree

	gravity float64

	tm           *tmx.Map
	geom         *tmx.ObjectGroup
	scenery      *tmx.ObjectGroup
	meta         *tmx.Object
	batches      []*pixel.Batch
	batchIndices map[string]int
	sprites      map[string]*pixel.Sprite

	objects     map[uint32]*tmx.Object
	objectTiles map[uint32]*tmx.LayerTile
	//	phys        map[uint32]*tmx.Object
	tiles      map[uint32]*tmx.LayerTile
	imdrawrect *imdraw.IMDraw

	viewport      pixel.Rect
	creaturesmeta []*tmx.Object

	visibleTiles []common.Objecter
	visibleObjs  []common.Objecter
	visiblePhys  []common.Objecter
	visibleSpec  []common.Objecter

	uObjects    []mgl32.Vec4 // = 250 rectangles
	uNumObjects int32
	uLight      mgl32.Vec2

	IsDebug bool
}

var colors = make(map[string]color.Color)

type Option func(*World)

func WithLoader(l *common.Loader) Option {
	return func(w *World) {
		w.loader = l
	}
}

func New(source string, rect pixel.Rect, opts ...Option) (*World, error) {
	w := World{
		batches:      make([]*pixel.Batch, 0),
		batchIndices: make(map[string]int),
		sprites:      make(map[string]*pixel.Sprite),
		tiles:        make(map[uint32]*tmx.LayerTile),
		//		phys:         make(map[uint32]*tmx.Object),
		objects: make(map[uint32]*tmx.Object),

		objectTiles:   make(map[uint32]*tmx.LayerTile),
		creaturesmeta: make([]*tmx.Object, 0),
		viewport:      rect,
		imdrawrect:    imdraw.New(nil),
	}

	for _, opt := range opts {
		opt(&w)
	}

	mapsource, err := w.loader.Open(source)
	if err != nil {
		return nil, err
	}
	defer mapsource.Close()

	tm, err := tmx.LoadReader("assets", mapsource)
	if err != nil {
		return nil, fmt.Errorf("%s 111", err)
	}

	w.tm = tm

	w.init()

	return &w, nil
}

func (w *World) init() {
	w.Height = float64(w.tm.TileHeight * w.tm.Height)
	w.Width = float64(w.tm.TileWidth * w.tm.Width)

	r := pixel.R(0.0, 0.0, w.Width, w.Height)
	w.qtTile = common.New(1, 10, r)
	w.qtPhys = common.New(1, 10, r)
	w.qtObjs = common.New(1, 10, r)
	w.qtSpec = common.New(1, 10, r)

	for _, og := range w.tm.ObjectGroups {
		if og.Name == "geom" {
			w.geom = og
		}
		if og.Name == "meta" {
			for _, o := range og.Objects {
				if o.Class == "hero" {
					w.meta = o
				}

				if o.Class == "enemy" || o.Class == "npc" || o.Class == "coin" {
					o.Y = w.Height - o.Y
					w.creaturesmeta = append(w.creaturesmeta, o)
				}

				if o.Class == "water" {
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

					w.qtSpec.Insert(common.Objecter{ID: o.GID, R: rc, Type: common.WATER})

					w.drawRectIfNeeded(o, rc)
				}
			}
		}
		if og.Name == "scenery" {
			w.scenery = og
		}
	}

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

	w.viewport = w.viewport.Moved(rect.Center().Sub(pixel.V(w.viewport.W()/2, w.viewport.H()/2)))

	var wg sync.WaitGroup
	w.initProps()
	go func() {
		w.initSets()
		wg.Add(1)
		defer wg.Done()
	}()
	go func() {
		w.initTiles()
		wg.Add(1)
		defer wg.Done()
	}()
	go func() {
		w.initPhys()
		wg.Add(1)
		defer wg.Done()
	}()
	go func() {
		w.initObjs()
		wg.Add(1)
		defer wg.Done()
	}()
	go func() {
		w.initShader("shader/world.glsl")
		wg.Add(1)
		defer wg.Done()
	}()
	// go func() {
	// 	w.initEnemies()
	// 	w.initNpcs()
	// 	wg.Add(1)
	// 	defer wg.Done()
	// }()
	wg.Wait()
}

func (w *World) initShader(shadername string) {
	w.cnv = pixelgl.NewCanvas(pixel.R(0, 0, w.viewport.W(), w.viewport.H()))
	w.cnv.SetSmooth(true)

	w.cnv2 = pixelgl.NewCanvas(pixel.R(0, 0, w.viewport.W(), w.viewport.H()))
	w.cnv2.SetSmooth(true)

	w.uObjects = make([]mgl32.Vec4, 0)
	w.uLight = [2]float32{float32(-w.viewport.W()/2 + 100.0), float32(w.viewport.H()/2 - 100.0)}

	w.cnv.SetUniform("uLight", &w.uLight)
	w.cnv.SetUniform("uObjects", &w.uObjects)
	w.cnv.SetUniform("uNumObjects", &w.uNumObjects)

	fragSource, err := w.loader.LoadFileToString(shadername)
	if err != nil {
		panic(err)
	}

	w.cnv.SetFragmentShader(fragSource)
}

func (w *World) initProps() {
	if w.tm.Properties == nil {
		return
	}
	w.gravity = w.tm.Properties.GetFloat("gravity")
}

func (w *World) GetMetas() []*tmx.Object {
	return w.creaturesmeta
}

func (w *World) initSets() {
	batchCounter := 0
	for _, tileset := range w.tm.Tilesets {
		if len(tileset.Tiles) > 0 { // tileset of pictures
			for _, tile := range tileset.Tiles {
				if _, alreadyLoaded := w.sprites[tile.Image.Source]; !alreadyLoaded {
					pictureData, _ := w.loader.LoadPicture(tile.Image.Source)
					sprite := pixel.NewSprite(pictureData, pictureData.Bounds())
					w.sprites[tile.Image.Source] = sprite
					w.batches = append(w.batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
					w.batchIndices[tile.Image.Source] = batchCounter
					batchCounter++
				}
			}
		} else {
			if _, alreadyLoaded := w.sprites[tileset.Image.Source]; !alreadyLoaded {
				pictureData, _ := w.loader.LoadPicture(tileset.Image.Source)
				sprite := pixel.NewSprite(pictureData, pictureData.Bounds())
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
		for tileIndex, tile := range layer.Tiles {
			if tile.Nil {
				continue
			}
			ts := tile.Tileset

			tile.ID = ts.FirstGID + tile.ID - 1

			// Calculate the framing for the tile within its tileset's source image
			gamePos := indexToGamePos(tileIndex, w.tm.Width, w.tm.Height)
			pos := gamePos.ScaledXY(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight)))
			w.tiles[tile.ID] = tile
			res := w.qtTile.Insert(common.Objecter{ID: tile.ID, R: pixel.R(pos.X, pos.Y, pos.X+float64(ts.TileWidth), pos.Y+float64(ts.TileHeight))})
			if !res {
				panic("cannot insert tile!")
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
		phType := common.GROUND
		if o.Class == "barier" {
			phType = common.BARRIER
		}
		//		w.phys[o.GID] = o
		w.qtPhys.Insert(common.Objecter{ID: o.GID, R: rc, Type: phType})

		w.drawRectIfNeeded(o, rc)
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

		dTile, err := w.tm.TileGIDToTile(o.GID) // type(dTile) = *tmx.LayerTile
		if err != nil {
			panic(err) // TODO!
		}
		w.objectTiles[o.GID] = dTile

		w.objects[o.GID] = o
		w.qtObjs.Insert(common.Objecter{ID: o.GID, R: rc})
	}
}

func (w *World) GetViewport() pixel.Rect {
	return w.viewport
}

func (w *World) GetVisiblePhys() []common.Objecter {
	return w.visiblePhys
}

func (w *World) Update(rect pixel.Rect, dt float64) {
	w.viewport = rect

	w.visibleTiles = w.qtTile.Retrieve(w.viewport)
	w.visibleObjs = w.qtObjs.Retrieve(w.viewport)
	w.visibleSpec = w.qtSpec.Retrieve(w.viewport)

	w.visiblePhys = w.qtPhys.Retrieve(w.viewport)
	w.uObjects = make([]mgl32.Vec4, 0)
	c := w.viewport.Center()
	for _, o := range w.visiblePhys {
		if o.Type != common.GROUND {
			continue
		}
		uObject := mgl32.Vec4{float32(-c.X + o.R.Min.X), float32(-c.Y + o.R.Min.Y - 150), float32(-c.X + o.R.Max.X), float32(-c.Y + o.R.Max.Y - 150)}
		w.uObjects = append(w.uObjects, uObject)
	}
	w.uNumObjects = int32(len(w.uObjects))

	creatures.Update(dt, w.visiblePhys, w.visibleSpec)
	loot.Update(dt, w.visiblePhys)
	particles.Update(dt, w.visiblePhys)
	projectiles.Update(dt, w.visiblePhys, w.visibleSpec)
}

func (w *World) GetGravity() float64 {
	return w.gravity
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

func (w *World) AddSpell(owner *actor.Actor, t pixel.Vec, spell string, objs []common.Objecter) {
	AddSpell(owner, t, spell, objs)
}

func (w *World) GetCenter() pixel.Vec {
	return w.viewport.Center()
}

func (w *World) SetBackground(b *background.Back) {
	w.b = b
}

//func (w *World) Draw(win *pixelgl.Window, hpos pixel.Vec, cam pixel.Vec) {
func (w *World) Draw(t pixel.Target, hpos pixel.Vec, cam pixel.Vec, center pixel.Vec) {

	w.cnv.Clear(color.RGBA{0, 0, 0, 1})
	w.cnv2.Clear(color.RGBA{240, 248, 255, 1})

	w.cnv2.SetMatrix(pixel.IM.Moved(w.cnv2.Bounds().Center().Sub(cam)))

	w.b.Draw(w.cnv2, hpos, cam)

	w.imdrawrect.Draw(w.cnv2)

	for _, batch := range w.batches {
		batch.Clear()
	}

	for _, t := range w.visibleTiles {
		tile := w.tiles[t.ID]
		ts := tile.Tileset
		tID := int(tile.ID - ts.FirstGID + 1)

		numRows := ts.TileCount / ts.Columns
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

		tile, _ := dTile.Tileset.GetTilesetTile(dTile.ID)

		iX := 0.0
		fX := float64(tile.Image.Width)
		iY := 0.0
		fY := float64(tile.Image.Height)

		sprite := w.sprites[tile.Image.Source]
		sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
		flip := 1.0
		if dTile.HorizontalFlip {
			flip = -1.0
		}
		sprite.Draw(w.batches[w.batchIndices[tile.Image.Source]], pixel.IM.ScaledXY(pixel.ZV, pixel.V(flip, 1)).Moved(obj.R.Center()))
	}

	for _, batch := range w.batches {
		batch.Draw(w.cnv2)
	}

	if w.IsDebug {
		imd := imdraw.New(nil)
		imd.Color = color.RGBA{255, 0, 0, 1}
		for _, p := range w.visiblePhys {
			vertices := p.R.Vertices()
			for _, v := range vertices {
				imd.Push(v)
			}
			imd.Rectangle(1)
		}
		imd.Draw(w.cnv2)
	}

	drawSpells(w.cnv2)

	creatures.Draw(w.cnv2)
	loot.Draw(w.cnv2)
	particles.Draw(w.cnv2)
	projectiles.Draw(w.cnv2)

	w.cnv2.Draw(w.cnv, pixel.IM.Moved(w.cnv.Bounds().Center()))
	w.cnv.Draw(t, pixel.IM.Moved(center))
}

func (w *World) drawRectIfNeeded(o *tmx.Object, r pixel.Rect) {
	var col color.Color
	var err error
	var ok bool
	cstr := o.Properties.GetString("color")
	if cstr == "" {
		return
	}

	if col, ok = colors[cstr]; !ok {
		col, err = ParseHexColor(cstr)
		if err != nil {
			panic(err)
		}
		colors[cstr] = col
	}
	w.imdrawrect.Color = col
	w.imdrawrect.Push(r.Min, r.Max)
	w.imdrawrect.Rectangle(0)
}
