---
title: Referencia de errores
---

# Referencia de errores

Cuando una petición falla, `alegra` imprime `Error: alegra: HTTP <code>: <message>` y
sale con un código distinto de cero. Aquí tienes qué significa cada estado y qué hacer.

## Códigos de estado HTTP

| Código | Significado | Qué hacer |
|------|---------|------------|
| **400** | Petición incorrecta — el cuerpo está mal formado o le faltan campos requeridos. | Lee el mensaje; normalmente nombra el campo. Para documentos, revisa los requisitos del país (resolución, tipo de identificación, régimen). Vuelve a ejecutar con `--dry-run` para inspeccionar lo que enviaste. |
| **401** | Falló la autenticación. | Email/token incorrecto. Vuelve a ejecutar `alegra auth login`, o revisa `ALEGRA_EMAIL`/`ALEGRA_TOKEN`. Verifica con `alegra auth status`. |
| **402** | Pago requerido — tu **plan no incluye** esta acción, o la cuenta está suspendida. | Mejora el plan de Alegra, o evita el endpoint. Común en algunos reportes y en funciones de complementos/marketplace. |
| **403** | Prohibido — el usuario no tiene permiso para esta acción. | Revisa el rol/permisos del usuario en Alegra (`alegra users self`). |
| **404** | No encontrado — el id no existe (o la cuenta está suspendida). | Verifica el id con un `list`. |
| **405** | Método no permitido para este endpoint. | Suele ser un bug de uso/CLI — p. ej. una acción que el recurso no soporta. |
| **429** | Demasiadas peticiones — alcanzaste el tope por minuto de tu plan. | La CLI reintenta automáticamente con backoff y se adapta a los headers `X-Rate-Limit-*`. Si estás haciendo scripting, baja el ritmo; prefiere `--count` sobre `--all`. |
| **500** | Error del servidor de Alegra. | Transitorio; la CLI reintenta. Si persiste, prueba más tarde. |
| **503** | Servicio no disponible (mantenimiento). | Reintenta más tarde. |

La CLI reintenta automáticamente `429` y `5xx` con backoff exponencial; los `4xx`
(excepto `429`) se devuelven de inmediato, ya que reintentar no servirá de nada.

## Límites de tasa

Cada respuesta incluye:

| Header | Significado |
|--------|---------|
| `X-Rate-Limit-Limit` | El tope por minuto de tu plan. |
| `X-Rate-Limit-Remaining` | Llamadas restantes en el minuto actual. |
| `X-Rate-Limit-Reset` | Segundos hasta que la ventana se reinicia. |

Ejecuta con `-v/--verbose` para ver el logging a nivel de petición.

## Fallos al timbrar (emisión electrónica)

Cuando la autoridad fiscal rechaza un documento, el error suele traer un código del
proveedor. Las familias más comunes (Colombia / DIAN vía la capa del e-Provider):

| Código(s) | Significado | Solución |
|---------|---------|-----|
| `AEP1001` / `AEP1005` / `AEP1003` | Certificado inválido / expirado / faltante | Renueva o sube el certificado digital en la configuración de Alegra. |
| `AEP6008` | No autorizado para emitir facturas electrónicas | Completa la habilitación ante la DIAN. |
| `AEP6009` / `AEP6011` | Numeración rechazada (duplicada / fuera de rango — "Regla 90") | Usa un rango de resolución válido; no vuelvas a emitir un número. |
| `AEP6012`–`AEP6017` | Falta un campo requerido (tipo de organización, tipo de id, nombre, régimen, email, dirección) | Agrega el campo al payload/contacto. |
| `AEP2011` | Ya enviado / en proceso | **Revisa el estado antes de reintentar** — puede que ya se haya emitido. |
| `EPR500` / `EPR503` | Error de comunicación con la DIAN / timeout | Transitorio — verifica el estado y luego reintenta. |

> La emisión **no es idempotente** del lado del servidor. Emite con `alegra invoices emit`
> — mantiene una guarda de idempotencia local y omite las facturas ya emitidas a menos que
> pases `--force`. Después de cualquier error al timbrar, ejecuta `alegra invoices get <id>`
> para ver si en realidad se emitió **antes** de reintentar. Mira
> [Facturación electrónica](../guides/electronic-invoicing.md).

## Cómo obtener más detalle

- `--dry-run` muestra la petición exacta que enviaría la CLI.
- `-o json` en un `get` muestra el objeto completo del servidor (incluidos los detalles de
  `stamp`/error que Alegra haya adjuntado).
- `-v` habilita el logging de depuración (peticiones, reintentos, presupuesto de límite de
  tasa).
