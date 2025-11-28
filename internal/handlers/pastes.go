package handlers

import "pastebin/internal/services"

type PasteHandler struct {
	pasteSvc *services.PasteService
}
