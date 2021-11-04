package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type timeValues struct {
	minute int
	second int
}

// checkTime will get the output from running the wg command,
// check if minutes > than max handshake minutes, and execute
// an up/down of the wg0 interface if true.
func checkTime(maxHandshakeMinutes int, intf string) error {
	var err error
	var t timeValues

	text, err := getHandshakeValue(intf)
	if err != nil {
		// There was an error getting the handshake string, and
		// most likely the interface is not up. Set high values
		// so we force a down and then up of the interface.
		t = timeValues{
			minute: 10000,
			second: 10000,
		}
	} else {
		// All ok, parse the string containing the handshake info.
		t, err = parseTimeValues(text)
		if err != nil {
			return err
		}
	}

	if t.minute >= maxHandshakeMinutes {
		// execute wg up/down here...
		log.Printf("not ok, got minute > %v, initiating wg down/up", maxHandshakeMinutes)

		cmdArg := fmt.Sprintf("wg-quick down %v", intf)
		cmd := exec.Command("bash", "-c", cmdArg)
		err := cmd.Run()
		if err != nil {
			log.Printf("error: executing command down of interface: %v", err)
		}

		cmdArg = fmt.Sprintf("wg-quick up %v", intf)
		cmd = exec.Command("bash", "-c", cmdArg)
		err = cmd.Run()
		if err != nil {
			log.Printf("error: executing command down of interface: %v", err)
		}
	} else {
		log.Printf("ok, got minute < %v, currently %v minutes, %v seconds\n", maxHandshakeMinutes, t.minute, t.second)
	}

	return nil
}

// getHandshakeValue gets the current last handshake value from
// running the command 'wg show xxx'
func getHandshakeValue(wgInterface string) (string, error) {
	cmd := exec.Command("bash", "-c", "wg show "+wgInterface)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: executing command: %v", err)
	}

	bReader := bytes.NewReader(b)
	scanner := bufio.NewScanner(bReader)
	var text string

	for scanner.Scan() {
		if v := scanner.Text(); strings.Contains(scanner.Text(), "latest handshake") {
			text = v
		}
	}

	return text, nil
}

// Parse the time values from string variable given, and return
// a struct with the result.
func parseTimeValues(s string) (timeValues, error) {
	t := timeValues{}

	split := strings.Split(s, " ")
	for i, v := range split {
		if strings.Contains(v, "minute") {
			var err error

			t.minute, err = strconv.Atoi(split[i-1])
			if err != nil {
				return t, fmt.Errorf("error: failed to convert minute to int: %v", err)
			}
		}

		if strings.Contains(v, "second") {
			var err error

			t.second, err = strconv.Atoi(split[i-1])
			if err != nil {
				return t, fmt.Errorf("error: failed to convert second to int: %v", err)
			}
		}
	}

	return t, nil
}

func main() {
	maxHandshakeMinutes := flag.Int("handshakeMinutes", 4, "the minutes to wait for handshake before taking action")
	checkSeconds := flag.Int("checkSeconds", 30, "the seconds to wait between doing a check")
	intf := flag.String("intf", "wg0", "the name of the wireguard interface")
	flag.Parse()

	var timedInterval = *checkSeconds

	ticker := time.NewTicker(time.Second * time.Duration(timedInterval))
	for {
		<-ticker.C

		err := checkTime(*maxHandshakeMinutes, *intf)
		if err != nil {
			log.Printf("%v\n", err)
		}
	}
}
