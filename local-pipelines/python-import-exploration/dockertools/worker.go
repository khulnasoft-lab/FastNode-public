package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/khulnasoft-lab/fastnode/local-pipelines/python-import-exploration/internal/docker"
)

const (
	statusInterval = 5 * time.Minute
	minStatus      = 15 * time.Minute
	timeout        = time.Hour
)

func buildImage(ctx context.Context, machine *docker.Machine, dockerfilePath string) (err error) {
	// docker image name
	name := strings.ToLower(filepath.Base(dockerfilePath))
	return docker.BuildWithContext(ctx, machine, name, dockerfilePath)
}

func buildImages(numgo int, machine *docker.Machine, dockerfiles []string) []namedError {
	start := time.Now()

	var jobs []*job
	for _, dockerfile := range dockerfiles {
		jobs = append(jobs, &job{
			dockerfile: dockerfile,
			machine:    machine,
		})
	}

	queue := make(chan *job)
	go func() {
		for _, job := range jobs {
			queue <- job
		}
		close(queue)
	}()

	var wg sync.WaitGroup
	wg.Add(numgo)
	for i := 0; i < numgo; i++ {
		go func() {
			defer wg.Done()
			for job := range queue {
				run(timeout, job)
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		// signal we are done by closing the channel
		close(done)
	}()

	status := time.NewTicker(statusInterval)
	defer status.Stop()

	// wait until all jobs done, meanwhile print status
	for {
		select {
		case <-status.C:
			printStatus(start, jobs)
		case <-done: // recieve on a closed channel returns zero value immediately
			var errs []namedError
			for _, job := range jobs {
				if job.err.Err != nil {
					errs = append(errs, job.err)
				}
			}
			return errs
		}
	}
}

func printStatus(start time.Time, jobs []*job) {
	var done int
	var long []string
	var erred int
	for _, job := range jobs {
		js, jd := job.Status()
		switch {
		case !jd.IsZero():
			done++
			if job.err.Err != nil {
				erred++
			}
		case !js.IsZero() && time.Since(js) > minStatus:
			long = append(long, fmt.Sprintf("job %s has taken %v", job.dockerfile, time.Since(js)))
		}
	}
	fmt.Printf("Finished %d (%d errored) jobs in %v\n", done, erred, time.Since(start))
	for _, msg := range long {
		fmt.Println(msg)
	}
	fmt.Println(strings.Repeat("*", 20))
}

type job struct {
	// read only
	dockerfile string
	machine    *docker.Machine

	// only written once job is done running
	// should only be read once job is completed.
	err namedError

	mu    sync.Mutex
	start time.Time // guarded by mu
	done  time.Time // guarded by mu
}

func (j *job) Status() (time.Time, time.Time) {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.start, j.done
}

func (j *job) Start() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.start = time.Now()
}

func (j *job) Done() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.done = time.Now()
}

func run(timeout time.Duration, j *job) {
	j.Start()
	defer j.Done()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := buildImage(ctx, j.machine, j.dockerfile); err != nil {
		j.err = namedError{
			Name: j.dockerfile,
			Err:  fmt.Errorf("error building image: %v", err),
		}
	}
}
