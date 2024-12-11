package controller

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type GatewayWebhook struct {
	client  client.Client
	decoder admission.Decoder
	log     logr.Logger
}

func (w *GatewayWebhook) Handle(_ context.Context, req admission.Request) admission.Response {
	// TODO implement me!
	return admission.Response{}
}

func (w *GatewayWebhook) SetupWithWebhook(mgr manager.Manager) error {
	w.client = mgr.GetClient()
	w.log = mgr.GetLogger().WithName("gateway-webhook")
	w.decoder = admission.NewDecoder(mgr.GetScheme())
	mgr.GetWebhookServer().Register("/gatewayapi-kubesphere-io-v1alpha1-gateway", &webhook.Admission{Handler: w})
	return nil
}
