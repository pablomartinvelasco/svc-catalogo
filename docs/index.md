# svc-catalogo

API REST de solo lectura con el catálogo de productos de una tienda ficticia. Es uno de los dos
servicios de ejemplo de la Internal Developer Platform del TFG y sigue la estructura estándar que
genera la plantilla de autoservicio.

## API

| Método y ruta | Descripción |
|---|---|
| `GET /` | Nombre del servicio, versión (SHA del commit) y entorno |
| `GET /api/productos` | Lista de productos; `?categoria=hogar` filtra por categoría |
| `GET /api/productos/{id}` | Un producto (`404` si no existe) |
| `GET /healthz`, `GET /readyz` | Sondas de *liveness* y *readiness* |
| `GET /metrics` | Métricas en formato Prometheus |

Los datos son sintéticos y viven en memoria.

## Métricas

| Métrica | Tipo | Etiquetas |
|---|---|---|
| `http_requests_total` | contador | `handler`, `method`, `code` |
| `http_request_duration_seconds` | histograma | `handler`, `method`, `code` |
| `app_build_info` | gauge (1) | `service`, `version`, `environment` |

Además, las métricas estándar del runtime de Go y del proceso.

## Configuración

| Variable | Por defecto | Uso |
|---|---|---|
| `PORT` | `8080` | Puerto HTTP |
| `APP_ENV` | `local` | Entorno que informa el servicio (`dev`, `staging`) |

## Despliegue

- Cada push a `main` ejecuta las pruebas y publica `ghcr.io/pablomartinvelasco/svc-catalogo:<sha>`.
- La configuración de despliegue está en
  [`idp-tfg/gitops/apps/svc-catalogo`](https://github.com/pablomartinvelasco/idp-tfg/tree/main/gitops/apps/svc-catalogo):
  Argo CD la sincroniza en `apps-dev` y `apps-staging`.
- Promoción: PR al repositorio de plataforma que cambia la etiqueta de imagen del overlay del entorno.

## Desarrollo local

```bash
go test ./...
go run .                      # http://localhost:8080
docker build -t svc-catalogo .
```
