package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
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
	r.PUT("/api/matches/:id", upmatchDateMatch)
	r.DELETE("/api/matches/:id", deleteMatch)
	r.PATCH("/api/matches/:id/goals", updateGoals)
	r.PATCH("/api/matches/:id/yellowcards", updateYellowCard)
	r.PATCH("/api/matches/:id/redcards", updateRedCard)
	r.PATCH("/api/matches/:id/extratime", updateExtraTime)
	
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
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
// Obtener todos los partidos
func getMatches(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT id, hometeam, awayteam, score1, score2, matchdate FROM matches")
	if err != nil {
		log.Printf("Error en consulta SQL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error obteniendo partidos",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	var matches []map[string]interface{}

	for rows.Next() {
		var id, score1, score2 int
		var homeTeam, awayTeam string
		var matchDate time.Time

		err := rows.Scan(&id, &homeTeam, &awayTeam, &score1, &score2, &matchDate)
		if err != nil {
			log.Printf("Error escaneando fila: %v", err)
			continue 
		}

		matches = append(matches, map[string]interface{}{
			"id":        id,
			"homeTeam":  strings.ToLower(homeTeam),
			"awayTeam":  strings.ToLower(awayTeam),
			"score1":    score1,
			"score2":    score2,
			"matchDate": matchDate.Format("2006-01-02 15:04:05"), 
		})
	}

	c.JSON(http.StatusOK, matches)
}

// Obtener un partido por ID
func getMatchByID(c *gin.Context) {
	id := c.Param("id")
	var matchID, score1, score2 int
	var homeTeam, awayTeam string
	var matchDate time.Time

	err := db.QueryRow(context.Background(), "SELECT id, homeTeam, awayTeam, score1, score2, matchDate FROM matches WHERE id=$1", id).
		Scan(&matchID, &homeTeam, &awayTeam, &score1, &score2, &matchDate)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
		} else {
			log.Printf("Error obteniendo partido: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener partido"})
		}
		return
	}

	match := map[string]interface{}{
		"id":        matchID,
		"homeTeam":  strings.ToLower(homeTeam),
		"awayTeam":  strings.ToLower(awayTeam),
		"score1":    score1,
		"score2":    score2,
		"matchDate": matchDate.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, match)
}


// Crear un nuevo partido
func createMatch(c *gin.Context) {
	var match struct {
		homeTeam  string `json:"homeTeam"`
		awayTeam  string `json:"awayTeam"`
		Score1 int    `json:"score1"`
		Score2 int    `json:"score2"`
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	_, err := db.Exec(context.Background(), "INSERT INTO matches (homeTeam, awayTeam, score1, score2) VALUES ($1, $2, $3, $4)",
		match.homeTeam, match.awayTeam, match.Score1, match.Score2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear partido"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Partido creado"})
}

// Actualizar un partido existente
func upmatchDateMatch(c *gin.Context) {
	id := c.Param("id")

	var match struct {
		Score1 int `json:"score1"`
		Score2 int `json:"score2"`
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	_, err := db.Exec(context.Background(), "UPmatchDate matches SET score1=$1, score2=$2 WHERE id=$3", match.Score1, match.Score2, id)
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

// Actualizar goles de un partido
func updateGoals(c *gin.Context) {
	id := c.Param("id")

	var scores struct {
		Score1 int `json:"score1"`
		Score2 int `json:"score2"`
	}

	if err := c.ShouldBindJSON(&scores); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	_, err := db.Exec(context.Background(), "UPDATE matches SET score1=$1, score2=$2 WHERE id=$3", scores.Score1, scores.Score2, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar goles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goles actualizados"})
}

// Registrar tarjeta amarilla
func updateYellowCard(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec(context.Background(), "UPDATE matches SET yellow_cards = yellow_cards + 1 WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar tarjeta amarilla"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tarjeta amarilla registrada"})
}

// Registrar tarjeta roja
func updateRedCard(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec(context.Background(), "UPDATE matches SET red_cards = red_cards + 1 WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar tarjeta roja"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tarjeta roja registrada"})
}

// Registrar tiempo extra
func updateExtraTime(c *gin.Context) {
	id := c.Param("id")

	var extraTime struct {
		ExtraMinutes int `json:"extra_minutes"`
	}

	if err := c.ShouldBindJSON(&extraTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	_, err := db.Exec(context.Background(), "UPDATE matches SET extra_time = extra_time + $1 WHERE id=$2", extraTime.ExtraMinutes, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar tiempo extra"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tiempo extra registrado"})
}
