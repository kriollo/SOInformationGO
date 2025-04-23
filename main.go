package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"io"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemInfo struct {
	Hostname    string                 `json:"hostname"`
	OS          string                 `json:"os"`
	Platform    string                 `json:"platform"`
	PlatformVer string                 `json:"platform_version"`
	Arch        string                 `json:"arch"`
	CPU         []cpu.InfoStat         `json:"cpu"`
	CPUTimes    []cpu.TimesStat        `json:"cpu_times"`
	CPUCores    int                    `json:"cpu_cores"`
	Memory      *mem.VirtualMemoryStat `json:"memory"`
	Disk        []disk.UsageStat       `json:"disk"`
	Uptime      uint64                 `json:"uptime_seconds"`
	BootTime    uint64                 `json:"boot_time"`
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

	return &SystemInfo{
		Hostname:    hostInfo.Hostname,
		OS:          hostInfo.OS,
		Platform:    hostInfo.Platform,
		PlatformVer: hostInfo.PlatformVersion,
		Arch:        hostInfo.KernelArch,
		CPU:         cpuInfo,
		CPUTimes:    cpuTimes,
		CPUCores:    cpuCores,
		Memory:      memInfo,
		Disk:        diskUsages,
		Uptime:      hostInfo.Uptime,
		BootTime:    hostInfo.BootTime,
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
	fmt.Fprintf(w, "Uptime: %s\n", humanDuration(info.Uptime))
	fmt.Fprintf(w, "Boot Time: %s\n", humanBootTime(info.BootTime))
	fmt.Fprintf(w, "\nCPU Info:\n")
	for _, cpu := range info.CPU {
		fmt.Fprintf(w, "  Model: %s, Cores: %d, Mhz: %.2f, Cache Size: %s\n", cpu.ModelName, cpu.Cores, cpu.Mhz, humanBytes(uint64(cpu.CacheSize)*1024))
	}
	fmt.Fprintf(w, "  Logical CPUs: %d\n", info.CPUCores)
	fmt.Fprintf(w, "\nMemory:\n")
	fmt.Fprintf(w, "  Total: %s\n", humanBytes(info.Memory.Total))
	fmt.Fprintf(w, "  Used: %s\n", humanBytes(info.Memory.Used))
	fmt.Fprintf(w, "  Free: %s\n", humanBytes(info.Memory.Available))
	fmt.Fprintf(w, "\nDisk(s):\n")
	for _, d := range info.Disk {
		fmt.Fprintf(w, "  Mount: %s, Total: %s, Used: %s, Free: %s, FS: %s\n", d.Path, humanBytes(d.Total), humanBytes(d.Used), humanBytes(d.Free), d.Fstype)
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
	}
}
