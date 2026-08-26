package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"impactgraph/internal/repository"
)

type Handler struct {
	store repository.Store
	mode  string
}

func New(store repository.Store, mode string) *Handler {
	return &Handler{store: store, mode: mode}
}

func (h *Handler) Health(c *gin.Context) {
	err := h.store.Health(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "graph database unavailable", "mode": h.mode})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "database": h.mode})
}

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	res, err := h.store.Search(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "graph database unavailable"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetEntity(c *gin.Context) {
	typ := c.Param("type")
	id := c.Param("id")
	res, err := h.store.GetEntity(c.Request.Context(), typ, id)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "entity not found"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "graph database unavailable"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetEntityGraph(c *gin.Context) {
	typ := c.Param("type")
	id := c.Param("id")
	res, err := h.store.GetEntityGraph(c.Request.Context(), typ, id)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "graph database unavailable"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetImpact(c *gin.Context) {
	id := c.Param("id")
	res, err := h.store.GetImpact(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "package not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "package not found"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "graph database unavailable"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetPath(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params required"})
		return
	}
	res, err := h.store.GetPath(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "graph database unavailable"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetStats(c *gin.Context) {
	res, err := h.store.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "graph database unavailable"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllNodes(c *gin.Context) {
	res, err := h.store.GetAllNodes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "graph database unavailable"})
		return
	}
	c.JSON(http.StatusOK, res)
}
