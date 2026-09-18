package main

import (
	"log"
	"pulse/repository"
	"time"
)

const (
	stuckProcessingAfter = 2 * time.Minute
	sweeperInterval      = 30 * time.Second
)

func RunSweeper(repo *repository.OutboxRepository) {
	for {
		n, err := repo.RequeueStuckProcessing(stuckProcessingAfter)
		if err != nil {
			log.Println("outbox sweeper:", err)
		} else if n > 0 {
			log.Printf("outbox sweeper: requeued %d stuck processing row(s)", n)
		}
		time.Sleep(sweeperInterval)
	}
}
