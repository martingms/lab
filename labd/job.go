package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Job struct {
	ID     string
	CmdStr string
	Status string
	cmd    *Command
	ch     chan error
}

// TODO: Run same command n times on same host, etcetera.
func StartJobs(hosts []*Host, cmdStr string, args ...string) ([]*Job, error) {
	jobs := []*Job{}
	for _, host := range hosts {
		cmd, err := host.StartCmd(cmdStr, args...)
		if err != nil {
			log.Println("unable to start command:", cmdStr, strings.Join(args, " "), "on host:", host.Name, ":", err)
			continue
		}
		log.Println("successfully started command:", cmdStr, strings.Join(args, " "), "on host:", host.Name)

		jobId, err := GenUUID()
		if err != nil {
			log.Println("unable to generate UUID for job:", err)
		}

		job := &Job{ID: jobId, CmdStr: cmdStr + strings.Join(args, " "), Status: "Running", cmd: cmd}
		job.handleJob()

		jobs = append(jobs, job)
	}

	if len(jobs) == 0 {
		return nil, errors.New("no jobs started, either all failed or no hosts")
	}

	return jobs, nil
}

func (job *Job) handleJob() {
	job.ch = make(chan error)

	go func() {
		b, err := ioutil.ReadAll(job.cmd.StdoutPipe)
		if err != nil {
			job.ch <- err
			return
		}

		fmt.Println(string(b))

		job.ch <- job.cmd.Wait()
	}()
}
