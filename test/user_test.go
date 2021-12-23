package test

import (
	"bytes"
	"github.com/gorilla/mux"
	"go-example/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUser(t *testing.T) {
	t.Run("returns create user result", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/user", nil)
		request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"id":"1","name":"1","password":"1","status":1}`))
		response := httptest.NewRecorder()

		service.CreateUser(response, request)

		got := response.Body.String()
		want := `Create user success!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})

	t.Run("returns get user result", func(t *testing.T) {
		vars := map[string]string{
			"id": "1",
		}
		request, _ := http.NewRequest(http.MethodGet, "/user", nil)
		request = mux.SetURLVars(request, vars)

		response := httptest.NewRecorder()

		service.GetUser(response, request)

		got := response.Body.String()
		want := `{"id":"1","name":"1","password":"1","status":1}`

		if got != want {
			t.Errorf("got '%s' size '%d', want %s size %d", got, len(got), want, len(want))
		}
	})

	t.Run("returns update user result", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/user", nil)
		request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"id":"1","name":"1","password":"1","status":2}`))
		response := httptest.NewRecorder()

		service.UpdateUser(response, request)

		got := response.Body.String()
		want := `Update user success!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})

	t.Run("returns delete user result", func(t *testing.T) {
		vars := map[string]string{
			"id": "1",
		}
		request, _ := http.NewRequest(http.MethodGet, "/user", nil)
		request = mux.SetURLVars(request, vars)

		response := httptest.NewRecorder()

		service.DeleteUser(response, request)

		got := response.Body.String()
		want := `Delete user success!`

		if got != want {
			t.Errorf("got '%s', want %s", got, want)
		}
	})
}
