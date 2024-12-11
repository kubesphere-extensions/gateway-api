package v1alpha1

import (
	"github.com/gin-gonic/gin"
	apiruntime "github.com/kubesphere-extensions/gateway-api/pkg/apiserver/runtime"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func AddRouterGroup(engin *gin.Engine, client rtclient.Client) {
	group := apiruntime.NewRouterGroup("gatewayapi.kubesphere.io", "v1alpha1", engin)
	handler := NewHandler(client)

	group.GET("/gateways/:gateway", handler.GetGateway)
	group.GET("/gateways", handler.ListGateways)
	group.POST("/gateways", handler.CreateGateway)
	group.PUT("/gateways", handler.UpdateGateway)
	group.DELETE("/gateways/:gateway", handler.DeleteGateway)

	group.GET("/workspaces/:workspace/gateways/:gateway", handler.GetGateway)
	group.GET("/workspaces/:workspace/gateways", handler.ListGateways)
	group.POST("/workspaces/:workspace/gateways", handler.CreateGateway)
	group.PUT("/workspaces/:workspace/gateways", handler.UpdateGateway)
	group.DELETE("/workspaces/:workspace/gateways/:gateway", handler.DeleteGateway)

	group.GET("/namespaces/:namespace/gateways/:gateway", handler.GetGateway)
	group.GET("/namespaces/:namespace/gateways", handler.ListGateways)
	group.POST("/namespaces/:namespace/gateways", handler.CreateGateway)
	group.PUT("/namespaces/:namespace/gateways", handler.UpdateGateway)
	group.DELETE("/namespaces/:namespace/gateways/:gateway", handler.DeleteGateway)

}
