package response
import (
	"github.com/gin-gonic/gin"
)
type Response struct{Code int `json:"code"`;Message string `json:"message"`;Detail interface{} `json:"detail"`;Data interface{} `json:"data"`}
type PaginatedData struct{Items interface{} `json:"items"`;Total int64 `json:"total"`;Page int `json:"page"`;PageSize int `json:"page_size"`;TotalPages int `json:"total_pages"`}
func OK(c *gin.Context, data interface{}){c.JSON(200,Response{Code:0,Message:"ok",Data:data})}
func Created(c *gin.Context, data interface{}){c.JSON(201,Response{Code:0,Message:"created",Data:data})}
func Fail(c *gin.Context, httpStatus int, code int, message string){c.JSON(httpStatus,Response{Code:code,Message:message})}
func FailWithDetail(c *gin.Context, httpStatus int, code int, message string, detail interface{}){c.JSON(httpStatus,Response{Code:code,Message:message,Detail:detail})}
