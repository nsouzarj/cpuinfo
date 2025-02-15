//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func main() {
	fmt.Println("Iniciando coleta de informações do sistema...")

	// Coleta de informações do sistema
	printSystemInfo()
	printOtherInfo()
	printNetInfo()

	fmt.Println("\nColeta de informações completa.")
}

func printSystemInfo() {
	fmt.Println("\n--- Informações Gerais do Sistema ---")

	// Informações do disco
	disks, err := disk.Partitions(true) // set `all` to true to include all partitions.
	if err != nil {
		log.Println("Erro ao obter partições do disco:", err)
		fmt.Println("Erro ao obter informações do disco")
	} else {
		for _, disk := range disks {
			usage, err := disk.Usage(disk.Mountpoint)
			if err != nil {
				log.Printf("Erro ao obter uso para %s: %v", disk.Mountpoint, err)
				continue // Skip this disk if we can't get its usage.
			}
			fmt.Printf("Disco: %s - %s - %s - Total: %s - Livre: %s\n", disk.Device, disk.Mountpoint, disk.Fstype, humanize.Bytes(usage.Total), humanize.Bytes(usage.Free))
		}
	}

	// Informações da rede
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("Erro ao obter interfaces de rede:", err)
		fmt.Println("Erro ao obter informações de rede")
	} else {
		for _, iface := range ifaces {
			// Obtenha a interface de rede real do pacote 'net'
			netIface, err := net.InterfaceByName(iface.Name)
			if err != nil {
				log.Printf("Erro ao obter informações para a interface %s: %v\n", iface.Name, err)
				continue
			}

			// Verifique se a interface está ativa usando o método Up()
			isUp := (netIface.Flags & net.FlagUp) != 0

			var addrString string
			if len(iface.Addrs) > 0 {
				addrString = iface.Addrs[0].Addr
			} else {
				addrString = "No address"
			}

			if isUp {
				fmt.Printf("Interface: %s - %s\n", iface.Name, addrString)
			} else {
				fmt.Printf("Interface: %s - DOWN\n", iface.Name)
			}
		}
	}

	// Informações da memória
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Erro ao obter informações de memória:", err)
		fmt.Println("Erro ao obter informações de memória")
	} else {
		fmt.Printf("Memória: %s\n", humanize.Bytes(memInfo.Total))
	}

	// Informações da CPU
	percentualDeUso, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Erro ao obter uso da CPU:", err)
		fmt.Println("Erro ao obter uso da CPU")
	} else {
		fmt.Printf("CPU: %.2f%%\n", percentualDeUso[0])
	}
}

func printOtherInfo() {
	fmt.Println("\n--- Outras Informações do Sistema ---")

	cpuinfo, err := cpu.Info()
	if err != nil {
		log.Println("Erro ao obter informações da CPU:", err)
		fmt.Println("Erro ao obter informações da CPU")
	} else {
		if len(cpuinfo) > 0 {
			fmt.Printf("Fabricante: - %s\n", cpuinfo[0].VendorID)
			fmt.Printf("Modelo: - %s\n", cpuinfo[0].Model)
			fmt.Printf("Mhz: - %.2f\n", cpuinfo[0].Mhz)
			fmt.Printf("Quant Cores: - %d\n", cpuinfo[0].Cores)
			fmt.Printf("Familia: - %s\n", cpuinfo[0].Family)
			fmt.Printf("Micro Code: - %s\n", cpuinfo[0].Microcode)
			fmt.Printf("Nome do Modelo: - %s\n", cpuinfo[0].ModelName)
		} else {
			fmt.Println("Nenhuma informação da CPU disponível")
		}
	}

	// Add OS information
	fmt.Printf("OS: - %s\n", runtime.GOOS)
	fmt.Printf("Architecture: - %s\n", runtime.GOARCH)
	fmt.Printf("Go Version: - %s\n", runtime.Version())
	fmt.Printf("Number of CPUs: - %d\n", runtime.NumCPU())
	fmt.Printf("Number of Goroutines: - %d\n", runtime.NumGoroutine())
}

func printNetInfo() {
	fmt.Println("\n--- Informações de Rede ---")

	// Get connection information
	cons, err := net.Connections("tcp")
	if err != nil {
		log.Println("Erro ao obter conexões de rede:", err)
		fmt.Println("Erro ao obter conexões de rede")
	} else {
		for _, con := range cons {
			// Skip connections with empty addresses.  This can happen.
			if con.Laddr.IP == "" || con.Raddr.IP == "" {
				continue
			}
			localAddr := con.Laddr.IP
			remoteAddr := con.Raddr.IP
			port := con.Laddr.Port
			fmt.Printf("Endereço Local: %s, Endereço Remoto: %s  Porta: %d, Status: %s\n",
				localAddr, remoteAddr, port, con.Status)
		}
	}
}
