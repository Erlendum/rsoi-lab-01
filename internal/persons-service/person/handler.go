package person

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

//go:generate mockgen -source=handler.go  -destination=handler_mocks.go -self_package=github.com/Erlendum/rsoi-lab-01/internal/persons-service/person -package=person

type storage interface {
	CreatePerson(ctx context.Context, person Person) (int, error)
	UpdatePerson(ctx context.Context, id int, person *Person) error
	DeletePerson(ctx context.Context, id int) (bool, error)
	GetPersons(ctx context.Context) ([]Person, error)
	GetPerson(ctx context.Context, id int) (Person, error)
}

type handler struct {
	storage storage
}

func NewHandler(storage storage) *handler {
	return &handler{storage: storage}
}

func (h *handler) Register(echo *echo.Echo) {
	api := echo.Group("/api/v1")

	api.GET("/persons/:id", h.GetPerson)
	api.GET("/persons", h.GetPersons)
	api.POST("/persons", h.CreatePerson)
	api.PATCH("/persons/:id", h.UpdatePerson)
	api.DELETE("/persons/:id", h.DeletePerson)
}

func (h *handler) CreatePerson(c echo.Context) error {
	type createPersonRequest struct {
		Name    *string `json:"name" validate:"required"`
		Age     *int    `json:"age"`
		Address *string `json:"address"`
		Work    *string `json:"work"`
	}

	req := &createPersonRequest{}
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("reading request body error")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "unmarshalling error",
		})
	}

	if err = json.Unmarshal(body, &req); err != nil {
		log.Error().Err(err).Msg("unmarshalling error")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "unmarshalling error",
		})
	}

	if err = c.Validate(req); err != nil {
		log.Error().Err(err).Msg("validation error")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "validation error",
		})
	}

	p := Person{
		Name:    req.Name,
		Address: req.Address,
		Age:     req.Age,
		Work:    req.Work,
	}

	id, err := h.storage.CreatePerson(c.Request().Context(), p)
	if err != nil {
		log.Error().Err(err).Msg("creating person error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": "creating person error",
		})
	}

	c.Response().Header().Set("Location", "/api/v1/persons/"+strconv.Itoa(id))

	return c.NoContent(http.StatusCreated)
}

func (h *handler) UpdatePerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "wrong id",
		})
	}

	type createPersonRequest struct {
		Name    *string `json:"name" validate:"required"`
		Age     *int    `json:"age"`
		Address *string `json:"address"`
		Work    *string `json:"work"`
	}

	req := &createPersonRequest{}
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("reading request body error")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "unmarshalling error",
		})
	}

	if err = json.Unmarshal(body, &req); err != nil {
		log.Error().Err(err).Msg("unmarshalling error")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "unmarshalling error",
		})
	}

	if err = c.Validate(req); err != nil {
		log.Error().Err(err).Msg("validation error")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "validation error",
		})
	}

	p := Person{
		Name:    req.Name,
		Address: req.Address,
		Age:     req.Age,
		Work:    req.Work,
	}

	err = h.storage.UpdatePerson(c.Request().Context(), id, &p)
	if err != nil {
		log.Error().Err(err).Msg("updating person error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": "updating person error",
		})
	}

	type updatePersonResponse struct {
		ID      *int    `json:"id"`
		Name    *string `json:"name"`
		Age     *int    `json:"age"`
		Address *string `json:"address"`
		Work    *string `json:"work"`
	}

	resp := updatePersonResponse{
		ID:      p.ID,
		Name:    p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *handler) DeletePerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "wrong id",
		})
	}

	isDeleted, err := h.storage.DeletePerson(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("deleting person error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": "deleting person error",
		})
	}

	if !isDeleted {
		log.Error().Err(err).Msgf("person with id = %d not found", id)
		return c.JSON(http.StatusNotFound, echo.Map{
			"errors": "person not found",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) GetPerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "wrong id",
		})
	}

	p, err := h.storage.GetPerson(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("getting person error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": "getting person error",
		})
	}

	if p.ID == nil {
		log.Error().Err(err).Msg("person with not found")
		return c.JSON(http.StatusNotFound, echo.Map{
			"errors": "person not found",
		})
	}

	type getPersonResponse struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Address string `json:"address"`
		Work    string `json:"work"`
	}

	resp := getPersonResponse{
		ID:      *p.ID,
		Name:    *p.Name,
		Age:     *p.Age,
		Address: *p.Address,
		Work:    *p.Work,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *handler) GetPersons(c echo.Context) error {
	persons, err := h.storage.GetPersons(c.Request().Context())
	if err != nil {
		log.Error().Err(err).Msg("getting persons error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": "getting persons error",
		})
	}

	type getPersonResponse struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Address string `json:"address"`
		Work    string `json:"work"`
	}

	personsResp := make([]getPersonResponse, len(persons), len(persons))
	for i, p := range persons {
		personsResp[i] = getPersonResponse{
			ID:      *p.ID,
			Name:    *p.Name,
			Age:     *p.Age,
			Address: *p.Address,
			Work:    *p.Work,
		}
	}

	return c.JSON(http.StatusOK, personsResp)
}
