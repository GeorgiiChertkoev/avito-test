package httpRouter

import (
	"pr-reviewer/internal/httpRouter/handlers"
	"pr-reviewer/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(teamSvc *service.TeamService, userSvc *service.UserService, prSvc *service.PRService) *chi.Mux {
	r := chi.NewRouter()

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

	return r
}
