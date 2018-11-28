package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homework-juan3674-1532739/servers/gateway/models/users"
	"homework-juan3674-1532739/servers/gateway/sessions"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createFakeContext() *HandlerContext {
	ctx := &HandlerContext{
		"test",
		sessions.NewMemStore(time.Hour, time.Minute),
		users.NewMockStore(),
	}
	return ctx
}
func makeNewUser() *users.User {
	newUser := &users.NewUser{
		Email:        "test@test.com",
		Password:     "test1234",
		PasswordConf: "test1234",
		UserName:     "test",
		FirstName:    "juan",
		LastName:     "oh",
	}
	user, err := newUser.ToUser()
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return nil
	}
	return user
}

func TestUsersHandler(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		header             string
		user               *users.NewUser
		expectedStatusCode int
		fakeID             bool
	}{
		{
			"Not Valid Method",
			"PATCH",
			"",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "juan",
				LastName:     "oh",
			},
			http.StatusMethodNotAllowed,
			false,
		},
		{
			"Not Valid Header",
			"POST",
			"text/plain",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "juan",
				LastName:     "oh",
			},
			http.StatusUnsupportedMediaType,
			false,
		},
		{
			"Valid Method and Header",
			"POST",
			"application/json",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "juan",
				LastName:     "oh",
			},
			http.StatusCreated,
			false,
		},
		{
			"wrong new user",
			"POST",
			"application/json",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "4321",
				UserName:     "test",
				FirstName:    "juan",
				LastName:     "oh",
			},
			http.StatusBadRequest,
			false,
		},
		{
			"wrong signing key",
			"POST",
			"application/json",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "juan",
				LastName:     "oh",
			},
			http.StatusInternalServerError,
			true,
		},
		{
			"insertion error",
			"POST",
			"application/json",
			&users.NewUser{
				Email:        "juan3674@naver.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "juan",
				LastName:     "oh",
			},
			http.StatusBadRequest,
			true,
		},
	}
	for _, c := range cases {
		var nu bytes.Buffer
		jsonStr := c.user
		json.NewEncoder(&nu).Encode(jsonStr)
		req, _ := http.NewRequest(c.method, "/v1/user", bytes.NewBuffer(nu.Bytes()))
		req.Header.Set("Content-Type", c.header)
		respRec := httptest.NewRecorder()

		ctx := createFakeContext()
		if c.fakeID {
			ctx.SigningKey = ""
		}
		ctx.UsersHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestSessionHandler(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		header             string
		credential         *users.Credentials
		expectedStatusCode int
		fakeID             bool
	}{
		{
			"Not Valid Method",
			"GET",
			"application/json",
			&users.Credentials{
				Email:    "test@test.com",
				Password: "test1234",
			},
			http.StatusMethodNotAllowed,
			false,
		},
		{
			"Not Valid Content-Type",
			"POST",
			"text/plain",
			&users.Credentials{
				Email:    "test@test.com",
				Password: "test1234",
			},
			http.StatusUnsupportedMediaType,
			false,
		},
		{
			"Invalid Email",
			"POST",
			"application/json",
			&users.Credentials{
				Email:    "1234@tgmail.com",
				Password: "test1234",
			},
			http.StatusUnauthorized,
			false,
		},
		{
			"Invalid Password",
			"POST",
			"application/json",
			&users.Credentials{
				Email:    "test@test.com",
				Password: "testing1234",
			},
			http.StatusUnauthorized,
			false,
		},
		{
			"Valid Request",
			"POST",
			"application/json",
			&users.Credentials{
				Email:    "test@test.com",
				Password: "test1234",
			},
			http.StatusCreated,
			false,
		},
		{
			"Wrong Session Signing Key",
			"POST",
			"application/json",
			&users.Credentials{
				Email:    "test@test.com",
				Password: "test1234",
			},
			http.StatusInternalServerError,
			true,
		},
	}
	for _, c := range cases {
		var nu bytes.Buffer
		jsonStr := c.credential
		json.NewEncoder(&nu).Encode(jsonStr)
		req, _ := http.NewRequest(c.method, "/v1/user", bytes.NewBuffer(nu.Bytes()))
		req.Header.Set("Content-Type", c.header)
		respRec := httptest.NewRecorder()

		ctx := createFakeContext()
		if c.fakeID {
			ctx.SigningKey = ""
		}
		ctx.SessionHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestSpecificSessionHandler(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		url                string
		authorization      bool
		expectedStatusCode int
	}{
		{
			"Not Valid Method",
			"GET",
			"/v1/user/mine",
			true,
			http.StatusMethodNotAllowed,
		},
		{
			"Not Valid URL",
			"DELETE",
			"/v1/user/me",
			true,
			http.StatusForbidden,
		},
		{
			"Valid Method and URL",
			"DELETE",
			"/v1/user/mine",
			true,
			http.StatusOK,
		},
		{
			"Invalid authorization key",
			"DELETE",
			"/v1/user/mine",
			false,
			http.StatusInternalServerError,
		},
	}
	for _, c := range cases {
		ctx := createFakeContext()
		req, _ := http.NewRequest(c.method, c.url, nil)
		if !c.authorization {
			req.Header.Set("Authorization", "whatever")
		} else {
			newSessionState := &SessionState{
				time.Now(),
				makeNewUser(),
			}
			id, err := sessions.NewSessionID(ctx.SigningKey)
			if err != nil {
				t.Fatal("unexpected error: ", err)
				return
			}
			err = ctx.Session.Save(id, newSessionState)
			if err != nil {
				t.Fatal("unexpected error: ", err)
				return
			}
			req.Header.Set("Authorization", "Bearer "+id.String())
		}
		respRec := httptest.NewRecorder()
		ctx.SpecificSessionHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}
	}
}
