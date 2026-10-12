package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(newServer("test", prometheus.NewRegistry()))
	t.Cleanup(srv.Close)
	return srv
}

func get(t *testing.T, url string) (int, string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func TestInfo(t *testing.T) {
	srv := newTestServer(t)
	code, body := get(t, srv.URL+"/")
	if code != http.StatusOK || !strings.Contains(body, `"servicio":"svc-catalogo"`) || !strings.Contains(body, `"entorno":"test"`) {
		t.Fatalf("respuesta inesperada: %d %s", code, body)
	}
}

func TestListProductos(t *testing.T) {
	srv := newTestServer(t)
	code, body := get(t, srv.URL+"/api/productos")
	var got []Producto
	if err := json.Unmarshal([]byte(body), &got); err != nil || code != http.StatusOK {
		t.Fatalf("respuesta inesperada: %d %s (%v)", code, body, err)
	}
	if len(got) != len(productos) {
		t.Fatalf("se esperaban %d productos, hay %d", len(productos), len(got))
	}
}

func TestListProductosPorCategoria(t *testing.T) {
	srv := newTestServer(t)
	_, body := get(t, srv.URL+"/api/productos?categoria=hogar")
	var got []Producto
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("se esperaban 2 productos de hogar, hay %d", len(got))
	}
}

func TestGetProducto(t *testing.T) {
	srv := newTestServer(t)
	if code, body := get(t, srv.URL+"/api/productos/P-003"); code != http.StatusOK || !strings.Contains(body, "Auriculares") {
		t.Fatalf("respuesta inesperada: %d %s", code, body)
	}
	if code, _ := get(t, srv.URL+"/api/productos/NO-EXISTE"); code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, llegó %d", code)
	}
}

func TestSondas(t *testing.T) {
	srv := newTestServer(t)
	for _, path := range []string{"/healthz", "/readyz"} {
		if code, _ := get(t, srv.URL+path); code != http.StatusOK {
			t.Fatalf("%s devolvió %d", path, code)
		}
	}
}

func TestMetricas(t *testing.T) {
	srv := newTestServer(t)
	get(t, srv.URL+"/api/productos")
	get(t, srv.URL+"/api/productos/NO-EXISTE")
	_, body := get(t, srv.URL+"/metrics")
	for _, want := range []string{
		`http_requests_total{code="200",handler="productos_list",method="get"} 1`,
		`http_requests_total{code="404",handler="productos_get",method="get"} 1`,
		`app_build_info{environment="test",service="svc-catalogo",version="dev"} 1`,
		`http_request_duration_seconds_bucket`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("falta en /metrics: %s", want)
		}
	}
}
