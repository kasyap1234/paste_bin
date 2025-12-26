package handlers

import (
	"net/http"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"pastebin/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ProfileHandler struct {
	profileService *services.ProfileService
	logger         *zerolog.Logger
}

func NewProfileHandler(profileService *services.ProfileService, logger *zerolog.Logger) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         logger,
	}
}

func (p *ProfileHandler) GetProfileHandler(c echo.Context) error {
	userID, err := auth.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		p.logger.Err(err).Msg("failed to get user id from context")
		return utils.SendError(c, http.StatusInternalServerError, "failed to get user id from context")
	}
	user, err := p.profileService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		p.logger.Err(err).Msg("failed to get profile")
		return utils.SendError(c, http.StatusInternalServerError, "failed to get profile")
	}
	return utils.SendSuccess(c, http.StatusOK, user, "profile retrieved successfully")
}

func (p *ProfileHandler) UpdateProfileHandler(c echo.Context) error {
	userID, err := auth.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		p.logger.Err(err).Msg("failed to get userID from context ")
		return utils.SendError(c, http.StatusInternalServerError, "failed to get userID from context")
	}
	var patchProfile models.PatchProfile
	if err := c.Bind(&patchProfile); err != nil {
		p.logger.Err(err).Msg("failed to bind patch profile")
		return utils.SendError(c, http.StatusBadRequest, "invalid request")
	}
	user, err := p.profileService.UpdateProfile(c.Request().Context(), userID, &patchProfile)
	if err != nil {
		p.logger.Err(err).Msg("failed to update profile")
		return utils.SendError(c, http.StatusInternalServerError, "failed to update profile")
	}
	return utils.SendSuccess(c, http.StatusOK, user, "profile updated successfully")
}
