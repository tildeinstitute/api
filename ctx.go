package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
)

type ipCtxKey int

const ctxKey ipCtxKey = iota

// Creates a new context containing the request's originating IP.
// Preference: X-Forwarded-For, X-Real-IP, RemoteAddress
func newCtxUserIP(ctx context.Context, r *http.Request) context.Context {
	split := strings.Split(r.RemoteAddr, ":")
	uip := split[0]

	xFwdFor := http.CanonicalHeaderKey("X-Forwarded-For")
	if _, ok := r.Header[xFwdFor]; ok {
		fwdaddr := r.Header[xFwdFor]
		split := strings.Split(fwdaddr[len(fwdaddr)-1], ":")
		uip = split[0]

		return context.WithValue(ctx, ctxKey, uip)
	}

	xRealIP := http.CanonicalHeaderKey("X-Real-IP")
	if _, ok := r.Header[xRealIP]; ok {
		realip := r.Header[xRealIP]
		split := strings.Split(realip[len(realip)-1], ":")
		uip = split[0]
	}

	return context.WithValue(ctx, ctxKey, uip)
}

// Yoinks the IP address from the request's context
func getIPFromCtx(ctx context.Context) net.IP {
	uip, ok := ctx.Value(ctxKey).(string)
	if !ok {
		log.Printf("Couldn't retrieve IP From Request\n")
	}
	return net.ParseIP(uip)
}

// Middleware to yeet the remote IP address into the request struct
func ipMiddleware(hop http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtxUserIP(r.Context(), r)
		hop.ServeHTTP(w, r.WithContext(ctx))
	})
}
