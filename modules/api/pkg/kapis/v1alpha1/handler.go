package v1alpha1

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kubesphere-extensions/gateway-api/pkg/api"
	"k8s.io/apimachinery/pkg/api/errors"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
	apisv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	scopeNamespace = "namespace"
	scopeWorkspace = "workspace"
	scopeCluster   = "cluster"
	paramWorkspace = scopeWorkspace
	paramNamespace = scopeNamespace

	resourceNameGateway = "gateway"

	kubesphereControlsSystem = "kubesphere-controls-system"
	defaultWorkingNamespace  = kubesphereControlsSystem

	workingNamespace        = "gatewayapi.kubesphere.io/working-namespace"
	workingWorkspace        = "gatewayapi.kubesphere.io/working-workspace"
	gatewayApiScope         = "gatewayapi.kubesphere.io/scope"
	gatewayListener         = "gatewayapi.kubesphere.io/listener"
	gatewayListenerProtocol = "gatewayapi.kubesphere.io/listener.%s.protocols"
	gatewayListenerPort     = "gatewayapi.kubesphere.io/listener.%s.port"
)

type Handler struct {
	client rtclient.Client
}

type Listener struct {
	Name      string   `json:"name"`
	Protocols []string `json:"protocols,omitempty"`
	Port      int32    `json:"port,omitempty"`
}

type ResourceParams struct {
	Scope        string
	Workspace    string
	Namespace    string
	ResourceName string
}

func NewHandler(client rtclient.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) getGateway(ctx context.Context, params ResourceParams) (*apisv1.Gateway, error) {
	list := &apisv1.GatewayList{}
	labelMap := map[string]string{}
	labelMap[gatewayApiScope] = params.Scope
	if params.Workspace != "" {
		labelMap[workingWorkspace] = params.Workspace
	}
	if params.Namespace != "" {
		labelMap[workingNamespace] = params.Namespace
	}

	err := h.client.List(ctx, list, rtclient.MatchingLabels(labelMap), rtclient.InNamespace(""))
	if err != nil {
		return nil, err
	}
	var gateway *apisv1.Gateway
	for _, item := range list.Items {
		if item.Name == params.ResourceName {
			gateway = &item
			break
		}
	}

	if gateway == nil {
		return nil, errors.NewNotFound(apisv1.Resource(resourceNameGateway), params.ResourceName)
	}

	return gateway, nil
}

func (h *Handler) GetGateway(c *gin.Context) {
	gwParams := handleRequestParams(c, resourceNameGateway)
	gateway, err := h.getGateway(c.Request.Context(), gwParams)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gateway)
}

func (h *Handler) ListGateways(c *gin.Context) {
	gwParams := handleRequestParams(c, resourceNameGateway)
	list := &apisv1.GatewayList{}
	labelMap := map[string]string{}
	labelMap[gatewayApiScope] = gwParams.Scope
	if gwParams.Workspace != "" {
		labelMap[workingWorkspace] = gwParams.Workspace
	}
	if gwParams.Namespace != "" {
		labelMap[workingNamespace] = gwParams.Namespace
	}

	err := h.client.List(c.Request.Context(), list, rtclient.MatchingLabels(labelMap), rtclient.InNamespace(""))
	if err != nil {
		api.HandleError(c, err)
		return
	}
}

func (h *Handler) CreateGateway(c *gin.Context) {
	params := handleRequestParams(c, resourceNameGateway)
	gateway := &apisv1.Gateway{}
	err := c.ShouldBind(gateway)
	if err != nil {
		api.HandleBadRequest(c, err)
		return
	}
	if gateway.Labels == nil {
		gateway.Labels = map[string]string{}
	}
	if _, ok := gateway.Labels[workingNamespace]; !ok && params.Namespace != "" {
		gateway.Labels[workingNamespace] = params.Namespace
	}
	if _, ok := gateway.Labels[workingWorkspace]; !ok && params.Workspace != "" {
		gateway.Labels[workingWorkspace] = params.Workspace
	}
	if _, ok := gateway.Labels[gatewayApiScope]; !ok && params.Scope != "" {
		gateway.Labels[gatewayApiScope] = params.Scope
	}
	if gateway.Namespace == "" {
		gateway.Namespace = defaultWorkingNamespace
	}

	for _, listener := range gateway.Spec.Listeners {
		if listener.AllowedRoutes == nil {
			routes, err := h.newAllowedRoutesByGateway(c.Request.Context(), gateway)
			if err != nil {
				api.HandleError(c, err)
				return
			}
			listener.AllowedRoutes = routes
		}
	}

	err = h.client.Create(c.Request.Context(), gateway)
	if err != nil {
		api.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gateway)
}

func (h *Handler) newAllowedRoutesByGateway(ctx context.Context, gateway *apisv1.Gateway) (*apisv1.AllowedRoutes, error) {
	// TODO implement me!
	return nil, nil
}

func (h *Handler) UpdateGateway(c *gin.Context) {
	// TODO implement me!
}

func (h *Handler) DeleteGateway(c *gin.Context) {
	gwParams := handleRequestParams(c, resourceNameGateway)
	gateway, err := h.getGateway(c.Request.Context(), gwParams)
	if err != nil {
		api.HandleError(c, err)
		return
	}
	err = h.client.Delete(c.Request.Context(), gateway)
	if err != nil {
		api.HandleError(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) GetListeners(c *gin.Context) {
	gatewayClassName := c.Param("gatewayclass")
	gatewayClass := &apisv1.GatewayClass{}
	err := h.client.Get(c.Request.Context(), types.NamespacedName{Name: gatewayClassName}, gatewayClass)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	listener, err := parseListeners(gatewayClass)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"listeners": listener})
}

func (h *Handler) GetGatewayClass(c *gin.Context) {

}

func (h *Handler) ListGatewayClass(c *gin.Context) {

}

/*
parseListeners get the listeners of GatewayClass.
The gateway need to specific the listener information by annotations.

Example:

	apiVersion: gateway.networking.k8s.io/v1
	kind: GatewayClass
	metadata:
	 annotations:
	   gatewayapi.kubesphere.io/listener: web,websecure
	   gatewayapi.kubesphere.io/listener.web.protocols: tcp,http
	   gatewayapi.kubesphere.io/listener.web.port: '8000'
	   gatewayapi.kubesphere.io/listener.websecure.protocols: tls,https
	   gatewayapi.kubesphere.io/listener.websecure.port: '8443'
	 name: traefik
	spec:
	 controllerName: traefik.io/gateway-controller
*/
func parseListeners(gatewayClass *apisv1.GatewayClass) ([]Listener, error) {
	if gatewayClass.Annotations == nil {
		return nil, fmt.Errorf("no listener can be used")
	}
	anno := gatewayClass.Annotations
	listeners := make([]Listener, 0)
	split := strings.Split(anno[gatewayListener], ",")
	if len(split) == 0 {
		return nil, fmt.Errorf("no listener can be used")
	}

	for _, l := range split {
		listener := Listener{Name: l}
		listener.Protocols = strings.Split(anno[fmt.Sprintf(gatewayListenerProtocol, l)], ",")
		pStr := anno[fmt.Sprintf(gatewayListenerPort, l)]
		if pStr != "" {
			port, err := strconv.ParseInt(pStr, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", err)
			}
			listener.Port = int32(port)
		}
		listeners = append(listeners, listener)
	}
	return listeners, nil
}

func handleRequestParams(c *gin.Context, resourceName string) ResourceParams {
	s := ResourceParams{
		Scope:        scopeCluster,
		ResourceName: c.Param(resourceName),
	}
	workspace := c.Param(paramWorkspace)
	namespace := c.Param(paramNamespace)

	if workspace != "" {
		s.Scope = scopeWorkspace
		s.Workspace = workspace
	}
	if namespace != "" {
		s.Scope = scopeNamespace
		s.Namespace = namespace
	}
	return s
}
