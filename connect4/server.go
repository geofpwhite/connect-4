package connect4

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

func serve() {

	r := gin.Default()
	r.LoadHTMLGlob("connect4/templates/*")
	r.GET("/rotate", rotate)
	r.GET("/insert/:team/:column", insert)

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
		c.HTML(http.StatusOK, "game.tmpl", gin.H{
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
		c.HTML(http.StatusOK, "game.tmpl", gin.H{
			"Rows":   rows,
			"Colors": colors,
		})
	})

	r.Run("localhost:8080")

}

func rotate(c *gin.Context) {
	Game.Rotate()
	Game.Fall()
	_, y := Game.scanForConnect4()
	if len(y) == 0 {
		Game.turn = (Game.turn % 2) + 1
	} else {
		Game.Clear()
	}
}

func insert(c *gin.Context) {

	teamString, columnString := c.Param("team"), c.Param("column")
	team, err := strconv.Atoi(teamString)
	if err != nil {
		panic("bad")
	}
	column, err := strconv.Atoi(columnString)
	if err != nil {
		panic("bad")
	}

	if Game.turn == team {
		worked := Game.Insert(team, column)
		if worked {
			Game.Fall()
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
