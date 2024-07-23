package router

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"just_for_fun/internal/server/response"
	"just_for_fun/internal/server/router/structures"
	"just_for_fun/internal/storage"
	"just_for_fun/internal/storage/db/models"
	"just_for_fun/internal/storage/db/repo"
	"net/http"

	"github.com/go-chi/render"
)

type UserRouter struct {
	router *MainRouter
	dbUser DBUserProvider
}

func NewUserRouter(r *MainRouter, storage *storage.Storage) *UserRouter {
	ur := &UserRouter{router: r, dbUser: storage.DB.User}
	ur.addRoute()
	return ur
}

type DBUserProvider interface {
	Set(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
}

func (ur *UserRouter) addRoute() {
	userPattern := fmt.Sprintf("%s/user", ur.router.BasePattern)

	ur.router.mr.Route(userPattern, func(r chi.Router) {
		r.Post("/", ur.SetUpdateUser())
	})
}

// SetUpdateUser
// @Summary Set or Update user
// @Tags users
// @Description Set or Update user
// @Accept  json
// @Produce  json
// @Param UserInfo body structures.UserSetRequest true "User"
// @Success 201 {object} response.UserResponse "User created"
// @Success 202 {object} response.UserResponse "User updated"
// @Failure 400 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Server error"
// @Router /api/user/ [post]
func (ur *UserRouter) SetUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req structures.UserSetRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			ur.router.log.Error("Error decoding request", zap.Error(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		user := &models.User{
			TgID:      req.TgID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			UserName:  req.UserName,
		}

		newUser, err := ur.dbUser.Set(user)
		if err != nil {
			if errors.Is(err, repo.ErrIsExists) {
				updatedUser, updatedErr := ur.dbUser.Update(user)
				if updatedErr != nil {
					ur.router.log.Error("Error updating user", zap.Error(updatedErr))
					render.Status(r, http.StatusInternalServerError)
					render.JSON(w, r, response.Error(updatedErr.Error()))
					return
				}

				render.Status(r, http.StatusAccepted)
				render.JSON(w, r, response.UserResponse{
					Response: response.OK(),
					User: structures.UserShow{
						TgID:      updatedUser.TgID,
						FirstName: updatedUser.FirstName,
						LastName:  updatedUser.LastName,
						UserName:  updatedUser.UserName,
					},
				})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.UserResponse{
			Response: response.OK(),
			User: structures.UserShow{
				TgID:      newUser.TgID,
				FirstName: newUser.FirstName,
				LastName:  newUser.LastName,
				UserName:  newUser.UserName,
			},
		})
		return
	}
}

func (ur *UserRouter) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
