package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"task111-shortlink/internal/store"
)

func TestActivityWindowOnlyIncludesRequestedCode(t *testing.T) {
	h := newServer(t)
	create := func(target string) store.Link {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(`{"target_url":"`+target+`"}`)))
		if rec.Code != http.StatusCreated {
			t.Fatalf("create status = %d", rec.Code)
		}
		var l store.Link
		if err := json.Unmarshal(rec.Body.Bytes(), &l); err != nil {
			t.Fatal(err)
		}
		return l
	}
	a, b := create("https://example.com/a"), create("https://example.com/b")
	for _, code := range []string{a.Code, b.Code} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/links/"+code+"/click", nil))
		if rec.Code != http.StatusCreated {
			t.Fatalf("click status = %d", rec.Code)
		}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/links/"+a.Code+"/activity", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("activity status = %d: %s", rec.Code, rec.Body.String())
	}
	var view struct{ Window struct{ Count int `json:"count"` } `json:"recent_window"` }
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if view.Window.Count != 1 {
		t.Fatalf("activity window count = %d, want 1", view.Window.Count)
	}
}
