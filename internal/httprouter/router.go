package httprouter

import (
	"fmt"
	"net/http"
	"pr-reviewer/internal/httprouter/handlers"
	"pr-reviewer/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(teamSvc *service.TeamService, userSvc *service.UserService, prSvc *service.PRService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	
	hTeam := handlers.NewTeamHandler(teamSvc)
	hUser := handlers.NewUserHandler(userSvc)
	hPR := handlers.NewPRHandler(prSvc)

	r.Post("/team/add", hTeam.AddTeam)
	r.Get("/team/get", hTeam.GetTeam)

	r.Post("/users/setIsActive", hUser.SetIsActive)
	r.Get("/users/getReview", hUser.GetReviewPRs)

	r.Post("/pullRequest/create", hPR.CreatePR)
	r.Post("/pullRequest/merge", hPR.MergePR)
	r.Post("/pullRequest/reassign", hPR.ReassignReviewer)

	r.Get("/health", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(writer, "service is healthy")
	})

	return r
}
