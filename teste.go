//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	psnet "github.com/shirou/gopsutil/net" // Usando um alias para o pacote gopsutil/net
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
	disks, err := disk.Partitions(true) // Obtém todas as partições do disco
	if err != nil {
		log.Println("Erro ao obter partições do disco:", err)
		fmt.Println("Erro ao obter informações do disco")
	} else {
		for _, d := range disks {
			usage, err := disk.Usage(d.Mountpoint) // Obtém o uso do disco para o ponto de montagem
			if err != nil {
				log.Printf("Erro ao obter uso para %s: %v", d.Mountpoint, err)
				continue // Pula para a próxima partição em caso de erro
			}
			fmt.Printf("Disco: %s - %s - %s - Total: %s - Livre: %s\n",
				d.Device, d.Mountpoint, d.Fstype,
				humanize.Bytes(usage.Total), humanize.Bytes(usage.Free))
		}
	}

	// Informações da rede
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("Erro ao obter interfaces de rede:", err)
		fmt.Println("Erro ao obter informações de rede")
	} else {
		for _, iface := range ifaces {
			// Obtém os endereços da interface
			addrs, err := iface.Addrs()
			if err != nil {
				log.Printf("Erro ao obter endereços para a interface %s: %v\n", iface.Name, err)
				continue // Pula para a próxima interface em caso de erro
			}

			// Verifica se a interface está ativa
			isUp := iface.Flags&net.FlagUp == net.FlagUp

			var addrString string
			if len(addrs) > 0 {
				addrString = addrs[0].String() // Obtém o primeiro endereço da interface
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

	// Informações detalhadas da CPU
	cpuinfo, err := cpu.Info()
	if err != nil {
		log.Println("Erro ao obter informações da CPU:", err)
		fmt.Println("Erro ao obter informações da CPU")
	} else {
		if len(cpuinfo) > 0 {
			fmt.Printf("Fabricante: %s\n", cpuinfo[0].VendorID)
			fmt.Printf("Modelo: %s\n", cpuinfo[0].ModelName)
			fmt.Printf("Mhz: %.2f\n", cpuinfo[0].Mhz)
			fmt.Printf("Quant Cores: %d\n", cpuinfo[0].Cores)
			fmt.Printf("Familia: %s\n", cpuinfo[0].Family)
			fmt.Printf("Micro Code: %s\n", cpuinfo[0].Microcode)
		} else {
			fmt.Println("Nenhuma informação da CPU disponível")
		}
	}

	// Informações do sistema operacional
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())
}

func printNetInfo() {
	fmt.Println("\n--- Informações de Rede ---")

	// Obtém informações sobre conexões de rede
	cons, err := psnet.Connections("tcp") // Usando o alias psnet para o pacote gopsutil/net
	if err != nil {
		log.Println("Erro ao obter conexões de rede:", err)
		fmt.Println("Erro ao obter conexões de rede")
	} else {
		for _, con := range cons {
			// Ignora conexões com endereços IP vazios
			if con.Laddr.IP == "" || con.Raddr.IP == "" {
				continue
			}
			fmt.Printf("Endereço Local: %s, Endereço Remoto: %s, Porta: %d, Status: %s\n",
				con.Laddr.IP, con.Raddr.IP, con.Laddr.Port, con.Status)
		}
	}
}
