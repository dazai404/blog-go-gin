package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) authMiddleware(ctx *gin.Context) {
	sessionToken, err := ctx.Cookie("session_cookie")
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	session, err := h.repository.GetSessionByToken(sessionToken)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	err = session.IsExpired()
	if err != nil {
		ctx.SetCookie("session_cookie", "", 0, "/", "localhost", false, true)
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		err = h.repository.DeleteSession(session.ID)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

	user, err := h.repository.GetUserByID(session.UserID)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set("user", user)
}

func (h *Handler) adminMiddleware(ctx *gin.Context) {
	sessionToken, err := ctx.Cookie("session_cookie")
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	session, err := h.repository.GetSessionByToken(sessionToken)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	err = session.IsExpired()
	if err != nil {
		ctx.SetCookie("session_cookie", "", 0, "/", "localhost", false, true)
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		err = h.repository.DeleteSession(session.ID)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

	user, err := h.repository.GetUserByID(session.UserID)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "session is expired")
		return
	}

	if user.Role == "admin" {
		ctx.Set("user", user)
		ctx.Next()
		return
	}

	newErrorResponse(ctx, http.StatusForbidden, "incorrect role")
}
