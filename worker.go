package main

import (
	"log"
	"pulse/repository"
	"pulse/ws"
	"time"
)

func Run(repo *repository.OutboxRepository, hub *ws.Hub, userRepo *repository.UserRepository) {
	for {
		processes, err := repo.GetProcesses(10)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		for _, p := range processes {
			// TODO: handle p (log / email / ws)
			follower, followerErr := userRepo.Me(p.CausedBy)

			if followerErr != nil {
				log.Println(followerErr)
				continue
			}
			msg := follower.UserName + " started following you"
			//Websocket
			hub.Send(p.ReceivedBy, 1, []byte(msg))
			if err := repo.MarkProcessed(p.ID); err != nil {
				log.Println(err)
			}
		}
		time.Sleep(time.Second)
	}
}
