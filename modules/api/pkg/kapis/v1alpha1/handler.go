package v1alpha1

import (
	"net/http"

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

	actualNamespace = "gatewayapi.kubesphere.io/actual-namespace"
	actualWorkspace = "gatewayapi.kubesphere.io/actual-workspace"
	gatewayApiScope = "gatewayapi.kubesphere.io/scope"
)

type Handler struct {
	client rtclient.Client
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

func (h *Handler) GetGateway(ctx *gin.Context) {
	gwParams := handleRequestParams(ctx, "gateway")
	list := &apisv1.GatewayList{}
	labelMap := map[string]string{}
	labelMap[gatewayApiScope] = gwParams.Scope
	if gwParams.Workspace != "" {
		labelMap[actualWorkspace] = gwParams.Workspace
	}
	if gwParams.Namespace != "" {
		labelMap[actualNamespace] = gwParams.Namespace
	}

	err := h.client.List(ctx.Request.Context(), list, rtclient.MatchingLabels(labelMap), rtclient.InNamespace(""))
	if err != nil {
		api.HandleError(ctx.Writer, ctx.Request, err)
		return
	}

	var gateway *apisv1.Gateway
	for _, item := range list.Items {
		if item.Name == gwParams.ResourceName {
			gateway = &item
			break
		}
	}

	if gateway == nil {
		api.HandleNotFound(ctx.Writer, ctx.Request, errors.NewNotFound(apisv1.Resource(resourceNameGateway), gwParams.ResourceName))
		return
	}

	ctx.JSON(http.StatusOK, gateway)
}

func (h *Handler) ListGateways(ctx *gin.Context) {
	gwParams := handleRequestParams(ctx, "gateway")
	list := &apisv1.GatewayList{}
	labelMap := map[string]string{}
	labelMap[gatewayApiScope] = gwParams.Scope
	if gwParams.Workspace != "" {
		labelMap[actualWorkspace] = gwParams.Workspace
	}
	if gwParams.Namespace != "" {
		labelMap[actualNamespace] = gwParams.Namespace
	}

	err := h.client.List(ctx.Request.Context(), list, rtclient.MatchingLabels(labelMap), rtclient.InNamespace(""))
	if err != nil {
		api.HandleError(ctx.Writer, ctx.Request, err)
		return
	}
}

func (h *Handler) CreateGateway(ctx *gin.Context) {
	// TODO implement me!
}

func (h *Handler) UpdateGateway(ctx *gin.Context) {
	// TODO implement me!
}

func (h *Handler) DeleteGateway(ctx *gin.Context) {
	// TODO implement me!
}

func handleRequestParams(ctx *gin.Context, resourceName string) ResourceParams {
	s := ResourceParams{
		Scope:        scopeCluster,
		ResourceName: ctx.Param(resourceName),
	}
	workspace := ctx.Param(paramWorkspace)
	namespace := ctx.Param(paramNamespace)

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
