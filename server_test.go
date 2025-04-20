package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}
func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}
func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	tests := []struct {
		name               string
		player             string
		expectedHTTPStatus int
		expectedScore      string
	}{
		{
			name:               "Returns Pepper score",
			player:             "Pepper",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "20",
		},
		{
			name:               "Returns Floyd score",
			player:             "Floyd",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "10",
		},
		{
			name:               "Returns Pepper score",
			player:             "Apollo",
			expectedHTTPStatus: http.StatusNotFound,
			expectedScore:      "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := newGetScoreRequest(tt.player)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertBodyResponse(t, response.Body.String(), tt.expectedScore)
		})
	}
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)
	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns 200 on /league", func(t *testing.T) {
		store := StubPlayerStore{}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
	})
	t.Run("it the league table as json", func(t *testing.T) {
		wantedLeague := []Player{
			{"chleo", 32},
			{"chris", 20},
			{"tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
	})

}

func newGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}
func newPostWinRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return request
}
func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}
func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}
