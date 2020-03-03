package project_package

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAllMessages(t *testing.T) {
	session := PrepareTestKesyspace()
	AddTestingData()
	defer session.Close()
	req, err := http.NewRequest("GET", "/api/messages", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllMessages)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	responseBody := rr.Body.String()
	responseBody = strings.Trim(responseBody, " \r\n")
	expected := fmt.Sprintf(`[{"id":1,"email":"jan.kowalski@example.com","title":"Interview","content":"simple text","magic_number":101,"created":%v},{"id":2,"email":"jan.kowalski@example.com","title":"Interview 2","content":"simple text 2","magic_number":22,"created":%v},{"id":3,"email":"anna.zajkowska@example.com","title":"Interview 3","content":"simple text 3","magic_number":101,"created":%v}]`, GetTime()-7, GetTime(), GetTime()-6)
	if responseBody != expected {
		t.Errorf("handler returned unexpected body:\n got %v want %v",
			responseBody, expected)
	}
	DropTestKeyspace()
}

func TestGetMessages(t *testing.T) {
	session := PrepareTestKesyspace()
	AddTestingData()
	defer session.Close()

	tt := []struct {
		routeVariable string
		shouldPass    bool
		expectedBody  string
	}{
		{
			"jan.kowalski@example.com",
			true,
			fmt.Sprintf(`[{"id":1,"email":"jan.kowalski@example.com","title":"Interview","content":"simple text","magic_number":101,"created":%v},{"id":2,"email":"jan.kowalski@example.com","title":"Interview 2","content":"simple text 2","magic_number":22,"created":%v}]`, GetTime()-7, GetTime()),
		},
		{
			"jan.kowalski@",
			false,
			"400 - invalid email address",
		},
		{
			"@example.com",
			false,
			"400 - invalid email address",
		},
		{
			"jan.kowalski.example.com",
			false,
			"400 - invalid email address",
		},
		{
			"no.such.person@example.com",
			true,
			"null",
		},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/api/messages/%s", tc.routeVariable)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		// Create a router to ba able to pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.HandleFunc("/api/messages/{emailValue}", GetMessages)
		router.ServeHTTP(rr, req)

		// Check status code
		if tc.shouldPass {
			if rr.Code != http.StatusOK {
				t.Errorf("handler returned wrong status code for url /api/messages/%s: got %v want %v",
					tc.routeVariable, rr.Code, http.StatusOK)
			}
		} else if !tc.shouldPass {
			if rr.Code != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code for url /api/messages/%s: got %v want %v",
					tc.routeVariable, rr.Code, http.StatusBadRequest)
			}
		}
		// Check the response body is what we expect.
		responseBody := rr.Body.String()
		responseBody = strings.Trim(responseBody, " \r\n")
		expected := tc.expectedBody
		if responseBody != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
	DropTestKeyspace()
}

func TestCreateMessage(t *testing.T) {
	session := PrepareTestKesyspace()
	defer session.Close()

	tt := []struct {
		body       string
		shouldPass bool
	}{
		{
			`{
			"email": "ania.kwiatkowska@example.com",
			"title": "15",
			"content": "lorem ipsum dolores",
			"magic_number": 123
			}`,
			true,
		},
		{
			`{
			"email": "bad.mail@",
			"title": "15",
			"content": "lorem ipsum dolores",
			"magic_number": 123
			}`,
			false,
		},
		{
			`{
			"email": "ania.kwiatkowska@example.com",
			"title": 15,
			"content": "lorem ipsum dolores",
			"magic_number": 123
			}`,
			false,
		},
		{
			`{
			"email": "ania.kwiatkowska@example.com",
			"title": "15",
			"content": "lorem ipsum dolores",
			"magic_number": "should be integer"
			}`,
			false,
		},
		{
			`{
			"title": "15",
			"content": "lorem ipsum dolores",
			"magic_number": 123
			}`,
			false,
		},
	}
	var expectedBody string

	for _, tc := range tt {
		path := fmt.Sprintf("/api/messages")
		var jsonStr = []byte(tc.body)
		req, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// Create a router to ba able to pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.HandleFunc("/api/messages", CreateMessage)
		router.ServeHTTP(rr, req)

		// Check status code
		if tc.shouldPass {
			if rr.Code != http.StatusCreated {
				t.Errorf("handler returned wrong status code for data:\n%s\ngot %v want %v",
					tc.body, rr.Code, http.StatusCreated)
			}
			expectedBody = "201 - Message created!"
		} else if !tc.shouldPass {
			if rr.Code != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code for data:\n%s\ngot %v want %v",
					tc.body, rr.Code, http.StatusBadRequest)
			}
			expectedBody = "400 - bad data"
		}
		// Check the response body is what we expect.
		responseBody := rr.Body.String()
		responseBody = strings.Trim(responseBody, " \r\n")
		if responseBody != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedBody)
		}
	}
	DropTestKeyspace()
}

func TestSendMessage(t *testing.T) {
	session := PrepareTestKesyspace()
	AddTestingData()
	defer session.Close()

	tt := []struct {
		body           string
		shouldPass     bool
		expectedStatus int
	}{
		{
			`{"magic_number": 999}`,
			true,
			204,
		},
		{
			`{"magic_number": 101}`,
			true,
			202,
		},
		{
			`{"magic_number": "123"}`,
			false,
			400,
		},
		{
			`{"magic_number": "not an int"}`,
			false,
			400,
		},
		{
			`{"magic_number": true}`,
			false,
			400,
		},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/api/send")
		var jsonStr = []byte(tc.body)
		req, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// Create a router to ba able to pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.HandleFunc("/api/send", SendMessage)
		router.ServeHTTP(rr, req)

		// Check status code
		if tc.shouldPass {
			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code for data:\n%s\ngot %v want %v",
					tc.body, rr.Code, tc.expectedStatus)
			}
		} else if !tc.shouldPass {
			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code for data:\n%s\ngot %v want %v",
					tc.body, rr.Code, tc.expectedStatus)
			}
		}
		// Check the response body is what we expect.
		responseBody := rr.Body.String()
		responseBody = strings.Trim(responseBody, " \r\n")
		var expectedBody string
		if tc.expectedStatus == 400 {
			expectedBody = "400 - bad data"
		} else if tc.expectedStatus == 202 {
			expectedBody = "202 - Message(s) sent!"
		} else if tc.expectedStatus == 204 {
			expectedBody = ""
		}
		if responseBody != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedBody)
		}
	}
	DropTestKeyspace()
}
