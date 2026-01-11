package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whu-ham/ham-gateway/internal/ham"
)

// HamHandler handles requests to HAM API
type HamHandler struct {
	hamClient ham.ClientInterface
}

// NewHamHandler creates a new HAM handler
func NewHamHandler(client ham.ClientInterface) *HamHandler {
	return &HamHandler{
		hamClient: client,
	}
}

// SearchCourse handles GET /api/v1/external/ham/course/search
// Query params: keyword (required), keyword_type (optional, default 0)
func (h *HamHandler) SearchCourse(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "400",
			"message": "keyword is required",
		})
		return
	}

	keywordTypeStr := c.DefaultQuery("keyword_type", "0")
	keywordType, err := strconv.ParseInt(keywordTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "400",
			"message": "invalid keyword_type",
		})
		return
	}

	resp, err := h.hamClient.SearchCourse(c.Request.Context(), keyword, int32(keywordType))
	if err != nil {
		log.Printf("[SearchCourse] HAM API call failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "500",
			"message": "Failed to connect to HAM API",
		})
		return
	}

	// Transform response to match gateway format
	courseItems := make([]gin.H, 0, len(resp.Item))
	for _, item := range resp.Item {
		courseItems = append(courseItems, gin.H{
			"name":       item.Value,
			"instructor": "",
			"type":       strconv.Itoa(int(item.Type)),
			"score":      item.Hit,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "00000",
		"message": "Success",
		"data":    courseItems,
	})
}

// GetCourseStat handles GET /api/v1/external/ham/score/stat
// Query params: course_name (required), instructor (required)
func (h *HamHandler) GetCourseStat(c *gin.Context) {
	courseName := c.Query("course_name")
	instructor := c.Query("instructor")

	if courseName == "" || instructor == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "400",
			"message": "course_name and instructor are required",
		})
		return
	}

	resp, err := h.hamClient.GetCourseScoreItem(c.Request.Context(), courseName, instructor)
	if err != nil {
		log.Printf("[GetCourseStat] HAM API call failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "500",
			"message": "Failed to connect to HAM API",
		})
		return
	}

	// Transform response to match gateway format
	scoreRanges := make([]gin.H, 0, len(resp.Item.Range))
	for _, r := range resp.Item.Range {
		scoreRanges = append(scoreRanges, gin.H{
			"from":  r.From,
			"to":    r.To,
			"total": r.Total,
			"color": r.Color,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "00000",
		"message": "Success",
		"data": gin.H{
			"name":       resp.Item.Name,
			"instructor": resp.Item.Instructor,
			"average":    resp.Item.Average,
			"total":      resp.Item.Total,
			"range":      scoreRanges,
		},
	})
}
