//go:build windows

package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func main() {
	// Create a new application
	a := app.New()

	// Create a new window
	w := a.NewWindow("Informacoes da CPI")
	// Set a smaller initial size
	w.Resize(fyne.NewSize(640, 480))

	// Create a new tab container
	tab1 := container.NewTabItem("Info Geral", createSystemInfoTab())
	tab2 := container.NewTabItem("Outras Info", createOtherInfoTab())
	tab3 := container.NewTabItem("Info de Rede", createNetInfoTab())

	// Create a new tab container for the system info
	tabs := container.NewAppTabs(tab1, tab2, tab3)
	w.SetContent(tabs)

	// Show the window
	w.ShowAndRun()
}

func createSystemInfoTab() fyne.CanvasObject {
	// Create a new box container
	box := container.NewVBox()

	// Get disk information
	disks, err := disk.Partitions(false)
	if err != nil {
		log.Println("Error getting disk partitions:", err)
		box.Add(widget.NewLabel("Error getting disk information"))
	} else {
		for _, disk := range disks {
			box.Add(widget.NewLabel(fmt.Sprintf("Disco: %s - %s - %s", disk.Device, disk.Mountpoint, disk.Fstype)))
		}
	}

	// Get network information
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting network interfaces:", err)
		box.Add(widget.NewLabel("Error getting network information"))
	} else {
		for _, iface := range ifaces {
			if len(iface.Addrs) > 0 {
				box.Add(widget.NewLabel(fmt.Sprintf("Interface: %s - %s", iface.Name, iface.Addrs[0].Addr)))
			} else {
				box.Add(widget.NewLabel(fmt.Sprintf("Interface: %s - no addresses", iface.Name)))
			}
		}
	}

	// Get memory information
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting memory information:", err)
		box.Add(widget.NewLabel("Error getting memory information"))
	} else {
		box.Add(widget.NewLabel(fmt.Sprintf("Memória: %s", humanize.Bytes(uint64(memInfo.Total)))))
	}

	// Get CPU information
	percentualDeUso, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Error getting CPU usage:", err)
		box.Add(widget.NewLabel("Error getting CPU usage"))
	} else {
		box.Add(widget.NewLabel(fmt.Sprintf("CPU: %.2f%%", percentualDeUso[0])))
	}

	// Create a scrollable container
	scroll := container.NewScroll(box)

	return scroll
}

func createOtherInfoTab() fyne.CanvasObject {
	// Create a new box container
	box := container.NewVBox()

	// Add other info widgets here
	cpuinfo, err := cpu.Info()
	if err != nil {
		log.Println("Error getting CPU information:", err)
		box.Add(widget.NewLabel("Error getting CPU information"))
	} else {
		if len(cpuinfo) > 0 {
			box.Add(widget.NewLabel(fmt.Sprintf("Fabricante: - %s", cpuinfo[0].VendorID)))
			box.Add(widget.NewLabel(fmt.Sprintf("Modelo - %s", cpuinfo[0].Model)))
			box.Add(widget.NewLabel(fmt.Sprintf("Mhz: - %.2f", cpuinfo[0].Mhz)))
			box.Add(widget.NewLabel(fmt.Sprintf("Quant Cores: - %d", cpuinfo[0].Cores)))
			box.Add(widget.NewLabel(fmt.Sprintf("Familia: - %s", cpuinfo[0].Family)))
			box.Add(widget.NewLabel(fmt.Sprintf("Micro Code: - %s ", cpuinfo[0].Microcode)))
			box.Add(widget.NewLabel(fmt.Sprintf("Nome do Modelo: - %s", cpuinfo[0].ModelName)))
		} else {
			box.Add(widget.NewLabel("No CPU information available"))
		}
	}

	// Create a scrollable container
	scroll := container.NewScroll(box)

	return scroll
}

func createNetInfoTab() fyne.CanvasObject {
	// Create a new box container
	box := container.NewVBox()

	// Get connection information
	cons, err := net.Connections("tcp")
	if err != nil {
		log.Println("Error getting network connections:", err)
		box.Add(widget.NewLabel("Error getting network connections"))
	} else {
		for _, con := range cons {
			localAddr := con.Laddr.IP
			remoteAddr := con.Raddr.IP
			port := con.Laddr.Port
			box.Add(widget.NewLabel(fmt.Sprintf("Endereço Local: %s, Endereço Remoto: %s  Porta: %d",
				localAddr, remoteAddr, port)))
		}
	}

	// Create a scrollable container
	scroll := container.NewScroll(box)

	return scroll
}
