package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	Name string
)

func init(){
	flag.StringVar(&Name, "n", "Centos", "Enter the name of the parallels vm")
	flag.StringVar(&Name, "name", "Centos", "Enter the name of the parallels vm")
}

func main(){
	err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(-1)
	}
}

func run() error {
	cmd := exec.Command("prlctl", []string{"list",  "-i", "-j", "name", Name}...)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	type VM struct {
		Hardware struct {
			Net struct {
				Mac string `json:"mac"`
			} `json:"net0"`
		} `json:"Hardware"`
	}

	out := make([]VM, 0)

	err = json.Unmarshal(output, &out)
	if err != nil {
		return err
	}

	mac := out[0].Hardware.Net.Mac
	buf := bytes.NewBuffer([]byte{})

	for i, c := range mac {
		if i % 2 == 0 && i != 0 {
			buf.WriteRune(':')
		}
		buf.WriteRune(c)
	}


	if string(mac[0]) == string('0') {
		mac = buf.String()[1:]
	}
	command := strings.ToLower(fmt.Sprintf("arp -an | grep \"%s\" | awk '{ print $2 }'", mac))
	cmd = exec.Command("bash",  "-c", command)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Print(string(output[1:len(output) - 2]))
	return nil
}