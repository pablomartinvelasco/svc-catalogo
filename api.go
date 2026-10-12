package main

import "net/http"

// Producto del catálogo. Los datos son sintéticos (RNF-03).
type Producto struct {
	ID        string  `json:"id"`
	Nombre    string  `json:"nombre"`
	Categoria string  `json:"categoria"`
	Precio    float64 `json:"precio"`
	Stock     int     `json:"stock"`
}

var productos = []Producto{
	{ID: "P-001", Nombre: "Taza de cerámica", Categoria: "hogar", Precio: 8.50, Stock: 120},
	{ID: "P-002", Nombre: "Cuaderno A5", Categoria: "papelería", Precio: 3.20, Stock: 340},
	{ID: "P-003", Nombre: "Auriculares inalámbricos", Categoria: "electrónica", Precio: 49.90, Stock: 35},
	{ID: "P-004", Nombre: "Mochila urbana", Categoria: "accesorios", Precio: 39.00, Stock: 60},
	{ID: "P-005", Nombre: "Lámpara de escritorio", Categoria: "hogar", Precio: 24.75, Stock: 0},
}

func registerAPI(mux *http.ServeMux, instrument func(string, http.HandlerFunc) http.Handler) {
	mux.Handle("GET /api/productos", instrument("productos_list", listProductos))
	mux.Handle("GET /api/productos/{id}", instrument("productos_get", getProducto))
}

// listProductos devuelve el catálogo; ?categoria= filtra por categoría.
func listProductos(w http.ResponseWriter, r *http.Request) {
	categoria := r.URL.Query().Get("categoria")
	out := make([]Producto, 0, len(productos))
	for _, p := range productos {
		if categoria == "" || p.Categoria == categoria {
			out = append(out, p)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func getProducto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for _, p := range productos {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}
	writeError(w, http.StatusNotFound, "producto no encontrado: "+id)
}
