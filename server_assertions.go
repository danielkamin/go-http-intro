package main

import (
	"net/http/httptest"
	"reflect"
	"testing"
)

func assertBodyResponse(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response is wrong, got %q, want %q", got, want)
	}
}
func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}
func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertContentType(t testing.TB, got *httptest.ResponseRecorder, want string) {
	t.Helper()
	if got.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of application/json, got %v", got.Result().Header.Get("content-type"))
	}
}
