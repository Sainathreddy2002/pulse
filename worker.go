package main

import (
	"context"
	"encoding/json"
	"log"
	"pulse/email"
	"pulse/repository"
	"pulse/ws"
	"time"

	"github.com/redis/go-redis/v9"
)

func Run(repo *repository.OutboxRepository, hub *ws.Hub, userRepo *repository.UserRepository, mailer *email.Sender, rdb *redis.Client) {
	ctx := context.Background()
	for {
		processes, err := repo.GetProcesses(10)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		for _, p := range processes {
			follower, followerErr := userRepo.Me(p.CausedBy)
			if followerErr != nil {
				log.Println(followerErr)
				continue
			}

			recipient, recipientErr := userRepo.Me(p.ReceivedBy)
			if recipientErr != nil {
				log.Println(recipientErr)
				continue
			}

			msg := follower.UserName + " started following you"
			payload, _ := json.Marshal(map[string]any{
				"userId": p.ReceivedBy,
				"body":   msg,
			})
			if err := rdb.Publish(ctx, "pulse:notify", payload).Err(); err != nil {
				log.Println(err)
			}

			subject := "New follower on Pulse"
			body := msg + "\n\n— Pulse"
			if err := mailer.Send(recipient.Email, subject, body); err != nil {
				log.Println("email send failed:", err)
				continue
			}

			if err := repo.MarkProcessed(p.ID); err != nil {
				log.Println(err)
			}
		}
		time.Sleep(time.Second)
	}
}
