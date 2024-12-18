package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kubesphere-extensions/gateway-api/pkg/kapis/v1alpha1"
	"k8s.io/klog/v2"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

type APIServer struct {
	Server *http.Server

	// webservice container, where all webservice defines
	Engine *gin.Engine

	// controller-runtime client
	RuntimeClient rtclient.Client
}

func (s *APIServer) installAPIs() {
	// add health check APIs
	s.Engine.GET("/healthz", func(c *gin.Context) {
		_ = healthz.Ping(c.Request)
	})
	s.Engine.GET("/readyz", func(c *gin.Context) {
		_ = healthz.Ping(c.Request)
	})

	v1alpha1.AddRouterGroup(s.Engine, s.RuntimeClient)
}

func (s *APIServer) PrepareRun() error {
	s.Engine = gin.New()
	s.Engine.Use(gin.Recovery())
	s.Engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("[%s] %s - \"%s %s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	s.installAPIs()

	s.Server.Handler = s.Engine

	return nil
}

func (s *APIServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = s.Server.Shutdown(ctx)
	}()

	s.Server.Handler = s.Engine

	klog.Infof("Start listening on %s", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
