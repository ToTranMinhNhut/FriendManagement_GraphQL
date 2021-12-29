package controllers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestControllers_GetUsers(t *testing.T) {
	tcs := map[string]struct {
		input          string
		expResult      string
		expError       error
		mockErr        error
		mockUserEmails []string
	}{
		"success with an empty input": {
			mockUserEmails: []string{"john@example.com", "andy@example.com"},
			expResult:      `{"count":2,"success":true,"users":["john@example.com","andy@example.com"]}`,
		},
		"failed with an unknow format input": {
			input:    `aaa`,
			expError: errors.New(`{"message":"Body request invalid format","success":false}`),
			mockErr:  errors.New(`{"message":"Body request invalid format","success":false}`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/v1/users", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)

			var mockService SpecService
			mockService.ExpectedCalls = []*mock.Call{
				mockService.On("GetUsers", mock.Anything).Return(tc.mockUserEmails, tc.mockErr),
			}
			friendController := NewFriendController(mockService)
			handler := http.HandlerFunc(friendController.GetUsers)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.expError != nil {
				require.EqualError(t, tc.expError, rr.Body.String())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, rr.Body.String())
			}
		})
	}
}

func TestControllers_CreateFriends(t *testing.T) {
	tcs := map[string]struct {
		input     string
		expResult string
		expError  error
		mockErr   error
	}{
		"success with an input": {
			input:     `{ "friends": ["andy@example.com","john@example.com"]}`,
			expResult: `{"success":true}`,
		},
		"failed with an unknow format input": {
			input:    `{}`,
			expError: errors.New(`{"message":"Request body is empty","success":false}`),
			mockErr:  errors.New(`{"message":"Request body is empty","success":false}`),
		},
		"failed with an input validation failure (two emails are similar)": {
			input:    `{ "friends": ["andy@example.com","andy@example.com"]}`,
			expError: errors.New(`{"message":"Two email addresses must be different","success":false}`),
			mockErr:  errors.New(`{"message":"Two email addresses must be different","success":false}`),
		},
		"failed with an input validation failure (email invalid format)": {
			input:    `{ "friends": ["andy@examplecom","andy@example.com"]}`,
			expError: errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
			mockErr:  errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
		},
		"failed with an input validation failure (number of emails are wrong)": {
			input:    `{ "friends": ["andy@examplecom"]}`,
			expError: errors.New(`{"message":"Number of email addresses must be 2","success":false}`),
			mockErr:  errors.New(`{"message":"Number of email addresses must be 2","success":false}`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/v1/friends", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)

			var mockService SpecService
			mockService.ExpectedCalls = []*mock.Call{
				mockService.On("CreateFriend", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockErr),
			}
			friendController := NewFriendController(mockService)
			handler := http.HandlerFunc(friendController.CreateFriend)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.expError != nil {
				require.EqualError(t, tc.expError, rr.Body.String())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, rr.Body.String())
			}
		})
	}
}

func TestControllers_GetFriends(t *testing.T) {
	tcs := map[string]struct {
		input            string
		expResult        string
		expError         error
		mockFriendEmails []string
		mockErr          error
	}{
		"success with an input": {
			input:            `{"Email":"andy@example.com"}`,
			expResult:        `{"count":1,"friends":["john@example.com"],"success":true}`,
			mockFriendEmails: []string{"john@example.com"},
		},
		"failed with an unknow format input": {
			input:    `{}`,
			expError: errors.New(`{"message":"Request body is empty","success":false}`),
			mockErr:  errors.New(`{"message":"Request body is empty","success":false}`),
		},
		"failed with an input validation failure": {
			input:    `{"Email":"andy@examplecom"}`,
			expError: errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
			mockErr:  errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/v1/friends", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)

			var mockService SpecService
			mockService.ExpectedCalls = []*mock.Call{
				mockService.On("GetFriends", mock.Anything, mock.Anything).Return(tc.mockFriendEmails, tc.mockErr),
			}
			friendController := NewFriendController(mockService)
			handler := http.HandlerFunc(friendController.GetFriends)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.expError != nil {
				require.EqualError(t, tc.expError, rr.Body.String())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, rr.Body.String())
			}
		})
	}
}

func TestControllers_GetCommonFriends(t *testing.T) {
	tcs := map[string]struct {
		input             string
		expResult         string
		expError          error
		mockCommonFriends []string
		mockErr           error
	}{
		"success with an input": {
			input:             `{ "friends": ["andy@example.com","john@example.com"]}`,
			mockCommonFriends: []string{"common@example.com"},
			expResult:         `{"count":1,"friends":["common@example.com"],"success":true}`,
		},
		"failed with an unknow format input": {
			input:    `{}`,
			expError: errors.New(`{"message":"Request body is empty","success":false}`),
			mockErr:  errors.New(`{"message":"Request body is empty","success":false}`),
		},
		"failed with an input validation failure (two emails are similar)": {
			input:    `{ "friends": ["andy@example.com","andy@example.com"]}`,
			expError: errors.New(`{"message":"Two email addresses must be different","success":false}`),
			mockErr:  errors.New(`{"message":"Two email addresses must be different","success":false}`),
		},
		"failed with an input validation failure (email invalid format)": {
			input:    `{ "friends": ["andy@examplecom","andy@example.com"]}`,
			expError: errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
			mockErr:  errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
		},
		"failed with an input validation failure (number of emails are wrong)": {
			input:    `{ "friends": ["andy@examplecom"]}`,
			expError: errors.New(`{"message":"Number of email addresses must be 2","success":false}`),
			mockErr:  errors.New(`{"message":"Number of email addresses must be 2","success":false}`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/v1/commonFriends", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)

			var mockService SpecService
			mockService.ExpectedCalls = []*mock.Call{
				mockService.On("GetCommonFriends", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.mockCommonFriends, tc.mockErr),
			}
			friendController := NewFriendController(mockService)
			handler := http.HandlerFunc(friendController.GetCommonFriends)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.expError != nil {
				require.EqualError(t, tc.expError, rr.Body.String())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, rr.Body.String())
			}
		})
	}
}

func TestControllers_CreateSubcription(t *testing.T) {
	tcs := map[string]struct {
		input     string
		expResult string
		expError  error
		mockErr   error
	}{
		"success with an input": {
			input:     `{"requestor": "andy@example.com","target": "lisa@example.com"}`,
			expResult: `{"success":true}`,
		},
		"failed with an unknow format input": {
			input:    `{}`,
			expError: errors.New(`{"message":"Request body is empty","success":false}`),
			mockErr:  errors.New(`{"message":"Request body is empty","success":false}`),
		},
		"failed with an input validation failure (two emails are similar)": {
			input:    `{"requestor": "andy@example.com","target": "andy@example.com"}`,
			expError: errors.New(`{"message":"Two email addresses must be different","success":false}`),
			mockErr:  errors.New(`{"message":"Two email addresses must be different","success":false}`),
		},
		"failed with an input validation failure (email invalid format)": {
			input:    `{"requestor": "andy@examplecom","target": "andy@example.com"}`,
			expError: errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
			mockErr:  errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
		},
		"failed with an input validation failure (target invalid)": {
			input:    `{"requestor": "andy@example.com"}`,
			expError: errors.New(`{"message":"Target field invalid format","success":false}`),
			mockErr:  errors.New(`{"message":"Target field invalid format","success":false}`),
		},
		"failed with an input validation failure (requestor invalid)": {
			input:    `{"target": "andy@example.com"}`,
			expError: errors.New(`{"message":"Requestor field invalid format","success":false}`),
			mockErr:  errors.New(`{"message":"Requestor field invalid format","success":false}`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/v1/subscription", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)

			var mockService SpecService
			mockService.ExpectedCalls = []*mock.Call{
				mockService.On("CreateSubscription", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockErr),
			}
			friendController := NewFriendController(mockService)
			handler := http.HandlerFunc(friendController.CreateSubcription)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.expError != nil {
				require.EqualError(t, tc.expError, rr.Body.String())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, rr.Body.String())
			}
		})
	}
}

func TestControllers_CreateUserBlocks(t *testing.T) {
	tcs := map[string]struct {
		input     string
		expResult string
		expError  error
		mockErr   error
	}{
		"success with an input": {
			input:     `{"requestor": "andy@example.com","target": "lisa@example.com"}`,
			expResult: `{"success":true}`,
		},
		"failed with an unknow format input": {
			input:    `{}`,
			expError: errors.New(`{"message":"Request body is empty","success":false}`),
			mockErr:  errors.New(`{"message":"Request body is empty","success":false}`),
		},
		"failed with an input validation failure (two emails are similar)": {
			input:    `{"requestor": "andy@example.com","target": "andy@example.com"}`,
			expError: errors.New(`{"message":"Two email addresses must be different","success":false}`),
			mockErr:  errors.New(`{"message":"Two email addresses must be different","success":false}`),
		},
		"failed with an input validation failure (email invalid format)": {
			input:    `{"requestor": "andy@examplecom","target": "andy@example.com"}`,
			expError: errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
			mockErr:  errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
		},
		"failed with an input validation failure (target invalid)": {
			input:    `{"requestor": "andy@example.com"}`,
			expError: errors.New(`{"message":"Target field invalid format","success":false}`),
			mockErr:  errors.New(`{"message":"Target field invalid format","success":false}`),
		},
		"failed with an input validation failure (requestor invalid)": {
			input:    `{"target": "andy@example.com"}`,
			expError: errors.New(`{"message":"Requestor field invalid format","success":false}`),
			mockErr:  errors.New(`{"message":"Requestor field invalid format","success":false}`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/v1/blocking", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)

			var mockService SpecService
			mockService.ExpectedCalls = []*mock.Call{
				mockService.On("CreateUserBlock", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockErr),
			}
			friendController := NewFriendController(mockService)
			handler := http.HandlerFunc(friendController.CreateUserBlock)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.expError != nil {
				require.EqualError(t, tc.expError, rr.Body.String())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, rr.Body.String())
			}
		})
	}
}

func TestControllers_GetRecipientEmails(t *testing.T) {
	tcs := map[string]struct {
		input               string
		text                string
		expResult           string
		expError            error
		mockRecipientEmails []string
		mockErr             error
	}{
		"success with an input": {
			input:               `{"sender": "andy@example.com","text": "Hello World! kate@example.com"}`,
			text:                "Hello World! kate@example.com",
			mockRecipientEmails: []string{"lisa@example.com", "kate@example.com"},
			expResult:           `{"recipients":["lisa@example.com","kate@example.com"],"success":true}`,
		},
		"failed with an unknow format input": {
			input:    `{}`,
			expError: errors.New(`{"message":"Request body is empty","success":false}`),
			mockErr:  errors.New(`{"message":"Request body is empty","success":false}`),
		},
		"failed with an input validation failure (email invalid format)": {
			input:    `{"sender": "andy@examplecom","text": "Hello World! kate@example.com"}`,
			expError: errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
			mockErr:  errors.New(`{"message":"andy@examplecom invalid format (ex: \"andy@example.com\")","success":false}`),
		},

		"failed with an input validation failure (target invalid)": {
			input:    `{"sender": "andy@example.com"}`,
			expError: errors.New(`{"message":"Text field invalid format","success":false}`),
			mockErr:  errors.New(`{"message":"Text field invalid format","success":false}`),
		},
		"failed with an input validation failure (requestor invalid)": {
			input:    `{"text": "Hello World! kate@example.com"}`,
			expError: errors.New(`{"message":"Sender field invalid format","success":false}`),
			mockErr:  errors.New(`{"message":"Sender field invalid format","success":false}`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/v1/recipients", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)

			var mockService SpecService
			mockService.ExpectedCalls = []*mock.Call{
				mockService.On("GetRecipientEmails", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.mockRecipientEmails, tc.mockErr),
			}
			friendController := NewFriendController(mockService)
			handler := http.HandlerFunc(friendController.GetRecipientEmails)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.expError != nil {
				require.EqualError(t, tc.expError, rr.Body.String())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, rr.Body.String())
			}
		})
	}
}
