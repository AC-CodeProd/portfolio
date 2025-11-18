package usecases

import (
	"context"
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
)

type SkillUseCase struct {
	skillRepo interfaces.SkillRepository
	userRepo  interfaces.UserRepository
	logger    *logger.Logger
}

func NewSkillUseCase(skillRepo interfaces.SkillRepository, userRepo interfaces.UserRepository, logger *logger.Logger) *SkillUseCase {
	return &SkillUseCase{
		skillRepo: skillRepo,
		userRepo:  userRepo,
		logger:    logger,
	}
}

func (uc *SkillUseCase) CreateSkill(ctx context.Context, skill *entities.Skill) (*entities.Skill, error) {
	if !skill.HasRequiredFields() {
		uc.logger.Error("Required fields are missing for skill: %v", skill)
		return nil, domain.NewValidationError("skill", "skill name, level (1-5), and user ID are required", nil)
	}

	if !skill.IsValidLevel() {
		uc.logger.Error("Invalid skill level for skill: %v", skill)
		return nil, domain.NewValidationError("level", "skill level must be between 1 and 5", nil)
	}

	userExists, err := uc.userRepo.ExistsByID(ctx, skill.UserID)
	if err != nil {
		uc.logger.Error("Failed to check if user exists: %v", err)
		return nil, domain.NewInternalError("failed to validate user", err)
	}
	if !userExists {
		uc.logger.Error("User %d does not exist", skill.UserID)
		return nil, domain.NewNotFoundError("User", fmt.Sprint(skill.UserID))
	}

	exists, err := uc.skillRepo.ExistsByNameAndUserID(ctx, skill.Name, skill.UserID)
	if err != nil {
		uc.logger.Error("Failed to check if skill exists: %v", err)
		return nil, domain.NewInternalError("failed to check skill existence", err)
	}
	if exists {
		uc.logger.Error("Skill %s already exists for user %d", skill.Name, skill.UserID)
		return nil, domain.NewAlreadyExistsError("Skill", skill.Name)
	}

	createdSkill, err := uc.skillRepo.Create(ctx, skill)
	if err != nil {
		uc.logger.Error("Failed to create skill: %v", err)
		return nil, domain.NewInternalError("failed to create skill", err)
	}

	return createdSkill, nil
}

func (uc *SkillUseCase) GetSkillByID(ctx context.Context, skillID int) (*entities.Skill, error) {
	if skillID <= 0 {
		uc.logger.Error("Skill ID is required")
		return nil, domain.NewValidationError("skillID", "skill ID must be positive", nil)
	}

	skill, err := uc.skillRepo.GetByID(ctx, skillID)
	if err != nil {
		uc.logger.Error("Failed to get skill by ID %d: %v", skillID, err)
		return nil, domain.NewInternalError("failed to retrieve skill", err)
	}

	if skill == nil {
		uc.logger.Error("Skill not found: %d", skillID)
		return nil, domain.NewNotFoundError("Skill", fmt.Sprint(skillID))
	}

	return skill, nil
}

func (uc *SkillUseCase) GetSkillsByUserID(ctx context.Context, userID int) ([]*entities.Skill, error) {
	if userID <= 0 {
		uc.logger.Error("User ID is required")
		return nil, domain.NewValidationError("userID", "user ID must be positive", nil)
	}

	skills, err := uc.skillRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get skills for user %d: %v", userID, err)
		return nil, domain.NewInternalError("failed to retrieve skills", err)
	}

	return skills, nil
}

func (uc *SkillUseCase) UpdateSkill(ctx context.Context, skillID int, skill *entities.Skill) (*entities.Skill, error) {
	if !skill.HasRequiredFields() {
		uc.logger.Error("Required fields are missing for skill: %v", skill)
		return nil, domain.NewValidationError("skill", "skill name, level (1-5), and user ID are required", nil)
	}

	if !skill.IsValidLevel() {
		uc.logger.Error("Invalid skill level for skill: %v", skill)
		return nil, domain.NewValidationError("level", "skill level must be between 1 and 5", nil)
	}

	exists, err := uc.skillRepo.ExistsByID(ctx, skill.SkillID)
	if err != nil {
		uc.logger.Error("Failed to check if skill exists: %v", err)
		return nil, domain.NewInternalError("failed to check skill existence", err)
	}
	if !exists {
		uc.logger.Error("Skill not found: %d", skill.SkillID)
		return nil, domain.NewNotFoundError("Skill", fmt.Sprint(skill.SkillID))
	}

	skill.MarkAsUpdated()
	updatedSkill, err := uc.skillRepo.Update(ctx, skillID, skill)
	if err != nil {
		uc.logger.Error("Failed to update skill: %v", err)
		return nil, domain.NewInternalError("failed to update skill", err)
	}

	return updatedSkill, nil
}

func (uc *SkillUseCase) PatchSkill(ctx context.Context, skillID int, patchData *entities.Skill) (*entities.Skill, error) {
	if skillID <= 0 {
		uc.logger.Error("Skill ID is required")
		return nil, domain.NewValidationError("skillID", "skill ID must be positive", nil)
	}

	existingSkill, err := uc.skillRepo.GetByID(ctx, skillID)
	if err != nil {
		uc.logger.Error("Failed to get skill: %v", err)
		return nil, domain.NewInternalError("failed to retrieve skill", err)
	}

	if existingSkill == nil {
		uc.logger.Error("Skill not found: %d", skillID)
		return nil, domain.NewNotFoundError("Skill", fmt.Sprint(skillID))
	}

	if patchData.Name != "" {
		existingSkill.Name = patchData.Name
	}
	if patchData.Level > 0 {
		if !patchData.IsValidLevel() {
			uc.logger.Error("Invalid skill level for skill: %v", patchData)
			return nil, domain.NewValidationError("level", "skill level must be between 1 and 5", nil)
		}
		existingSkill.Level = patchData.Level
	}

	existingSkill.MarkAsUpdated()
	patchedSkill, err := uc.skillRepo.Patch(ctx, skillID, existingSkill)
	if err != nil {
		uc.logger.Error("Failed to patch skill: %v", err)
		return nil, domain.NewInternalError("failed to patch skill", err)
	}

	return patchedSkill, nil
}

func (uc *SkillUseCase) DeleteSkill(ctx context.Context, skillID int) error {
	if skillID <= 0 {
		uc.logger.Error("Skill ID is required")
		return domain.NewValidationError("skillID", "skill ID must be positive", nil)
	}

	exists, err := uc.skillRepo.ExistsByID(ctx, skillID)
	if err != nil {
		uc.logger.Error("Failed to check if skill exists: %v", err)
		return domain.NewInternalError("failed to check skill existence", err)
	}
	if !exists {
		uc.logger.Error("Skill not found: %d", skillID)
		return domain.NewNotFoundError("Skill", fmt.Sprint(skillID))
	}

	err = uc.skillRepo.Delete(ctx, skillID)
	if err != nil {
		uc.logger.Error("Failed to delete skill: %v", err)
		return domain.NewInternalError("failed to delete skill", err)
	}

	return nil
}

func (uc *SkillUseCase) GetAllSkills(ctx context.Context) ([]*entities.Skill, error) {
	skills, err := uc.skillRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all skills: %v", err)
		return nil, domain.NewInternalError("failed to retrieve skills", err)
	}

	return skills, nil
}
