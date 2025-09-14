package app

import (
	"github.com/aidapedia/jabberwock/internal/interface/http"
	"github.com/google/wire"
)

var (
	httpSet = wire.NewSet(
		http.NewHTTPService,
	)
)
