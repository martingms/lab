package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

		job := &Job{
			ID:     jobId,
			CmdStr: cmdStr + " " + strings.Join(args, " "),
			Status: "Running",
			cmd:    cmd,
			ch:     make(chan error),
		}
		go job.handleJob()

		jobs = append(jobs, job)
	}

	if len(jobs) == 0 {
		return nil, errors.New("no jobs started, either all failed or no hosts")
	}

	return jobs, nil
}

func (job *Job) handleJob() {
	stdoutFile, err := os.Create(filepath.Join(*outputDir, job.ID+".stdout"))
	if err != nil {
		job.ch <- err
		return
	}
	defer stdoutFile.Close()

	stderrFile, err := os.Create(filepath.Join(*outputDir, job.ID+".stderr"))
	if err != nil {
		job.ch <- err
		return
	}
	defer stderrFile.Close()

	var wg sync.WaitGroup
	writeStream := func(file *os.File) {
		n, err := io.Copy(file, job.cmd.StdoutPipe)
		if err != nil {
			log.Println("some error in writing", file.Name()+":", err)
		}
		if n == 0 {
			os.Remove(file.Name())
		} else {
			log.Println(n, "bytes written to", file.Name())
		}

		wg.Done()
	}

	wg.Add(2)
	go writeStream(stdoutFile)
	go writeStream(stderrFile)

	wg.Wait()

	job.ch <- job.cmd.Wait()
}
