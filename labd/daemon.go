package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	configPath = flag.String("configPath", "labdrc.json", "The path to the JSON config file")
)

func init() {
	flag.Parse()
}

func main() {
	conf, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	cmd, err := conf.Hosts[0].Run("uname", "-a")
	if err != nil {
		log.Fatal(err)
	}

	out, err := ioutil.ReadAll(cmd.StdoutPipe)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
