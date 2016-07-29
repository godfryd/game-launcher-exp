package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type GameListScene struct{}

var (
	scrollSpeed float32 = 700

	worldWidth  int = 1400
	worldHeight int = 800

	fontRegular *common.Font
	fontSelected *common.Font
)

type Game struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	PkgName string
	Name string
	ScreenFile string
	ScreenImg common.Drawable
	Description string
}

type GameScreenshot struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (*GameListScene) Preload() {
	err := engo.Files.Load("SourceSansPro-Regular.ttf")
	if err != nil {
		panic(err)
	}
}

func loadGamesData() []*Game {
        games := []*Game{}

	file, err := os.Open("assets/games.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot open games.csv:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, ",")
                game := &Game{
			BasicEntity: ecs.NewBasic(),
			PkgName: data[0],
                        Name: data[1],
                        ScreenFile: data[2],
                        Description: data[3]}
		games = append(games, game)

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading games list:", err)
	}

	return games
}

type GamesListSystem struct {
	world *ecs.World
	games []*Game
	screenshot *GameScreenshot
	offset float32
	curPos int
	speed float32
	nextStop float32
}

func (gl *GamesListSystem) New(w *ecs.World) {
	gl.world = w
	gl.nextStop = 100000
}

func (gl *GamesListSystem) Add(game *Game) {
	gl.games = append(gl.games, game)
}

func (gl *GamesListSystem) computeOffset(dt float32) {
	if engo.Input.Button("moveup").Down() && gl.curPos > 0 {
		gl.speed = 200
		if gl.nextStop == 100000 {
			gl.nextStop = gl.offset + 20
			gl.curPos = gl.curPos - 1
		}
	} else if engo.Input.Button("movedown").Down() && gl.curPos < len(gl.games) - 1 {
		gl.speed = -200
		if gl.nextStop == 100000 {
			gl.nextStop = gl.offset - 20
			gl.curPos = gl.curPos + 1
		}
	}


	if gl.speed != 0 {
		gl.offset = gl.offset + (gl.speed * dt)
		if (gl.speed > 0 && gl.offset >= gl.nextStop) {
			gl.offset = gl.nextStop
			gl.speed = 0
			gl.nextStop = 100000
		} else if (gl.speed < 0 && gl.offset <= gl.nextStop) {
			gl.offset = gl.nextStop
			gl.speed = 0
			gl.nextStop = 100000
		}
	}
}

func (gl *GamesListSystem) Update(dt float32) {
	if engo.Input.Button("exit").Down() {
		engo.Exit()
        }

	gl.computeOffset(dt)


	for index, game := range gl.games {
		game.RenderComponent.Drawable.Close()
		if index == gl.curPos {
			game.RenderComponent.Drawable = fontSelected.Render(game.Name)
			gl.screenshot.RenderComponent.Drawable = game.ScreenImg
			gl.screenshot.RenderComponent.Hidden = false
		} else {
			game.RenderComponent.Drawable = fontRegular.Render(game.Name)
		}
		game.SpaceComponent.Position.Y = engo.GameHeight() / 20.0 + float32(index) * 30.0 + gl.offset
	}
}

func (gl *GamesListSystem) Remove(e ecs.BasicEntity) {
}

func (*GameListScene) Setup(w *ecs.World) {
	common.SetBackground(color.Black)

	renderSys := &common.RenderSystem{}
	w.AddSystem(renderSys)
	gamesListSys := &GamesListSystem{}
	w.AddSystem(gamesListSys)

	fontRegular = &common.Font{
		URL:  "SourceSansPro-Regular.ttf",
		FG:   color.RGBA{100, 255, 100, 0xff},
		Size: 30,
	}
	err := fontRegular.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	fontSelected = &common.Font{
		URL:  "SourceSansPro-Regular.ttf",
		FG:   color.RGBA{255, 50, 50, 0xff},
		Size: 30,
	}
	err = fontSelected.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	games := loadGamesData()

	for _, game := range games {
		fmt.Fprintln(os.Stderr, "game:", game.Name)
		game.RenderComponent.Drawable = fontRegular.Render(game.Name)
		game.SpaceComponent.Position.X = engo.GameWidth() / 20

		err := engo.Files.Load(game.ScreenFile)
		if err != nil {
			panic(err)
		}
		game.ScreenImg, err = common.PreloadedSpriteSingle(game.ScreenFile)
		if err != nil {
			panic(err)
		}

		renderSys.Add(&game.BasicEntity, &game.RenderComponent, &game.SpaceComponent)
		gamesListSys.Add(game)
	}

	gamesListSys.screenshot = &GameScreenshot{}
	gamesListSys.screenshot.SpaceComponent.Position.X = 5 * engo.GameWidth() / 20
	gamesListSys.screenshot.SpaceComponent.Position.Y = engo.GameHeight() / 20
	gamesListSys.screenshot.RenderComponent.Hidden = true
	renderSys.Add(&gamesListSys.screenshot.BasicEntity, &gamesListSys.screenshot.RenderComponent, &gamesListSys.screenshot.SpaceComponent)

	engo.Input.RegisterButton("moveup", engo.ArrowUp)
	engo.Input.RegisterButton("movedown", engo.ArrowDown)
	engo.Input.RegisterButton("exit", engo.Escape)
}

func (*GameListScene) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:          "Game Launcher",
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
	}
	engo.Run(opts, &GameListScene{})
}
