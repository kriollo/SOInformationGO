package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"systemInformation/internal/info"
	"systemInformation/internal/model"
	"systemInformation/internal/utils"
)

const (
	AppName = "System Information App"
	Version = "v1.0.3"
)

func printHumanReadable(info *model.SystemInfo) {
	writeHumanReadable(os.Stdout, info)
}

func writeHumanReadable(w *os.File, info *model.SystemInfo) {
	fmt.Fprintf(w, "Hostname: %s\n", info.Hostname)
	fmt.Fprintf(w, "OS: %s\n", info.OS)
	fmt.Fprintf(w, "Platform: %s %s\n", info.Platform, info.PlatformVer)
	fmt.Fprintf(w, "Arch: %s\n", info.Arch)
	if info.KernelVersion != "" {
		fmt.Fprintf(w, "Kernel: %s\n", info.KernelVersion)
	}
	fmt.Fprintf(w, "Uptime: %s\n", utils.HumanDuration(info.Uptime))
	fmt.Fprintf(w, "Boot Time: %s\n", utils.HumanBootTime(info.BootTime))

	// CPU Info
	fmt.Fprintf(w, "\nCPU Info:\n")
	if len(info.CPU) > 0 {
		fmt.Fprintf(w, "  CPUs detectados: %d\n", len(info.CPU))
		fmt.Fprintf(w, "  %-25s %-8s %-8s %-8s\n", "Modelo", "Cores", "MHz", "ID")
		for i, c := range info.CPU {
			fmt.Fprintf(w, "  %-25s %-8d %-8.2f %-8d\n", c.ModelName, c.Cores, c.Mhz, c.CPU)
			if i == 0 && len(info.CPU) > 1 {
				fmt.Fprintf(w, "\n")
			}
		}
	}

	// Memory
	fmt.Fprintf(w, "\nMemory:\n")
	fmt.Fprintf(w, "  Total: %s\n", utils.HumanBytes(info.Memory.Total))
	fmt.Fprintf(w, "  Used: %s\n", utils.HumanBytes(info.Memory.Used))
	fmt.Fprintf(w, "  Free: %s\n", utils.HumanBytes(info.Memory.Available))

	// Disks
	fmt.Fprintf(w, "\nDisk(s):\n")
	if len(info.Disk) > 0 {
		fmt.Fprintf(w, "  %-10s %-12s %-12s %-12s %-6s\n", "Mount", "Total", "Used", "Free", "FS")
		for _, d := range info.Disk {
			fmt.Fprintf(w, "  %-10s %-12s %-12s %-12s %-6s\n", d.Path, utils.HumanBytes(d.Total), utils.HumanBytes(d.Used), utils.HumanBytes(d.Free), d.Fstype)
		}
	} else {
		fmt.Fprintf(w, "  No se detectaron discos.\n")
	}

	// IPs activas
	if len(info.IPs) > 1 {
		fmt.Fprintf(w, "\nIPs activas:\n")
		fmt.Fprintf(w, "  %-25s %-40s\n", "Interfaz", "IP")
		for _, nip := range info.IPs {
			fmt.Fprintf(w, "  %-25s %-40s\n", nip.Interface, nip.IP)
		}
	} else if len(info.IPs) == 1 {
		nip := info.IPs[0]
		fmt.Fprintf(w, "\nIPs activas:\n  %s: %s\n", nip.Interface, nip.IP)
	}
}

func main() {
	fmt.Printf("\n=============================================================\n")
	fmt.Printf("  %s       |  Versión %s\n", AppName, Version)
	fmt.Printf("  Desarrollado por jorge Jara  |  https://github.com/kriollo\n")
	fmt.Printf("=============================================================\n\n")

	jsonFlag := flag.Bool("json", false, "Exportar salida en formato JSON")
	txtFlag := flag.Bool("txt", false, "Exportar salida legible a systeminfo.txt")
	flag.Parse()

	info, err := info.GetSystemInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error obteniendo información del sistema: %v\n", err)
		os.Exit(1)
	}

	if *jsonFlag {
		output, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error serializando a JSON: %v\n", err)
			os.Exit(1)
		}
		f, err := os.Create("systeminfo.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creando archivo systeminfo.json: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		_, err = f.Write(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error escribiendo archivo systeminfo.json: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Información guardada en systeminfo.json")
	} else if *txtFlag {
		f, err := os.Create("systeminfo.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creando archivo systeminfo.txt: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		writeHumanReadable(f, info)
		fmt.Println("Información guardada en systeminfo.txt")
	} else {
		printHumanReadable(info)
		fmt.Print("\nPresione ENTER para continuar...")
		var input string
		fmt.Scanln(&input)
	}
}
