package handler

import (
	"fmt"
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/usecases"
	dto "portfolio/dto/skill"
	"portfolio/logger"
	"portfolio/shared"
	"strconv"
	"time"
)

type skillHandler struct {
	AbstractHandler
	skillUseCase *usecases.SkillUseCase
	logger       *logger.Logger
}

func NewSkillHandler(settingUseCase *usecases.SettingUseCase, skillUseCase *usecases.SkillUseCase, logger *logger.Logger) []*routes.NamedRoute {
	skillHandler := skillHandler{
		AbstractHandler: AbstractHandler{
			settingUseCase: settingUseCase,
		},
		skillUseCase: skillUseCase,
		logger:       logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetSkillsHandler",
			Pattern: "GET /skills",
			Handler: skillHandler.GetSkills,
		},
		{
			Name:    "GetSkillHandler",
			Pattern: "GET /skills/{id}",
			Handler: skillHandler.GetSkill,
		},
	}
}

// GetSkills
//
//	@Summary		Get all skills
//	@Description	Retrieve all skills for the portfolio
//	@Tags			Skills
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.SkillListResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/skills [get]
func (sh *skillHandler) GetSkills(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolioOwnerID, err := utils.GetPortfolioOwnerID(sh.settingUseCase, ctx, w)

	if err != nil {
		sh.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	skills, err := sh.skillUseCase.GetSkillsByUserID(ctx, portfolioOwnerID)
	if err != nil {
		sh.logger.Error("Failed to get skills: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := dto.FromSkillsEntityToResponse(skills, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	if response == nil {
		sh.logger.Error("Skills response is nil for user ID %d", portfolioOwnerID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Skills", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetSkill
//
//	@Summary		Get a specific skill
//	@Description	Retrieve a specific skill by ID
//	@Tags			Skills
//	@Produce		json
//	@Param			id	path		int	true	"Skill ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.SkillResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/skills/{id} [get]
func (sh *skillHandler) GetSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	skillIDStr := r.PathValue("id")
	skillID, err := strconv.Atoi(skillIDStr)
	if err != nil || skillID <= 0 {
		sh.logger.Error("Invalid skill ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skill ID", "id", &err))
		return
	}

	skill, err := sh.skillUseCase.GetSkillByID(ctx, skillID)
	if err != nil {
		sh.logger.Error("Failed to get skill by ID %d: %v", skillID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	portfolioOwnerID, err := utils.GetPortfolioOwnerID(sh.settingUseCase, ctx, w)

	if err != nil {
		sh.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	if skill.UserID != portfolioOwnerID {
		sh.logger.Error("Unauthorized access to skill %d", skillID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Skill", fmt.Sprint(skillID)))
		return
	}

	response := dto.FromSkillEntityToResponse(skill, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	if response == nil {
		sh.logger.Error("Failed to create skill response for ID %d", skillID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Skill", fmt.Sprint(skillID)))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}
