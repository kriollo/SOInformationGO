package info

import (
	"net"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"systemInformation/internal/model"
)

func GetSystemInfo() (*model.SystemInfo, error) {
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

	ips := []model.NetIPInfo{}
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
					ips = append(ips, model.NetIPInfo{Interface: iface.Name, IP: ipstr})
				}
			}
		}
	}

	kernel := ""
	if hostInfo.OS == "linux" {
		kernel = hostInfo.KernelVersion
	}
	return &model.SystemInfo{
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
