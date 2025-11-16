package user

import (
	"encoding/json"
	"net/http"

	"github.com/avito/internship/pr-service/internal/handler"
	"github.com/avito/internship/pr-service/internal/model"
)

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}

func (controller *UserController) SetIsActive(w http.ResponseWriter, r *http.Request) error {
	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return model.ErrBadJSONRequest
	}
	user, err := controller.userService.SetActiveStatus(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		return err
	}
	handler.WriteJSONResponse(w, http.StatusOK, map[string]any{"user:": ToUserResponse(user)})
	return nil
}

func (controller *UserController) GetUsersPullRequests(w http.ResponseWriter, r *http.Request) error {
	userId := r.URL.Query().Get("user_id")
	pullRequests, err := controller.userService.GetUsersPullRequests(r.Context(), userId)
	if err != nil {
		return err
	}

	prResponse := make([]PullRequestShort, 0, len(pullRequests))
	for _, pr := range pullRequests {
		prResponse = append(prResponse, *ToPullRequestShort(&pr))
	}

	response := map[string]any{
		"user_id":       userId,
		"pull_requests": prResponse,
	}

	handler.WriteJSONResponse(w, http.StatusOK, response)
	return nil
}
