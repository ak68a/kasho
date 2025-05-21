package api

import (
	"fmt"
	db "github/kasho/backend/db/sqlc"
	"github/kasho/backend/utils"
	"net/http"

	"database/sql"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	queries *db.Queries
	router *gin.Engine
	config *utils.Config
}

var tokenController *utils.JWTToken

func NewServer(envPath string) *Server {
	config, err := utils.LoadConfig(envPath)
	if err != nil {
		panic(fmt.Sprintf("Could not load config: %v", err))
	}

	conn, err := sql.Open(config.DBdriver, config.DB_source_live)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to database: %v", err))
	}

	tokenController = utils.NewJWTToken(config)

	q := db.New(conn)

	g := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", currencyValidator)
	}

	g.Use(cors.Default())

	return &Server	{
		queries: q,
		router: g,
		config: config,
	}
}

func (s *Server) Start(port int) {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to Kasho!"})
	})

	User{}.router(s)
	Auth{}.router(s)
	Account{}.router(s)

	s.router.Run(fmt.Sprintf(":%v", port))
}

// NewServer (OLD)
// func NewServer(port int) {
// 	g := gin.Default()

// 	g.GET("/", func(ctx *gin.Context) {
// 		ctx.JSON(200, gin.H{"message": "Welcome to Kasho!"})
// 	})

// 	g.Run(fmt.Sprintf(":%v", port))
// }