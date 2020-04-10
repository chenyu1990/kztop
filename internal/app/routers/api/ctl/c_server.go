package ctl

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewServer 创建server控制器
func NewServer() *Server {
	return &Server{
	}
}

type Server struct {
}

func (a *Server) Query(c *gin.Context) {
	c.HTML(http.StatusOK, "server/index.tmpl", gin.H{
		"title": "Posts",
	})
}