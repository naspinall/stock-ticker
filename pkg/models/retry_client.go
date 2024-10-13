package models

import "net/http"

type RetryClient struct {
	http.Client
}
