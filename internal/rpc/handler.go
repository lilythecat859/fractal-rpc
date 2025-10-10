package rpc

import "net/http"

func NewHandler() http.Handler { return http.DefaultServeMux }
