package api

import (
	"encoding/json"
	"io"
	"mod36a41/pkg/storage"
	"mod36a41/pkg/storage/memdb"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// news to display
const wantLen = 4

func TestAPI_lastNewsHandler(t *testing.T) {
	dbase := memdb.New()
	api := New(dbase)

	req := httptest.NewRequest(http.MethodGet, "/news/"+strconv.Itoa(wantLen), nil)
	rr := httptest.NewRecorder()

	api.router.ServeHTTP(rr, req)
	// reads response code
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Error: got %d, need %d", rr.Code, http.StatusOK)
	}
	// reads server response
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Wasn't able to decode server response: %v", err)
	}
	// decodes JSON
	var news []storage.Post
	err = json.Unmarshal(b, &news)
	if err != nil {
		t.Fatalf("Wasn't able to decode server response: %v", err)
	}
	// checks the number of posts aquired
	if len(news) != wantLen {
		t.Fatalf("Error: got %d posts, want %d", len(news), wantLen)
	}
}
