package middlewares

import (
	"campaign-mgmt/app/usecases/dto"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	verifier "github.com/okta/okta-jwt-verifier-golang"
	logger "github.com/sirupsen/logrus"
)

func OktaAuthenticator(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, x-requested-with, origin, X-API-VERSION,Accept-Language")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS, PATCH")

		// OK for all pre-flight requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/campaigns/update-status") || r.Method == "GET" {
			inner.ServeHTTP(w, r)
			return
		}
		
		valid, ctx := validateToken(r)
		if !valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			r := dto.Response{StatusCode: http.StatusUnauthorized, Message: "Not Authorised"}
			err := json.NewEncoder(w).Encode(r)
			if err != nil {
				logger.Errorf("Failed to write response : %v", err)
			}
			return
		}
		if ctx != nil {
			*r = *r.WithContext(ctx)
		}
		inner.ServeHTTP(w, r)
	})
}

func validateToken(r *http.Request) (bool, context.Context) {
	authHeader := r.Header.Get("Authorization")
	var ctx context.Context = nil

	if authHeader == "" {
		return false, ctx
	}
	tokenParts := strings.Split(authHeader, "Bearer ")
	bearerToken := tokenParts[1]

	claimsToValidate := map[string]string{}
	claimsToValidate["aud"] = os.Getenv("OKTA_AUDIENCE")
	jv := verifier.JwtVerifier{
		Issuer:           os.Getenv("OKTA_ISSUER"),
		ClaimsToValidate: claimsToValidate,
	}

	jwt, err := jv.New().VerifyIdToken(bearerToken)
	if err != nil && strings.Contains(err.Error(), "token is not valid:") {
		// The case of legacy token. Letting the next middleware (legacy auth) to handle
		// But we won't be setting the orgId or userId here as we let the legacy auth to set
		// that into the ctx
		return true, ctx
	}

	if err != nil {
		logger.Errorf("JWT token validation failed with error : %v", err)
		return false, ctx
	}

	// Here we need to set the orgId, userId and user into the ctx
	if jwt != nil {
		uId := jwt.Claims["dbpUserId"].(string)
		ctx = context.WithValue(r.Context(), "userId", uId)
		ctx = context.WithValue(ctx, "user", fmt.Sprintf("{\"id\":%v}", uId))
		ctx = context.WithValue(ctx, "organizationId", "2")
	}
	return true, ctx
}
