package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func main() {
	// Conexion a la base de datos
	var err error
	db, err = pgx.Connect(context.Background(), "postgres://user:password@db:5432/matches_db")
	if err != nil {
		log.Fatalf("Error conectando a la base de datos", err)
	}
	defer db.Close(context.Background())

	// Configurar Gin con CORS
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Endpoints requeridos
	r.GET("/api/matches", getMatches)
	r.GET("/api/matches/:id", getMatchByID)
	r.POST("/api/matches", createMatch)
	r.PUT("/api/matches/:id", updateMatch)
	r.DELETE("/api/matches/:id", deleteMatch)

	// Servir el frontend
	r.StaticFile("/", "./LaLigaTracker.html")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico") 
	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor iniciado en :%s\n", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
// Middleware CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Next()
	}
}

// Obtener todos los partidos
func getMatches(c *gin.Context) {
    rows, err := db.Query(context.Background(), "SELECT id, team1, team2, score1, score2, date FROM matches")
    if err != nil {
        log.Printf("Error en consulta SQL: %v", err) // Log detallado
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Error obteniendo partidos",
            "details": err.Error(), // Muestra el error real
        })
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
