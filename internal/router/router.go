package router

import (
	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type routing struct {
	Router ports.Routing
}

func New(r ports.Routing) *routing {
	return &routing{
		Router: r,
	}
}
