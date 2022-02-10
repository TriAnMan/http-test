package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/TriAnMan/http-test/worker"
)

func main() {
	var workers int
	flag.IntVar(&workers, "parallel", 10, "maximum number of parallel http requests")
	flag.Parse()

	if workers < 1 {
		log.Fatal("`parallel` must be a positive integer")
	}

	jobs := dispatchJob(workers, flag.Args())
	results := processJob(workers, jobs)

	// wait for results processing
	<-processResult(results)
}

func dispatchJob(workers int, urls []string) chan worker.Job {
	jobs := make(chan worker.Job, workers)

	go func() {
		for _, rawUrl := range urls {
			jobs <- worker.Job{
				Url: rawUrl,
			}
		}
		close(jobs)
	}()

	return jobs
}

func processJob(workers int, jobs chan worker.Job) (results chan worker.Job) {
	results = make(chan worker.Job, workers)

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	// init workers
	for i := 0; i < workers; i++ {
		go func() {
			for job := range jobs {
				results <- worker.Do(job)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return
}

func processResult(results <-chan worker.Job) (done chan struct{}) {
	// start results reader
	done = make(chan struct{})

	go func() {
		logStdout := log.New(os.Stdout, log.Prefix(), 0)

		for result := range results {
			if result.Error != nil {
				log.Printf("%s %s", result.Url, result.Error.Error())
				continue
			}

			logStdout.Printf("%s %x", result.Url, result.Hash)
		}
		close(done)
	}()

	return done
}
