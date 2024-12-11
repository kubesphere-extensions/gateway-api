package runtime

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewRouterGroup(group, version string, engine *gin.Engine) *gin.RouterGroup {
	return engine.Group(fmt.Sprintf("/kapis/%s/%s", group, version))
}
