package team

import (
	"encoding/json"
	"net/http"

	"github.com/avito/internship/pr-service/internal/handler"
	"github.com/avito/internship/pr-service/internal/model"
)

type TeamController struct {
	teamService TeamService
}

func NewTeamController(teamService TeamService) *TeamController {
	return &TeamController{
		teamService: teamService,
	}
}

func (c *TeamController) CreateTeam(w http.ResponseWriter, r *http.Request) error {
	var createTeamReq TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&createTeamReq); err != nil {
		return model.ErrBadJSONRequest
	}
	team := createTeamReq.ToModel()
	createdTeam, err := c.teamService.CreateTeam(r.Context(), team)
	if err != nil {
		return err
	}
	response := map[string]any{"team": ModelToDTO(createdTeam)}
	handler.WriteJSONResponse(w, http.StatusCreated, response)
	return nil
}

func (c *TeamController) GetTeamByName(w http.ResponseWriter, r *http.Request) error {
	teamName := r.URL.Query().Get("team_name")
	team, err := c.teamService.GetTeamByName(r.Context(), teamName)
	if err != nil {
		return err
	}
	handler.WriteJSONResponse(w, http.StatusOK, ModelToDTO(team))
	return nil
}
