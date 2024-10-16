package user

import (
	"net/http"
	"strconv"

	"github.com/Tutuacs/pkg/password"
	"github.com/Tutuacs/pkg/resolver"
	"github.com/Tutuacs/pkg/routes"
	"github.com/go-playground/validator"
)

type Handler struct {
	subRoute string
}

func NewHandler() *Handler {
	return &Handler{subRoute: "/user"}
}

func (h *Handler) BuildRoutes(router routes.Route) {
	// TODO implement the routes call
	router.NewRoute(routes.POST, h.subRoute, h.create)
	router.NewRoute(routes.GET, h.subRoute, h.list)
	router.NewRoute(routes.GET, h.subRoute+"/{id}", h.getById)
	router.NewRoute(routes.PUT, h.subRoute+"/{id}", h.update)
	router.NewRoute(routes.DELETE, h.subRoute+"/{id}", h.delete)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {

	var payload NewUserDto

	// resolver.GetBody(r, &payload)

	if err := resolver.GetBody(r, &payload); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := resolver.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		resolver.WriteResponse(w, http.StatusBadRequest, errors.Error())
		return
	}

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	// TODO: Implement the auth validation after create

	exists, err := store.GetByEmail(payload.Email)
	if err == nil && exists.ID != 0 {
		resolver.WriteResponse(w, http.StatusConflict, map[string]string{"Error creating the user ": "User already exists"})
		return
	}

	pass, err := password.HashPassword(payload.Password)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error encrypting the users pass": err.Error()})
		return
	}
	payload.Password = pass

	usr, err := store.Create(payload)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error creating the user": err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusCreated, usr)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	// TODO: Implement the auth validation after list

	users, err := store.List()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error listing the users": err.Error()})
		return
	}

	if users == nil {
		resolver.WriteResponse(w, http.StatusOK, []User{})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, users)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {

	param := resolver.GetParam(r, "id")

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		if err != nil {
			resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": err.Error()})
		}
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": "id is invalid"})
		return
	}

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	user, err := store.GetByID(id)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error getting the user": err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, user)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {

	param := resolver.GetParam(r, "id")

	var body UpdateUserDto
	if err := resolver.GetBody(r, body); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error getting the body": err.Error()})
		return
	}

	if err := resolver.Validate.Struct(body); err != nil {
		errors := err.(validator.ValidationErrors)
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error validating the body": errors.Error()})
		return
	}

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		if err != nil {
			resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": err.Error()})
		}
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": "id is invalid"})
		return
	}

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	updated, err := store.Update(id, body)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error updating the user": err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, updated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {

	param := resolver.GetParam(r, "id")

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		if err != nil {
			resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": err.Error()})
		}
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": "id is invalid"})
		return
	}

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	deleted, err := store.Delete(id)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error deleting the user": err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, deleted)
}
