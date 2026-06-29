package world

import (
	"image/color"
	"sync"

	"github.com/shinomontaz/pixel"

	"platformer/actor"
	"platformer/common"
	"platformer/common/quadtree"
	"platformer/creatures"
	"platformer/inventory"
	"platformer/loot"
	"platformer/objects"
	"platformer/particles"
	"platformer/projectiles"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/shinomontaz/pixel/imdraw"
	"github.com/shinomontaz/pixel/pixelgl"

	tmx "github.com/lafriks/go-tiled"
)

type World struct {
	r      pixel.Rect
	b      Background
	loader *common.Loader

	cnv    *pixelgl.Canvas
	cnv2   *pixelgl.Canvas
	Height float64
	Width  float64
	qtPhys *quadtree.Qt
	qtTile *quadtree.Qt
	qtObjs *quadtree.Qt
	qtSpec *quadtree.Qt

	gravity float64

	tm *tmx.Map

	objMap map[groupType]*tmx.ObjectGroup

	meta pixel.Vec // here is a center of the map

	batches      []*pixel.Batch
	batchIndices map[string]int
	sprites      map[string]*pixel.Sprite

	specobjs    []object
	objects     []object
	objectTiles []*tmx.LayerTile
	physobjs    []object

	tiles      []tile
	imdrawrect *imdraw.IMDraw

	viewport      pixel.Rect
	creaturesmeta []*tmx.Object
	crtrs         *creatures.List

	visibleTiles []tile
	visibleObjs  []object
	visiblePhys  []common.O
	visibleSpec  []common.O

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

// func WithShader(l *common.Loader) Option {
// 	return func(w *World) {
// 		w.loader = l
// 	}
// }

func New(source string, rect pixel.Rect, opts ...Option) (*World, error) {
	w := World{
		batches:      []*pixel.Batch{},
		batchIndices: map[string]int{},
		sprites:      map[string]*pixel.Sprite{},
		tiles:        []tile{},
		objects:      []object{},
		physobjs:     []object{},
		specobjs:     []object{},

		objectTiles: []*tmx.LayerTile{},

		visibleTiles: []tile{},
		visiblePhys:  []common.O{},
		visibleObjs:  []object{},
		visibleSpec:  []common.O{},

		viewport:   rect,
		imdrawrect: imdraw.New(nil),
	}

	for _, opt := range opts {
		opt(&w)
	}

	var err error
	w.tm, err = w.loadMap(source)
	if err != nil {
		return nil, err
	}
	w.objMap = w.readObjects(w.tm)
	w.init()

	return &w, nil
}

func (w *World) init() {
	w.Height = float64(w.tm.TileHeight * w.tm.Height)
	w.Width = float64(w.tm.TileWidth * w.tm.Width)

	w.r = pixel.R(0.0, 0.0, w.Width, w.Height)
	w.qtTile = quadtree.New(1, 10, w.r)
	w.qtObjs = quadtree.New(1, 10, w.r)
	w.qtSpec = quadtree.New(1, 10, w.r)

	w.meta = pixel.Vec{X: w.tm.Properties.GetFloat("centerX"), Y: w.tm.Properties.GetFloat("centerY")}

	rect := pixel.Rect{
		Min: pixel.V(
			float64(w.meta.X),
			w.Height-float64(w.meta.Y),
		),
		Max: pixel.V(
			float64(w.meta.X),
			w.Height-float64(w.meta.Y),
		),
	}

	w.viewport = w.viewport.Moved(rect.Center().Sub(pixel.V(w.viewport.W()/2, w.viewport.H()/2)))

	w.initProps()
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		w.initSets()
		defer wg.Done()
	}()
	go func() {
		w.initTiles()
		defer wg.Done()
	}()
	go func() {
		w.initPhys()
		defer wg.Done()
	}()
	go func() {
		w.initObjs()
		defer wg.Done()
	}()
	go func() {
		w.initShader("shader/world.glsl")
		defer wg.Done()
	}()
	wg.Wait()

	if og := w.objMap[CREATURES]; og != nil {
		for _, o := range og.Objects {
			if o.Class == "enemy" || o.Class == "npc" || o.Class == "coin" || o.Class == "object" {
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
				oo := object{
					ID:     len(w.specobjs),
					R:      rc,
					Center: rc.Center(),
					T:      common.WATER,
				}
				w.specobjs = append(w.specobjs, oo)
				w.qtSpec.Insert(oo)
				w.drawRectIfNeeded(o, rc)
			}
		}
	}
}

func (w *World) initShader(shadername string) {
	w.cnv = pixelgl.NewCanvas(pixel.R(0, 0, w.viewport.W(), w.viewport.H()))
	w.cnv.SetSmooth(true)

	w.cnv2 = pixelgl.NewCanvas(pixel.R(0, 0, w.viewport.W(), w.viewport.H()))
	w.cnv2.SetSmooth(true)

	w.uObjects = make([]mgl32.Vec4, 0)
	w.uLight = [2]float32{float32(-w.viewport.W()/2 + 100.0), float32(w.viewport.H()/2 - 100.0)}

	// w.cnv.SetUniform("uLight", &w.uLight)
	// w.cnv.SetUniform("uObjects", &w.uObjects)
	// w.cnv.SetUniform("uNumObjects", &w.uNumObjects)

	//	fragSource, err := w.loader.LoadFileToString(shadername)
	// if err != nil {
	// 	panic(err)
	// }

	//	w.cnv.SetFragmentShader(fragSource)
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
				if tile.Image == nil {
					continue
				}
				if _, alreadyLoaded := w.sprites[tile.Image.Source]; !alreadyLoaded {
					pictureData, err := w.loader.LoadPicture(tile.Image.Source)
					if err != nil {
						panic(err)
					}
					sprite := pixel.NewSprite(pictureData, pictureData.Bounds())
					w.sprites[tile.Image.Source] = sprite
					w.batches = append(w.batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
					w.batchIndices[tile.Image.Source] = batchCounter
					batchCounter++
				}
			}
		} else {
			if _, alreadyLoaded := w.sprites[tileset.Image.Source]; !alreadyLoaded {
				pictureData, err := w.loader.LoadPicture(tileset.Image.Source)
				if err != nil {
					panic(err)
				}
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
	var ts *tmx.Tileset
	for _, layer := range w.tm.Layers {
		for tileIndex, tmxTile := range layer.Tiles {
			if tmxTile.Nil {
				continue
			}
			ts = tmxTile.Tileset

			// Calculate the framing for the tile within its tileset's source image
			gamePos := indexToGamePos(tileIndex, w.tm.Width, w.tm.Height)
			pos := gamePos.ScaledXY(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight)))

			tile := tile{
				Index:   len(w.tiles),
				ID:      ts.FirstGID + tmxTile.ID - 1,
				R:       pixel.R(pos.X, pos.Y, pos.X+float64(ts.TileWidth), pos.Y+float64(ts.TileHeight)),
				Tileset: ts,
			}
			w.tiles = append(w.tiles, tile)
			if !w.qtTile.Insert(tile) {
				panic("cannot insert tile!")
			}
		}
	}
}

func (w *World) initObjs() {
	og := w.objMap[SCENERY]
	if og == nil {
		return
	}
	for _, o := range og.Objects {
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
			panic(err)
		}

		oo := object{
			ID:     len(w.objects),
			R:      rc,
			Center: rc.Center(),
		}

		w.objectTiles = append(w.objectTiles, dTile)
		w.objects = append(w.objects, oo)
		w.qtObjs.Insert(oo)
	}
}

func (w *World) GetViewport() pixel.Rect {
	return w.viewport
}

func (w *World) Update(rect pixel.Rect, hpos pixel.Vec, dt float64) {
	w.viewport = rect

	w.updateTiles()
	w.updateObjects()
	w.updateSpecs()

	w.updatePhys()

	w.UpdateCreatures(dt)
	loot.Update(dt, w.visiblePhys)

	collectedloot := loot.Collect(w.crtrs.GetHero())
	inventory.Add(collectedloot)

	particles.Update(dt, w.visiblePhys)
	projectiles.Update(dt, w.visiblePhys, w.visibleSpec)

	objects.Update(dt, w.visiblePhys)

	w.b.Update(dt, hpos)
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
		if len(w.physobjs[o].R.IntersectionPoints(line)) > 0 {
			return false
		}
	}

	return true
}

func (w *World) AddSpell(owner *actor.Actor, t pixel.Vec, spell string, objs []common.O) {
	AddSpell(owner, t, spell, objs)
}

func (w *World) GetCenter() pixel.Vec {
	return w.viewport.Center()
}

func (w *World) SetBackground(b Background) {
	w.b = b
}

func (w *World) SetCreatures(c *creatures.List) {
	w.crtrs = c
}

// func (w *World) Draw(win *pixelgl.Window, hpos pixel.Vec, cam pixel.Vec) {
func (w *World) Draw(t pixel.Target, hpos pixel.Vec, cam pixel.Vec, center pixel.Vec) {
	w.cnv.Clear(color.RGBA{0, 0, 0, 1})
	//	w.cnv2.Clear(color.RGBA{240, 248, 255, 1})
	w.cnv2.Clear(color.RGBA{0, 0, 0, 1})

	w.cnv2.SetMatrix(pixel.IM.Moved(w.cnv2.Bounds().Center().Sub(cam)))

	w.b.Draw(w.cnv2)

	w.imdrawrect.Draw(w.cnv2) // filled rectangles

	for _, batch := range w.batches {
		batch.Clear()
	}

	for _, tile := range w.visibleTiles {
		ts := tile.Tileset
		tID := int(tile.ID - ts.FirstGID + 1)

		numRows := ts.TileCount / ts.Columns
		x, y := tileIDToCoord(tID, ts.Columns, numRows)

		iX := float64(x) * float64(ts.TileWidth)
		fX := iX + float64(ts.TileWidth)
		iY := float64(y) * float64(ts.TileHeight)
		fY := iY + float64(ts.TileHeight)

		sprite := w.sprites[ts.Image.Source]
		if sprite == nil {
			continue
		}
		batchIdx, ok := w.batchIndices[ts.Image.Source]
		if !ok {
			continue
		}

		sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
		pos := tile.R.Center()
		sprite.Draw(w.batches[batchIdx], pixel.IM.Moved(pos))
	}

	for _, obj := range w.visibleObjs {
		dTile := w.objectTiles[obj.ID]

		tile, err := dTile.Tileset.GetTilesetTile(dTile.ID)
		if err != nil || tile == nil {
			continue
		}

		iX := 0.0
		fX := float64(tile.Image.Width)
		iY := 0.0
		fY := float64(tile.Image.Height)

		sprite := w.sprites[tile.Image.Source]
		if sprite == nil {
			continue
		}
		batchIdx, ok := w.batchIndices[tile.Image.Source]
		if !ok {
			continue
		}

		sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
		flip := 1.0
		if dTile.HorizontalFlip {
			flip = -1.0
		}
		sprite.Draw(w.batches[batchIdx], pixel.IM.ScaledXY(pixel.ZV, pixel.V(flip, 1)).Moved(obj.Center))
	}

	for _, batch := range w.batches {
		batch.Draw(w.cnv2)
	}

	if w.IsDebug {
		imd := imdraw.New(nil)
		imd.Color = color.RGBA{255, 0, 0, 1}
		for _, p := range w.visiblePhys {
			vertices := p.Rect().Vertices()
			for _, v := range vertices {
				imd.Push(v)
			}
			imd.Rectangle(1)
		}
		imd.Draw(w.cnv2)
	}

	drawSpells(w.cnv2)

	if w.crtrs != nil {
		w.crtrs.Draw(w.cnv2)
	}
	loot.Draw(w.cnv2)
	particles.Draw(w.cnv2)
	projectiles.Draw(w.cnv2)
	objects.Draw(w.cnv2)

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
