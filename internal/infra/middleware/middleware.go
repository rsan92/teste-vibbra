package middleware

import "net/http"

type IMiddlewareFunction func(next http.HandlerFunc) http.HandlerFunc
