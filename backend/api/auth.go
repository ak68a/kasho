package api

import (
	"context"
	"database/sql"
	db "github/kasho/backend/db/sqlc"
	"github/kasho/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Auth struct {
	server *Server
}

func (a Auth) router(server *Server) {
	a.server = server

	serverGroup := server.router.Group("/auth")
	serverGroup.POST("login", a.login)
	serverGroup.POST("register", a.register)
}

type UserParams struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (a *Auth) register(c *gin.Context) {
	// user := New(UserParams{})
	var user UserParams

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateUserParams{
		Email: user.Email,
		HashedPassword: hashedPassword,
	}

	newUser, err := a.server.queries.CreateUser(context.Background(), arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code {
			case "23505":
				c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{}.toUserResponse(&newUser))
}

func (a Auth) login(c *gin.Context) {
	user := new(UserParams)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := a.server.queries.GetUserByEmail(context.Background(), user.Email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.VerifyPassword(user.Password, dbUser.HashedPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := tokenController.CreateToken(dbUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

