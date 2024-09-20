package person

import (
	"errors"
	"github.com/Erlendum/rsoi-lab-01/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type handlerTestFields struct {
	storage *Mockstorage
}

func createHandlerTestFields(ctrl *gomock.Controller) *handlerTestFields {
	return &handlerTestFields{
		storage: NewMockstorage(ctrl),
	}
}

func getPointerOnString(s string) *string {
	return &s
}

func getPointerOnInt(i int) *int {
	return &i
}

func Test_CreatePerson(t *testing.T) {
	type fields struct {
		reqBody                string
		expectedHTTPCode       int
		expectedLocationHeader string
	}

	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	tests := []struct {
		name    string
		fields  fields
		Prepare func(fields *handlerTestFields)
	}{
		{
			name: "http-code 400: wrong body",
			fields: fields{
				expectedHTTPCode:       http.StatusBadRequest,
				reqBody:                ``,
				expectedLocationHeader: ``,
			},

			Prepare: func(fields *handlerTestFields) {
			},
		},
		{
			name: "http-code 400: empty body",
			fields: fields{
				expectedHTTPCode:       http.StatusBadRequest,
				reqBody:                `{}`,
				expectedLocationHeader: ``,
			},

			Prepare: func(fields *handlerTestFields) {
			},
		},
		{
			name: "http-code 500: storage error",
			fields: fields{
				expectedHTTPCode:       http.StatusInternalServerError,
				reqBody:                `{"name": "test", "age": 1, "address": "test", "work": "test"}`,
				expectedLocationHeader: ``,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().CreatePerson(gomock.Any(), gomock.Any()).Return(0, errors.New(""))
			},
		},
		{
			name: "http-code 201",
			fields: fields{
				expectedHTTPCode:       http.StatusCreated,
				reqBody:                `{"name": "test", "age": 1, "address": "test", "work": "test"}`,
				expectedLocationHeader: `/api/v1/persons/1`,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().CreatePerson(gomock.Any(), gomock.Any()).Return(1, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			testFields := createHandlerTestFields(ctrl)
			tt.Prepare(testFields)

			h := &handler{storage: testFields.storage}

			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.fields.reqBody))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.CreatePerson(c)

			require.NoError(t, err)
			require.Equal(t, tt.fields.expectedHTTPCode, rec.Code)
			require.Equal(t, tt.fields.expectedLocationHeader, rec.Header().Get("Location"))
		})
	}
}

func Test_UpdatePerson(t *testing.T) {
	type fields struct {
		id               string
		reqBody          string
		expectedHTTPCode int
	}

	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	tests := []struct {
		name    string
		fields  fields
		Prepare func(fields *handlerTestFields)
	}{
		{
			name: "http-code 400: wrong id",
			fields: fields{
				expectedHTTPCode: http.StatusBadRequest,
				id:               "test",
				reqBody:          `{"name": "test", "age": 1, "address": "test", "work": "test"}`,
			},

			Prepare: func(fields *handlerTestFields) {
			},
		},
		{
			name: "http-code 400: wrong body",
			fields: fields{
				expectedHTTPCode: http.StatusBadRequest,
				id:               "1",
				reqBody:          ``,
			},

			Prepare: func(fields *handlerTestFields) {
			},
		},
		{
			name: "http-code 400: empty body",
			fields: fields{
				expectedHTTPCode: http.StatusBadRequest,
				id:               "1",
				reqBody:          `{}`,
			},

			Prepare: func(fields *handlerTestFields) {
			},
		},
		{
			name: "http-code 500: storage error",
			fields: fields{
				expectedHTTPCode: http.StatusInternalServerError,
				id:               "1",
				reqBody:          `{"name": "test", "age": 1, "address": "test", "work": "test"}`,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().UpdatePerson(gomock.Any(), 1, gomock.Any()).Return(errors.New(""))
			},
		},
		{
			name: "http-code 200",
			fields: fields{
				expectedHTTPCode: http.StatusOK,
				id:               "1",
				reqBody:          `{"name": "test", "age": 1, "address": "test", "work": "test"}`,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().UpdatePerson(gomock.Any(), 1, gomock.Any()).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			testFields := createHandlerTestFields(ctrl)
			tt.Prepare(testFields)

			h := &handler{storage: testFields.storage}

			req := httptest.NewRequest(http.MethodPatch, "/test", strings.NewReader(tt.fields.reqBody))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.fields.id)

			err := h.UpdatePerson(c)

			require.NoError(t, err)
			require.Equal(t, tt.fields.expectedHTTPCode, rec.Code)
		})
	}
}

func Test_DeletePerson(t *testing.T) {
	type fields struct {
		id               string
		expectedHTTPCode int
	}

	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	tests := []struct {
		name    string
		fields  fields
		Prepare func(fields *handlerTestFields)
	}{
		{
			name: "http-code 400: wrong id",
			fields: fields{
				expectedHTTPCode: http.StatusBadRequest,
				id:               "test",
			},

			Prepare: func(fields *handlerTestFields) {
			},
		},
		{
			name: "http-code 500: storage error",
			fields: fields{
				expectedHTTPCode: http.StatusInternalServerError,
				id:               "1",
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().DeletePerson(gomock.Any(), 1).Return(false, errors.New(""))
			},
		},
		{
			name: "http-code 404",
			fields: fields{
				expectedHTTPCode: http.StatusNotFound,
				id:               "1",
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().DeletePerson(gomock.Any(), 1).Return(false, nil)
			},
		},
		{
			name: "http-code 204",
			fields: fields{
				expectedHTTPCode: http.StatusNoContent,
				id:               "1",
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().DeletePerson(gomock.Any(), 1).Return(true, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			testFields := createHandlerTestFields(ctrl)
			tt.Prepare(testFields)

			h := &handler{storage: testFields.storage}

			req := httptest.NewRequest(http.MethodDelete, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.fields.id)

			err := h.DeletePerson(c)

			require.NoError(t, err)
			require.Equal(t, tt.fields.expectedHTTPCode, rec.Code)
		})
	}
}

func Test_GetPerson(t *testing.T) {
	type fields struct {
		id                   string
		expectedHTTPCode     int
		expectedResponseBody string
	}

	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	tests := []struct {
		name    string
		fields  fields
		Prepare func(fields *handlerTestFields)
	}{
		{
			name: "http-code 400: wrong id",
			fields: fields{
				expectedHTTPCode:     http.StatusBadRequest,
				id:                   "test",
				expectedResponseBody: ``,
			},

			Prepare: func(fields *handlerTestFields) {
			},
		},
		{
			name: "http-code 500: storage error",
			fields: fields{
				expectedHTTPCode:     http.StatusInternalServerError,
				id:                   "1",
				expectedResponseBody: ``,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().GetPerson(gomock.Any(), 1).Return(Person{}, errors.New(""))
			},
		},
		{
			name: "http-code 404",
			fields: fields{
				expectedHTTPCode:     http.StatusNotFound,
				id:                   "1",
				expectedResponseBody: ``,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().GetPerson(gomock.Any(), 1).Return(Person{}, nil)
			},
		},
		{
			name: "http-code 200",
			fields: fields{
				expectedHTTPCode: http.StatusOK,
				id:               "1",
				expectedResponseBody: `{"id":1,"name":"test","age":2,"address":"testaddress","work":"testwork"}
`,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().GetPerson(gomock.Any(), 1).Return(Person{
					ID:      getPointerOnInt(1),
					Name:    getPointerOnString("test"),
					Address: getPointerOnString("testaddress"),
					Age:     getPointerOnInt(2),
					Work:    getPointerOnString("testwork"),
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			testFields := createHandlerTestFields(ctrl)
			tt.Prepare(testFields)

			h := &handler{storage: testFields.storage}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.fields.id)

			err := h.GetPerson(c)

			require.NoError(t, err)
			require.Equal(t, tt.fields.expectedHTTPCode, rec.Code)
			if tt.fields.expectedHTTPCode == http.StatusOK {
				body, err := io.ReadAll(rec.Result().Body)
				require.NoError(t, err)
				require.Equal(t, tt.fields.expectedResponseBody, string(body))
			}
		})
	}
}

func Test_GetPersons(t *testing.T) {
	type fields struct {
		expectedHTTPCode     int
		expectedResponseBody string
	}

	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	tests := []struct {
		name    string
		fields  fields
		Prepare func(fields *handlerTestFields)
	}{
		{
			name: "http-code 500: storage error",
			fields: fields{
				expectedHTTPCode:     http.StatusInternalServerError,
				expectedResponseBody: ``,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().GetPersons(gomock.Any()).Return([]Person{{}}, errors.New(""))
			},
		},
		{
			name: "http-code 200",
			fields: fields{
				expectedHTTPCode: http.StatusOK,
				expectedResponseBody: `[{"id":1,"name":"test","age":2,"address":"testaddress","work":"testwork"}]
`,
			},

			Prepare: func(fields *handlerTestFields) {
				fields.storage.EXPECT().GetPersons(gomock.Any()).Return([]Person{{
					ID:      getPointerOnInt(1),
					Name:    getPointerOnString("test"),
					Address: getPointerOnString("testaddress"),
					Age:     getPointerOnInt(2),
					Work:    getPointerOnString("testwork"),
				}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			testFields := createHandlerTestFields(ctrl)
			tt.Prepare(testFields)

			h := &handler{storage: testFields.storage}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.GetPersons(c)

			require.NoError(t, err)
			require.Equal(t, tt.fields.expectedHTTPCode, rec.Code)
			if tt.fields.expectedHTTPCode == http.StatusOK {
				body, err := io.ReadAll(rec.Result().Body)
				require.NoError(t, err)
				require.Equal(t, tt.fields.expectedResponseBody, string(body))
			}
		})
	}
}
