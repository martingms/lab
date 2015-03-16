package main

import (
	"flag"
	"log"
)

var (
	configPath = flag.String("configPath", "labdrc.json", "The path to the JSON config file")
	outputDir  = flag.String("outputDir", "laboutput", "The directory in which to save job output")
)

func init() {
	flag.Parse()
}

var (
	jobs []*Job
)

func main() {
	conf, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	newJobs, err := StartJobs(conf.Hosts, "uname", "-a")
	if err != nil {
		log.Fatal(err)
	}

	jobs = append(jobs, newJobs...)

	err = <-jobs[0].ch
	if err != nil {
		log.Fatal(err)
	}
}
