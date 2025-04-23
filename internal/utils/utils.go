package utils

import (
	"fmt"
	"time"
)

func HumanBytes(bytes uint64) string {
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

func HumanDuration(seconds uint64) string {
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

func HumanBootTime(boot uint64) string {
	return time.Unix(int64(boot), 0).Format("2006-01-02 15:04:05")
}
