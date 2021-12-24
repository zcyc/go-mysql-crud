package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-example/service"

	"github.com/gorilla/mux"
)

func TestUser(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/user", nil)
		request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"id":"1","name":"1","password":"1","status":1}`))
		response := httptest.NewRecorder()

		service.CreateUser(response, request)

		got := response.Body.String()
		want := `Create user success!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		vars := map[string]string{
			"id": "1",
		}
		request, _ := http.NewRequest(http.MethodGet, "/user", nil)
		request = mux.SetURLVars(request, vars)

		response := httptest.NewRecorder()

		service.GetUser(response, request)

		got := response.Body.String()
		want := `{"result":true,"msg":"success","code":200,"data":{"id":"1","name":"1","password":"1","status":1}}`

		if got != want {
			t.Errorf("got '%s' size '%d', want %s size %d", got, len(got), want, len(want))
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/user", nil)
		request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"id":"1","name":"1","password":"1","status":2}`))
		response := httptest.NewRecorder()

		service.UpdateUser(response, request)

		got := response.Body.String()
		want := `Update user success!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		vars := map[string]string{
			"id": "1",
		}
		request, _ := http.NewRequest(http.MethodDelete, "/user", nil)
		request = mux.SetURLVars(request, vars)

		response := httptest.NewRecorder()

		service.DeleteUser(response, request)

		got := response.Body.String()
		want := `Delete user success!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})

	t.Run("CreateUserByNilName", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/user", nil)
		request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"id":"1","password":"1","status":1}`))
		response := httptest.NewRecorder()

		service.CreateUser(response, request)

		got := response.Body.String()
		want := `Name/Password is empty!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})

	t.Run("GetUserByNilId", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/user", nil)

		response := httptest.NewRecorder()

		service.GetUser(response, request)

		got := response.Body.String()
		want := `Get failed,id is nil!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})

	t.Run("UpdateUserByNilId", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/user", nil)
		request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"name":"1","password":"1","status":2}`))
		response := httptest.NewRecorder()

		service.UpdateUser(response, request)

		got := response.Body.String()
		want := `ID/Name/Password is empty!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})

	t.Run("DeleteUserByNilId", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/user", nil)

		response := httptest.NewRecorder()

		service.DeleteUser(response, request)

		got := response.Body.String()
		want := `Delete failed,id is nil!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})
}
