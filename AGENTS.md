# Reglas del proyecto

## 1. Carpetas fuera de alcance

El agente **nunca debe leer, analizar, modificar, crear archivos dentro de, ni ejecutar archivos de** las siguientes carpetas:

* `crud-collections/`
* `doc/`

Estas carpetas están completamente fuera del alcance del agente.

No deben utilizarse como fuente de contexto para implementar funcionalidades.

---

## 2. Tests obligatorios

Los tests del proyecto se encuentran en:

```text
test/
```

**Siempre debes ejecutar los tests antes de considerar una implementación como terminada.**

No avances a la siguiente fase de implementación si los tests existentes fallan.

La validación mínima después de realizar cambios debe incluir:

1. Ejecutar todos los tests de `test/`.
2. Verificar que los tests existentes continúan pasando.
3. Verificar que la nueva implementación no rompe funcionalidades existentes.

---

## 3. Toda nueva implementación requiere tests

Cada vez que implementes una nueva funcionalidad, comportamiento, endpoint, servicio o modificación relevante:

* Debes crear los tests correspondientes.
* Los tests **siempre deben estar dentro de `test/`**.
* No debes colocar tests fuera de `test/` salvo que una herramienta del lenguaje lo requiera explícitamente.
* Los tests deben cubrir el comportamiento esperado de la nueva implementación.

### Flujo obligatorio

Para cualquier nueva implementación:

```text
Implementar
    ↓
Crear tests
    ↓
Ejecutar tests
    ↓
Corregir errores
    ↓
Ejecutar tests nuevamente
    ↓
Continuar con la siguiente fase
```

Una implementación no se considera terminada hasta que sus tests hayan sido ejecutados correctamente.

---

## 4. No modificar tests para ocultar errores

Si una implementación hace que un test falle:

* Primero revisa la implementación.
* Corrige el código de producción cuando corresponda.
* No modifiques ni elimines un test únicamente para conseguir que pase.
* Si el comportamiento esperado cambió intencionalmente, actualiza el test de forma coherente con el nuevo requisito.

---

## 5. Regla general

Antes de realizar cambios importantes, analiza primero la estructura existente del proyecto y respeta las convenciones actuales.

Prioriza cambios pequeños, verificables y compatibles con la arquitectura existente.
