# svc-catalogo

Catálogo de productos de la tienda de demostración: servicio de ejemplo de la
[IDP del TFG](https://github.com/pablomartinvelasco/idp-tfg).

- API: `GET /api/productos`, `GET /api/productos/{id}`; métricas en `/metrics`.
- Imagen: `ghcr.io/pablomartinvelasco/svc-catalogo:<sha>` (CI en cada push a `main`).
- Despliegue: [`idp-tfg/gitops/apps/svc-catalogo`](https://github.com/pablomartinvelasco/idp-tfg/tree/main/gitops/apps/svc-catalogo).
- Documentación: [`docs/index.md`](docs/index.md) (publicada en el portal con TechDocs).
