package envs

import (
	"net/http"
	"strings"
	"testing"
)

func Test_getCookieValueReturnNilOnEmptirCookies(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)

	cookieValue := getCookieValue(request, "some")

	if cookieValue != nil {
		t.Errorf("Cookie value expected to be nil, got %v", cookieValue)
	}
}

func Test_getCookieValueReturnNilOnMissingCookie(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	cookie := &http.Cookie{Name: "other", Value: "value"}
	request.AddCookie(cookie)

	cookieValue := getCookieValue(request, "some")

	if cookieValue != nil {
		t.Errorf("Cookie value expected to be nil, got %v", cookieValue)
	}
}

func Test_getCookieValueReturnValueIfCookiePresent(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	cookie := &http.Cookie{Name: "some", Value: "value"}
	request.AddCookie(cookie)

	cookieValue := getCookieValue(request, "some")

	if *cookieValue != cookie.Value {
		t.Errorf("Cookie value expected to be %v, got %v", cookie.Value, cookieValue)
	}
}

type testWriter struct {
	TestHeader http.Header
}

func (responseWriter testWriter) Header() http.Header {
	return responseWriter.TestHeader
}
func (responseWriter testWriter) Write([]byte) (int, error) {
	panic("NOIMPL")
}
func (responseWriter testWriter) WriteHeader(statusCode int) {
	panic("NOIMPL")
}

func Test_setCookieValue(t *testing.T) {
	responseWriter := testWriter{
		TestHeader: http.Header{},
	}
	cookieValue := "value"
	setCookieValue(responseWriter, "some", &cookieValue)

	cookieHeader := responseWriter.TestHeader.Get("Set-Cookie")
	if strings.Contains(cookieHeader, "some=value;") == false {
		t.Errorf("Expected to have cookie some=value; , but got %v", cookieHeader)
	}
}

func TestHandlerFindsSelectedEnvironment(t *testing.T) {
	handler := Handler{Cookie: "some"}
	request, _ := http.NewRequest("GET", "/", nil)
	cookie := &http.Cookie{Name: "some", Value: "value"}
	request.AddCookie(cookie)

	cookieValue := handler.Selected(request)

	if *cookieValue != "value" {
		t.Errorf("Handler should select environmetn from cookie")
	}
}
