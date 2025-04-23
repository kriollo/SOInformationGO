# System Information App

Esta aplicación en Go permite obtener información detallada del sistema donde se ejecuta, incluyendo:
- Nombre del equipo (hostname)
- Sistema operativo y versión
- Arquitectura
- Uptime (tiempo encendido)
- Fecha/hora de booteo
- Información de CPU (modelo, núcleos, frecuencia, caché)
- Memoria RAM (total, usada, libre)
- Discos duros (total, usado, libre, tipo de sistema de archivos)

## Requisitos
- Go 1.18 o superior

## Instalación de dependencias

```
go mod tidy
```

## Uso

### Ejecutar directamente (modo desarrollo)

```
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

#### Para Linux (desde Windows o cualquier sistema):
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

El ejecutable generado (`systeminfo.exe` o `systeminfo`) puede copiarse y ejecutarse en cualquier máquina del sistema objetivo.

## Ejemplo de salida

```
Hostname: DESKTOP-B0EF0LU
OS: windows
Platform: Microsoft Windows 11 Pro 24H2
Arch: x86_64
Uptime: 1d 2h 29m 51s
Boot Time: 2025-04-22 09:54:04

CPU Info:
  Model: AMD Ryzen 7 4700U with Radeon Graphics         , Cores: 8, Mhz: 2000.00, Cache Size: 0 B
  Logical CPUs: 8

Memory:
  Total: 31.36 GB
  Used: 15.66 GB
  Free: 15.70 GB

Disk(s):
  Mount: C:, Total: 444.91 GB, Used: 299.27 GB, Free: 145.64 GB, FS: 
  Mount: E:, Total: 29.81 GB, Used: 25.67 GB, Free: 4.14 GB, FS: 
```

## Notas
- La información de usuarios activos no está disponible en gopsutil v4.
- El formato de tiempo y bytes es amigable para humanos.
- Puedes modificar el código para agregar más detalles si lo requieres.

---

**Autor:** Tu Nombre
