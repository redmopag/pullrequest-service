package team

import (
	"github.com/avito/internship/pr-service/internal/model"
)

func (dto *TeamDTO) ToModel() *model.Team {
	var users []model.User
	for _, member := range dto.Members {
		user := model.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: dto.TeamName,
			IsActive: member.IsActive,
		}
		users = append(users, user)
	}
	return &model.Team{
		TeamName: dto.TeamName,
		Users:    users,
	}
}

func ModelToDTO(team *model.Team) *TeamDTO {
	var members []TeamMember
	for _, user := range team.Users {
		member := TeamMember{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		}
		members = append(members, member)
	}
	return &TeamDTO{
		TeamName: team.TeamName,
		Members:  members,
	}
}
