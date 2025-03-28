package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func main() {
	// Conectar a la base de datos
	var err error
	db, err = pgx.Connect(context.Background(), "postgres://user:password@localhost:5432/matches_db")
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer db.Close(context.Background())

	// Crear el router de Gin
	r := gin.Default()

	// Definir rutas
	r.GET("/api/matches", getMatches)
	r.GET("/api/matches/:id", getMatchByID)
	r.POST("/api/matches", createMatch)
	r.PUT("/api/matches/:id", updateMatch)
	r.DELETE("/api/matches/:id", deleteMatch)

	// Ejecutar el servidor en el puerto 8080
	r.Run(":8080")
}

// Obtener todos los partidos
func getMatches(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT * FROM matches")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo partidos"})
		return
	}
	defer rows.Close()

	var matches []map[string]interface{}
	for rows.Next() {
		var id int
		var team1, team2 string
		var score1, score2 int
		var date string
		rows.Scan(&id, &team1, &team2, &score1, &score2, &date)

		matches = append(matches, map[string]interface{}{
			"id":     id,
			"team1":  team1,
			"team2":  team2,
			"score1": score1,
			"score2": score2,
			"date":   date,
		})
	}

	c.JSON(http.StatusOK, matches)
}

// Obtener un partido por ID
func getMatchByID(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow(context.Background(), "SELECT * FROM matches WHERE id=$1", id)

	var match map[string]interface{}
	var matchID int
	var team1, team2 string
	var score1, score2 int
	var date string

	err := row.Scan(&matchID, &team1, &team2, &score1, &score2, &date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
		return
	}

	match = map[string]interface{}{
		"id":     matchID,
		"team1":  team1,
		"team2":  team2,
		"score1": score1,
		"score2": score2,
		"date":   date,
	}

	c.JSON(http.StatusOK, match)
}

// Crear un nuevo partido
func createMatch(c *gin.Context) {
	var match struct {
		Team1  string `json:"team1"`
		Team2  string `json:"team2"`
		Score1 int    `json:"score1"`
		Score2 int    `json:"score2"`
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	_, err := db.Exec(context.Background(), "INSERT INTO matches (team1, team2, score1, score2) VALUES ($1, $2, $3, $4)",
		match.Team1, match.Team2, match.Score1, match.Score2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear partido"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Partido creado"})
}

// Actualizar un partido existente
func updateMatch(c *gin.Context) {
	id := c.Param("id")

	var match struct {
		Score1 int `json:"score1"`
		Score2 int `json:"score2"`
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	_, err := db.Exec(context.Background(), "UPDATE matches SET score1=$1, score2=$2 WHERE id=$3", match.Score1, match.Score2, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar partido"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partido actualizado"})
}

// Eliminar un partido
func deleteMatch(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec(context.Background(), "DELETE FROM matches WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar partido"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partido eliminado"})
}
