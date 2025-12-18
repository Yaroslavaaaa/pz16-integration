package httpapi

import (
	"net/http"
	"strconv"

	"example.com/pz16/internal/models"
	"example.com/pz16/internal/service"

	"github.com/gin-gonic/gin"
)

type Router struct{ Svc *service.Service }

func (rt Router) Register(r *gin.Engine) {
	r.POST("/notes", rt.createNote)
	r.GET("/notes/:id", rt.getNote)
}

func (rt Router) createNote(c *gin.Context) {
	var in struct{ Title, Content string }
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	n := models.Note{Title: in.Title, Content: in.Content}
	if err := rt.Svc.Create(c, &n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, n)
}

func (rt Router) getNote(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	n, err := rt.Svc.Get(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, n)
}
