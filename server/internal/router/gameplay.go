package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrspec7er/matchmind/server/internal/controller"
)

func Router() func(chi.Router) {
	c := &controller.Controller{}

	return func(r chi.Router) {
		r.Get("/{roomId}", c.SendResponse)
		r.Post("/rooms", c.CreateRoom)
	}
}
