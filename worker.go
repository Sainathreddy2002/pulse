package main

import (
	"log"
	"pulse/repository"
	"time"
)

func Run(repo *repository.OutboxRepository) {
	for {
		processes, err := repo.GetProcesses(10)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		for _, p := range processes {
			// TODO: handle p (log / email / ws)
			if err := repo.MarkProcessed(p.ID); err != nil {
				log.Println(err)
			}
		}
		time.Sleep(time.Second)
	}
}
