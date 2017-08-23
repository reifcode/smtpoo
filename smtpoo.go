package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/mailproto/smtpd"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "smtpoo"
	app.Version = "0.0.1"
	app.Usage = "a fake SMTP server caching outbound emails on Redis"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port,p",
			Value: 25,
			Usage: "SMTP port",
		},
		cli.StringFlag{
			Name:  "redis-addr",
			Value: "localhost",
			Usage: "Redis address",
		},
		cli.IntFlag{
			Name:  "redis-port",
			Value: 6379,
			Usage: "Redis port",
		},
		cli.IntFlag{
			Name:  "redis-db",
			Value: 0,
			Usage: "Redis db number",
		},
		cli.StringFlag{
			Name:  "redis-pass",
			Value: "",
			Usage: "Redis password",
		},
		cli.IntFlag{
			Name:  "redis-expiration",
			Value: 0,
			Usage: "Redis keys expiration time (seconds)",
		},
	}

	app.Action = func(c *cli.Context) {
		var (
			port            = c.Int("port")
			redisAddr       = c.String("redis-addr")
			redisPort       = c.Int("redis-port")
			redisDB         = c.Int("redis-db")
			redisPass       = c.String("redis-pass")
			redisExpiration = c.Int("redis-expiration")
		)

		log.Printf("connecting to Redis %s:%d (%d)", redisAddr, redisPort, redisDB)
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisAddr, redisPort),
			DB:       redisDB,
			Password: redisPass,
		})
		_, err := client.Ping().Result()
		if err != nil {
			panic(err)
		}

		server := smtpd.NewServer(func(msg *smtpd.Message) error {
			log.Printf("message from %s to %s", msg.From, msg.To)
			data, err := json.Marshal(msg)
			if err != nil {
				return err
			}
			ex := time.Second * time.Duration(redisExpiration)
			_, err = client.Set(makeRedisKey(), data, ex).Result()
			return err
		})

		log.Println("starting SMTP server on port", port)
		server.ListenAndServe(fmt.Sprintf(":%d", port))
	}

	app.Run(os.Args)
}

func makeRedisKey() string {
	return fmt.Sprintf("mail:%d", time.Now().UnixNano())
}
