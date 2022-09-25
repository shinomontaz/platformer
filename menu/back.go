package menu

import (
	"image/color"
	"math"
	"platformer/animation"
	"platformer/common"

	"github.com/salviati/go-tmx/tmx"
	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Scenery struct {
	anim *animation.Anims
	name string
}

var (
	uTime float32

	uLightX float32
	uLightY float32
)

type Back struct {
	rgba          color.Color
	animSpriteNum int
	currtime      float64
	rect          pixel.Rect
	intercanvas   *pixelgl.Canvas
	canvas        *pixelgl.Canvas
	bg            *pixel.Sprite

	//-------from world------------
	tm           *tmx.Map
	Height       float64
	Width        float64
	qtTile       *common.Quadtree
	qtObjs       *common.Quadtree
	scenery      tmx.ObjectGroup
	batches      []*pixel.Batch
	batchIndices map[string]int
	sprites      map[string]*pixel.Sprite

	objects map[int]tmx.Object
	//	objectTiles map[int]*tmx.DecodedTile
	objectAnims map[int]Scenery

	tiles map[int]*tmx.DecodedTile

	visibleTiles []common.Objecter
	visibleObjs  []common.Objecter
}

func NewBack(rect pixel.Rect) *Back {
	tm, err := tmx.ReadFile("menu.tmx")
	if err != nil {
		panic(err)
	}

	b := &Back{
		// intercanvas: pixelgl.NewCanvas(rect),
		// canvas:      pixelgl.NewCanvas(rect),
		intercanvas: pixelgl.NewCanvas(pixel.R(0, 0, rect.W(), rect.H())),
		canvas:      pixelgl.NewCanvas(pixel.R(0, 0, rect.W(), rect.H())),

		rect: rect,
		rgba: colornames.Black,
		tm:   tm,

		batches:      make([]*pixel.Batch, 0),
		batchIndices: make(map[string]int),
		sprites:      make(map[string]*pixel.Sprite),
		tiles:        make(map[int]*tmx.DecodedTile),
		objects:      make(map[int]tmx.Object),
		objectAnims:  make(map[int]Scenery),
	}

	b.init()

	return b
}

func (b *Back) init() {
	for _, og := range b.tm.ObjectGroups {
		if og.Name == "scenery" {
			b.scenery = og
		}
	}

	b.Height = float64(b.tm.TileHeight * b.tm.Height)
	b.Width = float64(b.tm.TileWidth * b.tm.Width)

	//	r := pixel.R(0.0+b.rect.Min.X, 0.0+b.rect.Min.Y, b.Width+b.rect.Min.X, b.Height+b.rect.Min.Y)
	r := pixel.R(0.0, 0.0, b.Width, b.Height)

	b.qtTile = common.New(1, r)
	b.qtObjs = common.New(1, r)

	fs, err := common.LoadFileToString("assets/shader/spotlight.frag.glsl")
	if err != nil {
		panic(err)
	}

	b.canvas.SetUniform("uTime", &uTime)
	b.canvas.SetUniform("uLightX", &uLightX)
	b.canvas.SetUniform("uLightY", &uLightY)

	b.canvas.SetFragmentShader(fs)

	b.initSets()
	b.initTiles()
	b.initObjs()

	bg, err := common.LoadPicture("assets/gamebackground.png")
	if err != nil {
		panic(err)
	}

	b.bg = pixel.NewSprite(bg, pixel.R(0, 0, bg.Bounds().W(), bg.Bounds().H()))
}

func (b *Back) initSets() {
	batchCounter := 0
	for _, tileset := range b.tm.Tilesets {
		if len(tileset.Tiles) > 0 && tileset.Image.Source == "" {
			for _, tile := range tileset.Tiles {
				if _, alreadyLoaded := b.sprites[tile.Image.Source]; !alreadyLoaded {
					sprite, pictureData := loadSprite(tile.Image.Source)
					b.sprites[tile.Image.Source] = sprite
					b.batches = append(b.batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
					b.batchIndices[tile.Image.Source] = batchCounter
					batchCounter++
				}
			}
		} else {
			if _, alreadyLoaded := b.sprites[tileset.Image.Source]; !alreadyLoaded {
				sprite, pictureData := loadSprite(tileset.Image.Source)
				b.sprites[tileset.Image.Source] = sprite
				b.batches = append(b.batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
				b.batchIndices[tileset.Image.Source] = batchCounter
				batchCounter++
			}
		}
	}
}

func (b *Back) initTiles() {
	for _, layer := range b.tm.Layers {
		for tileIndex, tile := range layer.DecodedTiles {
			if tile.Nil {
				continue
			}
			ts := tile.Tileset
			// Calculate the framing for the tile within its tileset's source image
			gamePos := indexToGamePos(tileIndex, b.tm.Width, b.tm.Height)
			pos := gamePos.ScaledXY(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight)))
			b.tiles[tileIndex] = tile

			// pos.X += b.rect.Min.X
			// pos.Y += b.rect.Min.Y

			res := b.qtTile.Insert(common.Objecter{ID: tileIndex, R: pixel.R(pos.X, pos.Y, pos.X+float64(ts.TileWidth), pos.Y+float64(ts.TileHeight))})

			if !res {
				panic("canot insert tile!")
			}
		}
	}
}

func (b *Back) initObjs() {
	for _, o := range b.scenery.Objects {
		min := pixel.V(
			float64(o.X),
			float64(b.Height)-float64(o.Y)-float64(o.Height)-b.rect.Min.Y,
		)
		max := pixel.Vec{
			X: min.X + float64(o.Width),
			Y: min.Y + float64(o.Height),
		}

		rc := pixel.Rect{
			Min: min,
			Max: max,
		}

		b.objectAnims[o.GID] = Scenery{anim: animation.Get("scenery"), name: o.Name}

		b.objects[o.GID] = o
		b.qtObjs.Insert(common.Objecter{ID: o.GID, R: rc})

		screenPos := pixel.Vec{float64(o.X) + o.Width/2, float64(b.Height) - float64(o.Height) - float64(o.Y) - b.rect.Min.Y}
		uLightX = float32(screenPos.X)
		uLightY = float32(screenPos.Y)
	}
}

func (b *Back) Update(dt float64) {
	b.currtime += dt
	b.animSpriteNum = int(math.Floor(b.currtime / 0.2))

	b.visibleTiles = b.qtTile.Retrieve(b.rect)
	b.visibleObjs = b.qtObjs.Retrieve(b.rect)

	uTime += float32(dt)
}

func (b *Back) Draw(t pixel.Target) {
	b.intercanvas.Clear(pixel.RGB(0, 0, 0))
	b.canvas.Clear(pixel.RGB(0, 0, 0))

	b.bg.Draw(b.intercanvas, pixel.IM.Moved(b.intercanvas.Bounds().Center()))

	for _, batch := range b.batches {
		batch.Clear()
	}

	for _, t := range b.visibleTiles {
		tile := b.tiles[t.ID]
		ts := tile.Tileset
		tID := int(tile.ID)

		numRows := ts.Tilecount / ts.Columns
		x, y := tileIDToCoord(tID, ts.Columns, numRows)

		iX := float64(x) * float64(ts.TileWidth)
		fX := iX + float64(ts.TileWidth)
		iY := float64(y) * float64(ts.TileHeight)
		fY := iY + float64(ts.TileHeight)

		sprite := b.sprites[ts.Image.Source]
		sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
		pos := t.R.Center()
		sprite.Draw(b.batches[b.batchIndices[ts.Image.Source]], pixel.IM.Moved(pos))
	}
	for _, batch := range b.batches {
		batch.Draw(b.intercanvas)
	}
	for _, obj := range b.visibleObjs {
		o := b.objects[obj.ID]
		scenery := b.objectAnims[o.GID]

		pic, rect := scenery.anim.GetSprite(scenery.name, b.animSpriteNum)
		sprite := pixel.NewSprite(nil, pixel.Rect{})
		sprite.Set(pic, rect)
		sprite.Draw(b.intercanvas, pixel.IM.Moved(obj.R.Center()))
	}

	b.intercanvas.Draw(b.canvas, pixel.IM.Moved(b.canvas.Bounds().Center()))
	b.canvas.Draw(t, pixel.IM.Moved(b.rect.Center()))
}
