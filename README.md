# System Information App

Aplicación en Go para obtener información detallada del sistema donde se ejecuta. Incluye:

- Nombre del equipo (hostname)
- Sistema operativo, versión y arquitectura
- Uptime (tiempo encendido) y fecha/hora de booteo
- Información de CPU (modelo, núcleos físicos/lógicos, frecuencia, caché)
- Memoria RAM (total, usada, libre)
- Discos duros (total, usado, libre, tipo de sistema de archivos)
- IPs activas de las interfaces de red
- Exportación de la información en formato legible, JSON o archivo TXT
- Cabecera con versión y autor

## Novedades recientes

- Opción `--json` para exportar la salida en formato JSON (guarda en `systeminfo.json`)
- Opción `--txt` para exportar la salida legible a un archivo `systeminfo.txt`
- Detección y listado de IPs activas en las interfaces de red
- Cabecera con nombre de la app, versión y autor

## Requisitos

- Go 1.18 o superior

## Instalación de dependencias

```
go mod tidy
```

## Uso rápido

### Ejecutar directamente (modo desarrollo)

```bash
go run main.go [--json] [--txt]
```

- Sin parámetros: muestra la información en consola en formato legible.
- `--json`: muestra la información en consola en formato JSON.
- `--txt`: exporta la información legible al archivo `systeminfo.txt`.

### Compilar para Windows y Linux

#### Para Windows (desde cualquier sistema):

```sh
GOOS=windows GOARCH=amd64 go build -o systeminfo.exe main.go
```

#### Para Linux (desde cualquier sistema):

```sh
GOOS=linux GOARCH=amd64 go build -o systeminfo main.go
```

> **En PowerShell de Windows:**

```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o systeminfo main.go
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o systeminfo.exe main.go
```

> **En CMD de Windows:**

```cmd
set GOOS=linux
set GOARCH=amd64
go build -o systeminfo main.go
```

---

## Instructivo para ejecutar en Linux

1. **Instala Go** (si no lo tienes):

   - En Ubuntu/Debian:
     ```bash
     sudo apt update && sudo apt install golang-go
     ```
   - O descarga desde [golang.org/dl](https://golang.org/dl/)

2. **Clona el repositorio o copia los archivos fuente**:

   ```bash
   git clone https://github.com/kriollo/SOInformationGO
   cd SOInformationGO
   ```

3. **Instala las dependencias:**

   ```bash
   go mod tidy
   ```

4. **Ejecuta la app:**

   ```bash
   go run main.go
   # o para salida JSON
   go run main.go --json
   # o para exportar a TXT
   go run main.go --txt
   ```

5. **(Opcional) Compila el ejecutable:**
   ```bash
   go build -o systeminfo main.go
   chmod +x  systeminfo
   ./systeminfo
   ```

---

El ejecutable generado (`systeminfo.exe` o `systeminfo`) puede copiarse y ejecutarse en cualquier máquina del sistema objetivo.

## Ejemplo de salida

```
=============================================================
  System Information App       |  Versión v1.0.3
  Desarrollado por jorge Jara  |  https://github.com/kriollo
=============================================================

Hostname: TU-HOSTNAME
OS: linux
Platform: Ubuntu 22.04
Arch: x86_64
Kernel: 6.5.0-27-generic
Uptime: 2d 3h 15m 42s
Boot Time: 2025-04-21 13:05:01

CPU Info:
  CPUs detectados: 8
  Modelo                    Cores    MHz      ID
  Intel(R) Core(TM)...      4        1992.00  0
  Intel(R) Core(TM)...      4        1992.00  1

Memory:
  Total: 15.57 GB
  Used: 7.23 GB
  Free: 8.34 GB

Disk(s):
  Mount      Total        Used         Free         FS
  /          100.00 GB    45.00 GB     55.00 GB    ext4
  /home      400.00 GB    120.00 GB    280.00 GB   ext4

IPs activas:
  Interfaz                 IP
  eth0                     192.168.1.100
  wlan0                    192.168.1.101
```

## Notas

- La información de usuarios activos no está disponible en gopsutil v4.
- El formato de tiempo y bytes es amigable para humanos.
- Puedes modificar el código para agregar más detalles si lo requieres.

---

**Autor:** Jorge Jara
