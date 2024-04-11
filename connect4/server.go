package connect4

import (
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var updateChannel chan [][]string = make(chan [][]string)

func serve() {

	r := gin.Default()
	r.LoadHTMLGlob("connect4/templates/*")

	defer close(updateChannel)
	// Define a route that renders the 'index.tmpl' template
	colors := map[string]string{
		"1": "blue",
		"2": "red",
	}
	r.GET("/", func(c *gin.Context) {
		// Create an 8x8 grid (you can customize this data as needed)
		rows := make([][]string, 8)
		for i := range rows {
			rows[i] = make([]string, 8)
		}
		for i := range Game.field {
			for j := range Game.field {
				rows[i][j] = strconv.Itoa(Game.field[i][j])
			}
		}

		slices.Reverse(rows)
		// Pass the grid data to the template
		c.HTML(http.StatusOK, "game.go.tmpl", gin.H{
			"Rows":   rows,
			"Colors": colors,
		})
	})
	r.GET("/connect4", func(c *gin.Context) {
		// Create an 8x8 grid (you can customize this data as needed)
		rows := make([][]string, 8)
		for i := range rows {
			rows[i] = make([]string, 8)
		}
		for i := range Game.field {
			for j := range Game.field {
				rows[i][j] = strconv.Itoa(Game.field[i][j])
			}
		}

		slices.Reverse(rows)
		// Pass the grid data to the template
		c.HTML(http.StatusOK, "game.go.tmpl", gin.H{
			"Rows":   rows,
			"Colors": colors,
		})
	})
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	connections := map[*websocket.Conn]bool{}

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		connections[conn] = true
		defer delete(connections, conn)
		defer conn.Close()

		for {
			x, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if x == websocket.TextMessage {
				switch string(msg) {
				case "r":
					rotate()
				default:
					msgStrings := strings.Split(string(msg), ",")
					team, _ := strconv.Atoi(msgStrings[0])
					column, _ := strconv.Atoi(msgStrings[1])
					insert(team, column)

				}
			}

		}
		// Handle WebSocket communication here
		// ...
	})
	go func() {
		for newGrid := range updateChannel {
			for c := range connections {
				(*c).WriteJSON(newGrid)
			}
		}
	}()
	r.Run("localhost:8080")

}

func rotate() {
	defer func() {
		rows := make([][]string, 8)
		for i := range rows {
			rows[i] = make([]string, 8)
		}
		for i := range Game.field {
			for j := range Game.field {
				rows[i][j] = strconv.Itoa(Game.field[i][j])
			}
		}
		updateChannel <- rows
	}()
	Game.Rotate()
	_, y := Game.scanForConnect4()
	if len(y) == 0 {
		Game.turn = (Game.turn % 2) + 1
	} else {
		Game.Clear()
	}
}

func insert(team, column int) {

	defer func() {
		rows := make([][]string, 8)
		for i := range rows {
			rows[i] = make([]string, 8)
		}
		for i := range Game.field {
			for j := range Game.field {
				rows[i][j] = strconv.Itoa(Game.field[i][j])
			}
		}
		updateChannel <- rows
	}()

	if Game.turn == team {
		worked := Game.Insert(team, column)
		if worked {
			x, y := Game.scanForConnect4()
			if x == team {
				// team x won the Game
				Game.Clear()
			} else if len(y) == 0 {
				Game.turn = (team % 2) + 1
			} else {
				Game.Clear()
			}
		}
	}
}
