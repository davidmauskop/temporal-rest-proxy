package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/render-examples/temporal-rest-proxy/server/generated/api/workflowservice/v1"
	"github.com/render-examples/temporal-rest-proxy/server/rpc"
)

func getBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header provided")
	}

	authHeaderParts := strings.SplitN(authHeader, " ", 2)
	if len(authHeaderParts) != 2 {
		return "", errors.New("malformed auth header")
	}

	if authHeaderParts[0] != "Bearer" {
		return "", errors.New("wrong auth type, expected \"Bearer\"")
	}

	return authHeaderParts[1], nil
}

func checkAuthHeader(r *http.Request) error {
	providedToken, err := getBearerToken(r)
	if err != nil {
		return fmt.Errorf("unable to get bearer token: %w", err)
	}

	authToken := strings.TrimSpace(os.Getenv("AUTH_TOKEN"))
	// to enable zero-downtime token rotation
	expiringAuthToken := strings.TrimSpace(os.Getenv("EXPIRING_AUTH_TOKEN"))

	if len(authToken) > 0 && providedToken == authToken {
		return nil
	}

	if len(expiringAuthToken) > 0 && providedToken == expiringAuthToken {
		return nil
	}

	return errors.New("bearer token auth failed")
}

func createGRPCGatewayHandler(temporalClientConn *grpc.ClientConn) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&gateway.JSONPb{
				EmitDefaults: true,
				Indent:       " ",
				OrigName:     false,
			},
		),
	)

	err := workflowservice.RegisterWorkflowServiceHandler(
		context.Background(),
		mux,
		temporalClientConn,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to register workflow service handler: %w", err)
	}

	return mux, nil
}

func main() {
	conn := rpc.CreateGRPCConnection(os.Getenv("TEMPORAL_CLUSTER_HOST")+":7233", nil)
	defer conn.Close()

	handler, err := createGRPCGatewayHandler(conn)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to create gRPC gateway handler: %w", err))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthHeader(r); err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})
	http.ListenAndServe(":10000", nil)
}
