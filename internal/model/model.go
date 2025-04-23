package model

// NetIPInfo representa una IP de una interfaz de red
 type NetIPInfo struct {
	Interface string `json:"interface"`
	IP        string `json:"ip"`
}

// SystemInfo almacena toda la información del sistema
 type SystemInfo struct {
	Hostname      string      `json:"hostname"`
	OS            string      `json:"os"`
	Platform      string      `json:"platform"`
	PlatformVer   string      `json:"platform_version"`
	Arch          string      `json:"arch"`
	KernelVersion string      `json:"kernel_version,omitempty"`
	IPs           []NetIPInfo `json:"ips"`
	CPU           interface{} `json:"cpu"`
	CPUTimes      interface{} `json:"cpu_times"`
	CPUCores      int         `json:"cpu_cores"`
	Memory        interface{} `json:"memory"`
	Disk          interface{} `json:"disk"`
	Uptime        uint64      `json:"uptime_seconds"`
	BootTime      uint64      `json:"boot_time"`
}
