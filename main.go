package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jawher/mow.cli"
	"os"
	"path/filepath"
	"strings"
)

type web struct {
	Host string      `json:"host,omitempty"`
	Http *httpStruct `json:"HTTP,omitempty"`
}

type httpStruct struct {
	Port int `json:"port,omitempty"`
}

type cbApp struct {
	Cfengine string `json:"cfengine"`
}

func main() {
	app := cli.App("cbh", "CommandBox server.json helper")

	port := app.IntOpt("p port", 0, "the port for CommandBox to run the application on")
	name := app.StringOpt("n name", "", "the app name for CommandBox to use")
	host := app.StringOpt("H host", "", "the host name, or address, for CommandBox to use")
	overwrite := app.BoolOpt("o overwrite", false, "overwrite server.json file if found in directory")
	destination := app.StringArg("PATH", ".", "the destination to make the file")

	app.Action = func() {
		config := make(map[string]interface{})
		web := web{}
		assigned := false
		if *name != "" {
			config["name"] = *name
		}

		if *host != "" && *host != "127.0.0.1" {
			web.Host = *host
			assigned = true
		}

		if *port > 0 {
			web.Http = &httpStruct{Port: *port}
			assigned = true
		}

		config["app"] = cbApp{Cfengine: "adobe@2018"}
		if assigned {
			config["web"] = web
		}

		j, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			panic(err)
		}

		dest, err := writeConfig(*destination, j, *overwrite)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}

		fmt.Println("server.json file successfully created at", dest)
	}

	err := app.Run(os.Args)

	if err != nil {
		panic(err)
	}
}

func writeConfig(dest string, con []byte, ow bool) (string, error) {
	path, err := filepath.Abs(dest)
	if err != nil {
		return "", nil
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return "", errors.New("specified directory does not exist")
	}

	ogPath := path
	path = filepath.Join(path, "server.json")
	_, err = os.Stat(path)
	if err == nil && !ow {
		if !prompt("A server.json file already exists at " + ogPath + ", do you want to overwrite it") {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}

	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// close the file with defer
	defer f.Close()

	_, err = f.Write(con)
	if err != nil {
		return "", err
	}

	return path, nil
}

func prompt(q string) bool {
	var response string
	quest := fmt.Sprint(q, "(y/N)?")
	fmt.Println(quest)
	for {
		_, err := fmt.Scanln(&response)
		if err != nil {
			panic(err)
		}
		response = strings.ToUpper(strings.TrimSpace(response))
		if response == "Y" {
			return true
		} else if response == "N" {
			return false
		} else {
			fmt.Println("Please enter either 'y' or 'N'")
		}
	}
}
