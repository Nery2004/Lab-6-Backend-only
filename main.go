package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func main() {
	// Conectar a la base de datos con reintentos
	var err error
	for i := 0; i < 5; i++ {
		db, err = pgx.Connect(context.Background(), "postgres://nery:161204@db:5432/mi_base_de_datos?sslmode=disable")
		if err == nil {
			break
		}
		log.Printf("Intento %d: Error conectando a DB: %v\n", i+1, err)
		time.Sleep(3 * time.Second)
	}
	
	if err != nil {
		log.Fatal("No se pudo conectar a PostgreSQL después de 5 intentos:", err)
	}
	defer db.Close(context.Background())

	// Verificar/Crear tabla al iniciar
	_, err = db.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS matches (
		id SERIAL PRIMARY KEY,
		team1 TEXT NOT NULL,
		team2 TEXT NOT NULL,
		score1 INTEGER,
		score2 INTEGER,
		date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Error creando tabla:", err)
	}

	// Configurar router Gin
	r := gin.Default()

	// Middleware para loggear peticiones
	r.Use(func(c *gin.Context) {
		log.Printf("Petición: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Ruta raíz
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API de Partidos",
			"endpoints": []string{
				"GET    /api/matches",
				"GET    /api/matches/:id",
				"POST   /api/matches",
				"PUT    /api/matches/:id",
				"DELETE /api/matches/:id",
			},
		})
	})

	// Rutas API
	r.GET("/api/matches", getMatches)
	r.GET("/api/matches/:id", getMatchByID)
	r.POST("/api/matches", createMatch)
	r.PUT("/api/matches/:id", updateMatch)
	r.DELETE("/api/matches/:id", deleteMatch)

	// Manejar favicon.ico
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// Iniciar servidor
	log.Println("Servidor iniciado en :8080")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
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
