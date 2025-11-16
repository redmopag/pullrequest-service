package pullrequest

import (
	"encoding/json"
	"net/http"

	"github.com/avito/internship/pr-service/internal/handler"
	"github.com/avito/internship/pr-service/internal/model"
)

type PullRequestController struct {
	pullRequestService PullRequestService
}

func NewPullRequestController(pullRequestService PullRequestService) *PullRequestController {
	return &PullRequestController{pullRequestService: pullRequestService}
}

func (c *PullRequestController) Create(w http.ResponseWriter, r *http.Request) error {
	var req CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return model.ErrBadJSONRequest
	}

	pr, err := c.pullRequestService.Create(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		return err
	}

	response := map[string]any{"pr": ToPullRequestDTO(pr)}
	handler.WriteJSONResponse(w, http.StatusCreated, response)
	return nil
}

func (c *PullRequestController) Merge(w http.ResponseWriter, r *http.Request) error {
	var req MergePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return model.ErrBadJSONRequest
	}

	pr, err := c.pullRequestService.Merge(r.Context(), req.PullRequestID)
	if err != nil {
		return err
	}

	handler.WriteJSONResponse(w, http.StatusOK, ToPullRequestDTO(pr))
	return nil
}

func (c *PullRequestController) Reassign(w http.ResponseWriter, r *http.Request) error {
	var req ReassignPullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return model.ErrBadJSONRequest
	}

	res, err := c.pullRequestService.Reassign(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		return err
	}

	responseDTO := ReassignResponseDTO{
		PullRequest: *ToPullRequestDTO(res.PullRequest),
		ReplacedBy:  res.ReplacedBy,
	}

	handler.WriteJSONResponse(w, http.StatusOK, responseDTO)
	return nil
}
