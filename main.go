package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

type icon struct {
	Base64  string
	Decoded []byte
}

var iconBase string = "PHN2ZyB2ZXJzaW9uPSIxLjEiIHZpZXdCb3g9IjAgMCAxNiAxNiIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KIDxwYXRoIGQ9Im02LjAwNDEgMHYxLjAxMTdjLTAuNzcxNDggMC4wMTYzOTYtMS4zOTgyIDAuMDc2Mjk0LTEuOTM2NCAwLjM3MzA1YTEuODc4OSAxLjg4IDAgMCAwLTAuODM3NCAwLjk5ODA1Yy0wLjE2NDkgMC40MzctMC4yMjQ0OCAwLjk2MDE5LTAuMjI0NDggMS42MTcydjdoLTAuMDA1ODU1OXYxYzAgMC42NTggMC4wNTk1NzQgMS4xNzkyIDAuMjI0NDggMS42MTcyIDAuMTYzOSAwLjQzOSAwLjQ2MTYyIDAuNzg5MDkgMC44Mzc0IDAuOTk2MDkgMC43NTI1NiAwLjQxNSAxLjY3MjcgMC4zNzE3MiAyLjkyOTkgMC4zODY3MmgyLjAxMDVjMS4yNTczLTAuMDE1IDIuMTc3NCAwLjAyODI4IDIuOTI5OS0wLjM4NjcyIDAuMzc1NzgtMC4yMDcgMC42NzM0OS0wLjU1NzA5IDAuODM3NC0wLjk5NjA5IDAuMTY0ODEtMC40MzggMC4yMjQ0Ny0wLjk1OTE5IDAuMjI0NDctMS42MTcyaDAuMDA1ODk3di04YzAtMC42NTctMC4wNTk1NzUtMS4xODAyLTAuMjI0NDgtMS42MTcyYTEuODc4OSAxLjg4IDAgMCAwLTAuODM3NDQtMC45OTgwNWMtMC41MzgxMy0wLjI5Njc1LTEuMTY0OS0wLjM1NjY1LTEuOTM2NC0wLjM3MzA1di0xLjAxMTd6bTAuOTk5NDEgMmgxLjk5ODhjMS4yNTkzIDAuMDE1IDIuMDg1OSAwLjA1OTcxOSAyLjQ1MTcgMC4yNjE3MiAwLjE4Mzg5IDAuMTAxIDAuMjg4NTUgMC4yMTI2NiAwLjM4NjQ5IDAuNDcyNjYgMC4wOTY5NDMgMC4yNiAwLjE2MDA2IDAuNjczNjIgMC4xNjAwNiAxLjI2NTZ2N2gtMC4wMDU5djFjMCAwLjU5Mi0wLjA2MzEyIDEuMDA1Ni0wLjE2MDA2IDEuMjY1Ni0wLjA5Nzk0IDAuMjU5LTAuMjAwNjUgMC4zNzI2Ni0wLjM4NDU0IDAuNDcyNjYtMC4zNjQ2MiAwLjIwMi0xLjE5NTMgMC4yNDY3Mi0yLjQ1MzYgMC4yNjE3MmgtMS45OTg4Yy0xLjI1ODMtMC4wMTUtMi4wODg4LTAuMDU5NzItMi40NTM2LTAuMjYxNzItMC4xODM4OS0wLjEtMC4yODY2LTAuMjEzNjYtMC4zODQ1NC0wLjQ3MjY2LTAuMDk2OTQzLTAuMjYtMC4xNjAwNi0wLjY3MzYyLTAuMTYwMDYtMS4yNjU2di0xaDAuMDA1ODZ2LTdjMC0wLjU5MiAwLjA2MzExOS0xLjAwNTYgMC4xNjAwNi0xLjI2NTYgMC4wOTc5NDItMC4yNiAwLjIwMjYtMC4zNzE2NiAwLjM4NjQ5LTAuNDcyNjYgMC4zNjU3OC0wLjIwMiAxLjE5MjQtMC4yNDY3MiAyLjQ1MTctMC4yNjE3MnoiIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXciIGZpbGw9IiM4MDgwODAiLz4KIDxwYXRoIG"
var icons = map[int]icon{
	10: {
		Base64: iconBase + "NsYXNzPSJlcnJvciIgZD0ibTUgMTJjMCAwLjU1NCAwLjQ0NiAxIDEgMWg0YzAuNTU0IDAgMS0wLjQ0NiAxLTF6IiBmaWxsPSIjZGExNjM2IiBzdHJva2Utd2lkdGg9IjUiLz4KPC9zdmc+Cg==",
	},
	20: {
		Base64: iconBase + "NsYXNzPSJ3YXJuaW5nIiBkPSJtNSAxMXYxYzAgMC41NTQgMC40NDYgMSAxIDFoNGMwLjU1NCAwIDEtMC40NDYgMS0xdi0xaC02eiIgZmlsbD0iI2ZkYzkyYiIgc3Ryb2tlLXdpZHRoPSI1Ii8+Cjwvc3ZnPgo=",
	},
	30: {
		Base64: iconBase + "Q9Im01IDEwdjJjMCAwLjU1NCAwLjQ0NiAxIDEgMWg0YzAuNTU0IDAgMS0wLjQ0NiAxLTF2LTJoLTZ6IiBmaWxsPSIjODA4MDgwIiBzdHJva2Utd2lkdGg9IjUiLz4KPC9zdmc+Cg==",
	},
	40: {
		Base64: iconBase + "Q9Im01IDl2M2MwIDAuNTU0IDAuNDQ2IDEgMSAxaDRjMC41NTQgMCAxLTAuNDQ2IDEtMXYtM2gtNnoiIGZpbGw9IiM4MDgwODAiIHN0cm9rZS13aWR0aD0iNSIvPgo8L3N2Zz4K",
	},
	50: {
		Base64: iconBase + "Q9Im01IDh2NGMwIDAuNTU0IDAuNDQ2IDEgMSAxaDRjMC41NTQgMCAxLTAuNDQ2IDEtMXYtNGgtNnoiIGZpbGw9IiM4MDgwODAiIHN0cm9rZS13aWR0aD0iNSIvPgo8L3N2Zz4K",
	},
	60: {
		Base64: iconBase + "Q9Im01IDd2NWMwIDAuNTU0IDAuNDQ2IDEgMSAxaDRjMC41NTQgMCAxLTAuNDQ2IDEtMXYtNWgtNnoiIGZpbGw9IiM4MDgwODAiIHN0cm9rZS13aWR0aD0iNSIvPgo8L3N2Zz4K",
	},
	70: {
		Base64: iconBase + "Q9Im01IDZ2NmMwIDAuNTU0IDAuNDQ2IDEgMSAxaDRjMC41NTQgMCAxLTAuNDQ2IDEtMXYtNmgtNnoiIGZpbGw9IiM4MDgwODAiIHN0cm9rZS13aWR0aD0iNSIvPgo8L3N2Zz4K",
	},
	80: {
		Base64: iconBase + "Q9Im01IDV2N2MwIDAuNTU0IDAuNDQ2IDEgMSAxaDRjMC41NTQgMCAxLTAuNDQ2IDEtMXYtN2gtNnoiIGZpbGw9IiM4MDgwODAiIHN0cm9rZS13aWR0aD0iNSIvPgo8L3N2Zz4K",
	},
	90: {
		Base64: iconBase + "Q9Im01IDR2OGMwIDAuNTU0IDAuNDQ2IDEgMSAxaDRjMC41NTQgMCAxLTAuNDQ2IDEtMXYtOGgtNnoiIGZpbGw9IiM4MDgwODAiIHN0cm9rZS13aWR0aD0iNSIvPgo8L3N2Zz4K",
	},
	100: {
		Base64: iconBase + "Q9Im02IDNoNGMwLjU1NCAwIDEgMC40NDYgMSAxdjhjMCAwLjU1NC0wLjQ0NiAxLTEgMWgtNGMtMC41NTQgMC0xLTAuNDQ2LTEtMXYtOGMwLTAuNTU0IDAuNDQ2LTEgMS0xeiIgZmlsbD0iIzgwODA4MCIgc3Ryb2tlLXdpZHRoPSI1Ii8+Cjwvc3ZnPgo=",
	},
}

func main() {
	syslog, err := syslog.New(syslog.LOG_INFO, "battery-percentage")
	if err != nil {
		panic("Unable to connect to syslog")
	}
	log.SetOutput(syslog)

	flag.Parse()

	systray.Run(func() {
		quit := systray.AddMenuItem("Quit", "")
		go func() {
			for {
				select {
				case <-quit.ClickedCh:
					systray.Quit()
				}
			}
		}()
		setPercentage()
		for range time.Tick(time.Minute * 1) {
			setPercentage()
		}
	}, func() {})
}

func getIcon(percentage string) []byte {
	p, err := strconv.ParseFloat(percentage, 64)
	if err != nil {
		log.Println(err)
	}
	if p < 1 {
		p = 10
	}
	img, err := decodedIcon(int(math.Ceil(p/10) * 10))
	if err != nil {
		log.Println(err)
	}

	return img
}

func decodedIcon(p int) ([]byte, error) {
	if len(icons[p].Decoded) < 1 {
		img, err := base64.StdEncoding.DecodeString(icons[p].Base64)
		if err != nil {
			return []byte(" "), fmt.Errorf("failed to get icon: %d", p)
		}
		icons[p] = icon{
			Base64:  icons[p].Base64,
			Decoded: img,
		}
	}

	return icons[p].Decoded, nil
}

func setPercentage() {
	cmd := "bluetoothctl info | grep 'Battery Percentage' | cut -d ' ' -f4 | sed 's/^.//;s/.$//'"
	rawPercentage, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Println(err)
		systray.SetTitle("")
		systray.SetIcon([]byte("test"))
		return
	}
	percentage := strings.TrimSuffix(string(rawPercentage), "\n")
	systray.SetTitle(percentage + "%")
	icon := getIcon(percentage)
	systray.SetIcon(icon)
}
