package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

const (
	AppName = "System Information App"
	Version = "v1.0.3"
)

type NetIPInfo struct {
	Interface string `json:"interface"`
	IP        string `json:"ip"`
}

type SystemInfo struct {
	Hostname      string                 `json:"hostname"`
	OS            string                 `json:"os"`
	Platform      string                 `json:"platform"`
	PlatformVer   string                 `json:"platform_version"`
	Arch          string                 `json:"arch"`
	KernelVersion string                 `json:"kernel_version,omitempty"`
	IPs           []NetIPInfo            `json:"ips"`
	CPU           []cpu.InfoStat         `json:"cpu"`
	CPUTimes      []cpu.TimesStat        `json:"cpu_times"`
	CPUCores      int                    `json:"cpu_cores"`
	Memory        *mem.VirtualMemoryStat `json:"memory"`
	Disk          []disk.UsageStat       `json:"disk"`
	Uptime        uint64                 `json:"uptime_seconds"`
	BootTime      uint64                 `json:"boot_time"`
}

func getSystemInfo() (*SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	cpuTimes, err := cpu.Times(false)
	if err != nil {
		return nil, err
	}
	cpuCores, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	diskPartitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	diskUsages := []disk.UsageStat{}
	for _, p := range diskPartitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err == nil {
			diskUsages = append(diskUsages, *usage)
		}
	}

	// Obtener IPs activas (no loopback)
	ips := []NetIPInfo{}
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip == nil || ip.IsLoopback() {
					continue
				}
				ipstr := ""
				if ip4 := ip.To4(); ip4 != nil {
					ipstr = ip4.String()
				} else if ip6 := ip.To16(); ip6 != nil {
					ipstr = ip6.String()
				}
				if ipstr != "" {
					ips = append(ips, NetIPInfo{Interface: iface.Name, IP: ipstr})
				}
			}
		}
	}

	kernel := ""
	if hostInfo.OS == "linux" {
		kernel = hostInfo.KernelVersion
	}
	return &SystemInfo{
		Hostname:      hostInfo.Hostname,
		OS:            hostInfo.OS,
		Platform:      hostInfo.Platform,
		PlatformVer:   hostInfo.PlatformVersion,
		Arch:          hostInfo.KernelArch,
		KernelVersion: kernel,
		IPs:           ips,
		CPU:           cpuInfo,
		CPUTimes:      cpuTimes,
		CPUCores:      cpuCores,
		Memory:        memInfo,
		Disk:          diskUsages,
		Uptime:        hostInfo.Uptime,
		BootTime:      hostInfo.BootTime,
	}, nil
}

func printHumanReadable(info *SystemInfo) {
	writeHumanReadable(os.Stdout, info)
}

func writeHumanReadable(w io.Writer, info *SystemInfo) {
	fmt.Fprintf(w, "Hostname: %s\n", info.Hostname)
	fmt.Fprintf(w, "OS: %s\n", info.OS)
	fmt.Fprintf(w, "Platform: %s %s\n", info.Platform, info.PlatformVer)
	fmt.Fprintf(w, "Arch: %s\n", info.Arch)
	if info.KernelVersion != "" {
		fmt.Fprintf(w, "Kernel: %s\n", info.KernelVersion)
	}
	fmt.Fprintf(w, "Uptime: %s\n", humanDuration(info.Uptime))
	fmt.Fprintf(w, "Boot Time: %s\n", humanBootTime(info.BootTime))

	// CPU Info
	fmt.Fprintf(w, "\nCPU Info:\n")
	if len(info.CPU) > 0 {
		cpu := info.CPU[0]
		fmt.Fprintf(w, "  %-15s: %s\n", "Modelo", cpu.ModelName)
		fmt.Fprintf(w, "  %-15s: %d\n", "Cores físicos", cpu.Cores)
		fmt.Fprintf(w, "  %-15s: %d\n", "Cores lógicos", info.CPUCores)
		fmt.Fprintf(w, "  %-15s: %.2f MHz\n", "Frecuencia", cpu.Mhz)
		fmt.Fprintf(w, "  %-15s: %s\n", "Cache", humanBytes(uint64(cpu.CacheSize)*1024))
	}
	if len(info.CPU) > 1 {
		fmt.Fprintf(w, "\n  %-3s %-30s %-7s %-9s %-11s\n", "#", "Modelo", "Cores", "Mhz", "Cache")
		for i, cpu := range info.CPU {
			fmt.Fprintf(w, "  %-3d %-30s %-7d %-9.2f %-11s\n", i+1, cpu.ModelName, cpu.Cores, cpu.Mhz, humanBytes(uint64(cpu.CacheSize)*1024))
		}
	}

	// Memory
	fmt.Fprintf(w, "\nMemory:\n")
	fmt.Fprintf(w, "  Total: %s\n", humanBytes(info.Memory.Total))
	fmt.Fprintf(w, "  Used: %s\n", humanBytes(info.Memory.Used))
	fmt.Fprintf(w, "  Free: %s\n", humanBytes(info.Memory.Available))

	// Disks
	fmt.Fprintf(w, "\nDisk(s):\n")
	if len(info.Disk) > 1 {
		fmt.Fprintf(w, "  %-10s %-12s %-12s %-12s %-6s\n", "Mount", "Total", "Used", "Free", "FS")
		for _, d := range info.Disk {
			fmt.Fprintf(w, "  %-10s %-12s %-12s %-12s %-6s\n", d.Path, humanBytes(d.Total), humanBytes(d.Used), humanBytes(d.Free), d.Fstype)
		}
	} else if len(info.Disk) == 1 {
		d := info.Disk[0]
		fmt.Fprintf(w, "  Mount: %s, Total: %s, Used: %s, Free: %s, FS: %s\n", d.Path, humanBytes(d.Total), humanBytes(d.Used), humanBytes(d.Free), d.Fstype)
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

func humanBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func humanDuration(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	days := d / (24 * time.Hour)
	hours := (d % (24 * time.Hour)) / time.Hour
	minutes := (d % time.Hour) / time.Minute
	sec := (d % time.Minute) / time.Second
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, sec)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, sec)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, sec)
	}
	return fmt.Sprintf("%ds", sec)
}

func humanBootTime(boot uint64) string {
	return time.Unix(int64(boot), 0).Format("2006-01-02 15:04:05")
}

func main() {
	fmt.Printf("\n=============================================================\n")
	fmt.Printf("  %s  |  Versión %s\n", AppName, Version)
	fmt.Printf("  Desarrollado por jorge Jara  |  https://github.com/kriollo\n")
	fmt.Printf("=============================================================\n\n")

	jsonFlag := flag.Bool("json", false, "Exportar salida en formato JSON")
	txtFlag := flag.Bool("txt", false, "Exportar salida legible a systeminfo.txt")
	flag.Parse()

	info, err := getSystemInfo()
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
		fmt.Println(string(output))
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
