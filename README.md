# Panadería Concurrente

Una simulación visual interactiva de una línea de producción de pastelería, construida en Go con [Ebiten](https://ebiten.org/). Este proyecto demuestra patrones de concurrencia fundamentales de Go, como **Pipeline**, **Worker Pool** y **Bounded Resource** en un entorno gráfico en tiempo real.


## Descripción General

Este programa simula las etapas de producción de una pastelería como un **pipeline concurrente**. Cada etapa de la producción (recepción, mezclado, horneado, decoración, empaquetado) se ejecuta en su propia goroutine, pasando los pasteles a la siguiente etapa a través de canales (channels).

El usuario puede hacer clic en el botón "MAKE ORDER" para inyectar nuevos pedidos en el inicio del pipeline. La simulación visualiza cada pastel (representado por un `*models.Cake`) moviéndose a través de las diferentes etapas de trabajo.

## Cómo Funciona: El Flujo del Pipeline

El corazón de la simulación se encuentra en `workers/pipeline.go`, que configura y conecta todas las etapas:

1.  **Recepción (`Receptionist`):** Cuando haces clic en "MAKE ORDER", se envía un ID de pedido al canal `OrdersChan`. La goroutine `Receptionist` lo recibe, simula un tiempo de procesamiento, crea un objeto `Cake` y lo envía al siguiente canal (`pending`).
2.  **Mezclado (`MixCoordinator`):** Esta etapa es un **Pool de Trabajadores**. Múltiples goroutines `MixerWorker` (5 en este caso) leen concurrentemente del canal `pending`. Esto permite procesar hasta 5 pasteles en paralelo. Una vez mezclados, los envían al canal `mixed`.
3.  **Horneado (`Bake`):** Esta etapa lee del canal `mixed`. El horno (`models.Oven`) es un **Recurso Limitado** con una capacidad fija (3). La goroutine `Bake` debe "adquirir" un lugar en el horno (usando un semáforo) antes de poder procesar el pastel. Si el horno está lleno, la goroutine espera. Una vez horneado, el pastel se envía al canal `baked`.
4.  **Decoración (`Decorator`):** Una sola goroutine que lee de `baked`, simula el trabajo y envía al canal `decorated`.
5.  **Empaquetado (`Packager`):** Una sola goroutine que lee de `decorated`, simula el trabajo y envía al canal `ready`.
6.  **Finalizado (`Finisher`):** La última etapa. Lee de `ready`, anima un "fade out" del pastel, y actualiza el contador de pasteles completados.

---

## Patrones de Concurrencia Utilizados

Este proyecto es una demostración práctica de varios patrones de concurrencia clave:

### 1. Pipeline (Tuberías)

El patrón arquitectónico principal. Toda la panadería es un pipeline de varias etapas.

* **Cómo funciona:** Cada etapa del proceso es una goroutine (o un grupo de ellas) que:
    1.  Recibe trabajo de un canal de entrada (`<-in`).
    2.  Procesa el trabajo (p.ej., `time.Sleep(...)`).
    3.  Envía el resultado a un canal de salida (`out <- ...`).
* **Código:** `workers/pipeline.go` define una serie de canales (`pending`, `mixed`, `baked`, etc.) que conectan las goroutines de cada etapa (`Receptionist`, `MixCoordinator`, `Bake`, ...).
* **Beneficio:** Desacopla la lógica de cada etapa y permite que todas las etapas se ejecuten concurrentemente. Mientras un pastel se está horneando, otro puede estar mezclándose y un tercero decorándose.

### 2. Worker Pool (Pool de Trabajadores)

La etapa de mezclado (`MixCoordinator`) implementa este patrón para paralelizar una tarea específica.

* **Cómo funciona:** En lugar de una sola goroutine leyendo del canal `pending`, el `MixCoordinator` inicia un número fijo de goroutines `MixerWorker` (p.ej., 5).
    * Todas estas goroutines leen del *mismo* canal de entrada (`pending`). Go se encarga de distribuir los pasteles entre los trabajadores disponibles.
    * Todas las goroutines escriben en el *mismo* canal de salida (`mixed`).
* **Código:** `workers/mixer.go`. El `MixCoordinator` usa un `sync.WaitGroup` para saber cuándo todos los `MixerWorker` han terminado (después de que el canal `pending` se cierra) antes de cerrar el canal `mixed`.
* **Beneficio:** Aumenta significativamente el *throughput* (rendimiento) de una etapa que de otra manera sería un cuello de botella.

### 3. Bounded Resource (Recurso Limitado)

El horno (`models.Oven`) es un recurso con una capacidad limitada (3 pasteles).

* **Cómo funciona:** Se utiliza una combinación de `sync.Mutex` y `sync.Cond` (Variable de Condición) para gestionar el acceso.
    * `Oven.Use()`: La goroutine `Bake` llama a esta función. La función bloquea el mutex y comprueba si `InUse < Capacity`.
        * Si hay espacio, incrementa `InUse` y continúa.
        * Si está lleno (`InUse >= Capacity`), llama a `cond.Wait()`. Esto **atómica y eficientemente** desbloquea el mutex y pone a la goroutine "en espera".
    * `Oven.Release()`: Cuando un pastel termina de hornearse, se llama a esta función. Bloquea el mutex, decrementa `InUse` y llama a `cond.Signal()`. Esto "despierta" a *una* de las goroutines que estaba esperando en `cond.Wait()`.
* **Código:** `models/oven.go` y `workers/oven.go` (la función `Bake`).
* **Beneficio:** Permite limitar la concurrencia en una sección específica de manera segura y eficiente, sin usar *busy-waiting* (bucles que consumen CPU).

### 4. Paso de Mensajes (para Sincronización de UI)

Este es un patrón crucial para integrar goroutines de fondo con un hilo principal (el bucle de Ebiten).

* **Cómo funciona:** Las goroutines del pipeline (Mezclador, Hornero, etc.) **nunca modifican directamente** la lista de pasteles del juego (`g.State.Cakes`). Hacerlo causaría *data races* (condiciones de carrera), ya que el hilo principal (`game.Update`) también está leyendo esa lista para dibujarla.
* En su lugar, cuando un trabajador necesita actualizar el estado de un pastel (p.ej., su posición o estado), envía un *mensaje* (el propio objeto `*models.Cake` actualizado) al canal `g.State.UpdatesChan`.
* El hilo principal (`game.Update`) es el **único** que lee de `UpdatesChan` (usando un `select` no bloqueante). Al ser el único "escritor" de la lista `g.State.Cakes`, se evita cualquier *race condition* sobre esa lista.
* **Código:** `game/game.go` (en `Update()`) y todas las funciones de `workers/` (especialmente `MoveAlongBelt`).
* **Beneficio:** Sigue la filosofía de Go: "No te comuniques compartiendo memoria; comparte memoria comunicándote".


## Cómo Ejecutarlo

1.  Asegúrate de tener [Go instalado](https://go.dev/doc/install).
2.  Clona el repositorio (o ten los archivos en un directorio).
3.  Abre una terminal en el directorio raíz del proyecto.
4.  Ejecuta:

    ```bash
    go run main.go -race
    ```

    (Go descargará e instalará automáticamente Ebiten y otras dependencias).

5.  Haz clic en el botón "MAKE ORDER" para empezar a hornear.