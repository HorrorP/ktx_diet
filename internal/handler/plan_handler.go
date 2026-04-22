package handler

import (
	"net/http"

	"ktx-diet/internal/domain"
	"ktx-diet/internal/service"

	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	planService service.PlanService
}

func NewPlanHandler(planService service.PlanService) *PlanHandler {
	return &PlanHandler{planService: planService}
}

func (h *PlanHandler) GeneratePlan(c *gin.Context) {
	var req domain.GeneratePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.planService.GeneratePlan(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PlanHandler) ReplaceMeal(c *gin.Context) {
	var req domain.ReplaceMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.planService.ReplaceMeal(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plan": plan})
}
