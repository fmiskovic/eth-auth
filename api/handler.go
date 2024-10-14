package api

import (
	"github.com/fmiskovic/eth-auth/store"
)

type handler struct {
	secret string // secret for signing JWT tokens
	store  store.Storer
}

func newHandler(store store.Storer, secret string) handler {
	return handler{store: store, secret: secret}
}
