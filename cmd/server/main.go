package main

import (
	"MatrixGo/api/server"
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"MatrixGo/internal/vector"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if _, err := os.Stat("web/frontend"); os.IsNotExist(err) {
		log.Fatal("Frontend directory does not exist")
	}

	r := gin.Default()

	// Настройка CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	// API endpoints
	api := r.Group("/api")
	{
		api.POST("/solve", handleSolveSystem)
		api.POST("/determinant", handleDeterminant)
		api.POST("/rank", handleRank)
		api.POST("/inverse", handleInverse)
	}

	// Обработка статических файлов
	r.GET("/", func(c *gin.Context) {
		c.File("web/frontend/index.html")
	})
	r.GET("/app.js", func(c *gin.Context) {
		c.File("web/frontend/app.js")
	})
	r.GET("/styles.css", func(c *gin.Context) {
		c.File("web/frontend/styles.css")
	})

	workDir, _ := filepath.Abs(".")
	log.Printf("Working directory: %s", workDir)
	log.Printf("Starting server at http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}

func handleSolveSystem(c *gin.Context) {
	// Логируем тело запроса
	body, err := c.GetRawData()
	if err != nil {
		fmt.Printf("Ошибка чтения тела запроса: %v\n", err)
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "Ошибка чтения запроса"})
		return
	}
	fmt.Printf("Получено тело запроса: %s\n", string(body))

	// Восстанавливаем тело запроса для последующего использования
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req server.SystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Ошибка разбора JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "Неверный формат данных"})
		return
	}

	fmt.Printf("Успешно разобран запрос: %+v\n", req)
	fmt.Printf("Тип матрицы: %s\n", req.Matrix.Type)
	fmt.Printf("Размеры матрицы: %dx%d\n", req.Matrix.Rows, req.Matrix.Cols)
	fmt.Printf("Данные матрицы: %v\n", req.Matrix.Data)
	fmt.Printf("Вектор: %v\n", req.Vector)

	// Парсим матрицу
	mat, err := server.ParseMatrix(req.Matrix)
	if err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: err.Error()})
		return
	}

	// Парсим вектор
	var vec interface{}
	switch req.Matrix.Type {
	case "float64":
		vec, err = server.ParseFloat64Vector(req.Vector)
	case "complex":
		vec, err = server.ParseComplexVector(req.Vector)
	case "rational":
		vec, err = server.ParseRationalVector(req.Vector)
	case "gf":
		vec, err = server.ParseGFVector(req.Vector, req.Matrix.ModP)
	default:
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "неподдерживаемый тип матрицы"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: fmt.Sprintf("ошибка парсинга вектора: %v", err)})
		return
	}

	// Решаем систему
	var result interface{}
	var solveErr error
	switch m := mat.(type) {
	case *matrix.Matrix[field.Float64]:
		result, solveErr = matrix.SolveSystemParallel(m, vec.(*vector.Vector[field.Float64]))
	case *matrix.Matrix[field.Complex]:
		result, solveErr = matrix.SolveSystemParallel(m, vec.(*vector.Vector[field.Complex]))
	case *matrix.Matrix[field.Rational]:
		result, solveErr = matrix.SolveSystemParallel(m, vec.(*vector.Vector[field.Rational]))
	case *matrix.Matrix[field.GF]:
		result, solveErr = matrix.SolveSystemParallel(m, vec.(*vector.Vector[field.GF]))
	default:
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "неподдерживаемый тип матрицы"})
		return
	}

	if solveErr != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: solveErr.Error()})
		return
	}

	// Конвертируем результат в строки
	var resultStrings []string
	switch res := result.(type) {
	case *vector.Vector[field.Float64]:
		resultStrings = server.VectorToStrings(res)
	case *vector.Vector[field.Complex]:
		resultStrings = server.VectorToStrings(res)
	case *vector.Vector[field.Rational]:
		resultStrings = server.VectorToStrings(res)
	case *vector.Vector[field.GF]:
		resultStrings = server.VectorToStrings(res)
	}

	c.JSON(http.StatusOK, server.MatrixResponse{Result: [][]string{resultStrings}})
}

func handleDeterminant(c *gin.Context) {
	var req server.MatrixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "Неверный формат данных"})
		return
	}

	mat, err := server.ParseMatrix(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: err.Error()})
		return
	}

	var result string
	switch m := mat.(type) {
	case *matrix.Matrix[field.Float64]:
		result = fmt.Sprintf("%v", m.DeterminantParallel())
	case *matrix.Matrix[field.Complex]:
		result = fmt.Sprintf("%v", m.DeterminantParallel())
	case *matrix.Matrix[field.Rational]:
		result = fmt.Sprintf("%v", m.DeterminantParallel())
	case *matrix.Matrix[field.GF]:
		result = fmt.Sprintf("%v", m.DeterminantParallel())
	default:
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "Неподдерживаемый тип матрицы"})
		return
	}

	c.JSON(http.StatusOK, server.MatrixResponse{Value: result})
}

func handleRank(c *gin.Context) {
	var req server.MatrixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "Неверный формат данных"})
		return
	}

	mat, err := server.ParseMatrix(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: err.Error()})
		return
	}

	var rank int
	switch m := mat.(type) {
	case *matrix.Matrix[field.Float64]:
		rank = m.RankParallel()
	case *matrix.Matrix[field.Complex]:
		rank = m.RankParallel()
	case *matrix.Matrix[field.Rational]:
		rank = m.RankParallel()
	case *matrix.Matrix[field.GF]:
		rank = m.RankParallel()
	default:
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "Неподдерживаемый тип матрицы"})
		return
	}

	c.JSON(http.StatusOK, server.MatrixResponse{Value: fmt.Sprintf("%d", rank)})
}

func handleInverse(c *gin.Context) {
	var req server.MatrixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "Неверный формат данных"})
		return
	}

	mat, err := server.ParseMatrix(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: err.Error()})
		return
	}

	var result interface{}
	switch m := mat.(type) {
	case *matrix.Matrix[field.Float64]:
		result, err = m.InverseParallel()
	case *matrix.Matrix[field.Complex]:
		result, err = m.InverseParallel()
	case *matrix.Matrix[field.Rational]:
		result, err = m.InverseParallel()
	case *matrix.Matrix[field.GF]:
		result, err = m.InverseParallel()
	default:
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: "Неподдерживаемый тип матрицы"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, server.MatrixResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, server.MatrixResponse{Result: server.MatrixToStrings(result)})
}
