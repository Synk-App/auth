package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"synk/gateway/app/util"
)

type Response struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	Info  any    `json:"info"`
	List  []any  `json:"list"`
}

type ResponseHeader struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func WriteErrorResponse(w http.ResponseWriter, response any, route string, message string, status int) {
	util.LogRoute(route, message)

	jsonResp, _ := json.Marshal(response)

	w.WriteHeader(status)
	w.Write(jsonResp)
}

func WriteSuccessResponse(w http.ResponseWriter, response any) {
	jsonResp, _ := json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func SetJsonContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allowedOriginsMap = map[string]struct{}{
			strings.TrimSuffix(os.Getenv("WEB_ENDPOINT"), "/"):     {},
			strings.TrimSuffix(os.Getenv("GATEWAY_ENDPOINT"), "/"): {},
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		origin := r.Header.Get("Origin")
		if _, ok := allowedOriginsMap[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
