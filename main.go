package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Davincible/goinsta/v3"
	"github.com/gin-gonic/gin"
)

var (
	username                = flag.String("username", os.Getenv("INSTAGRAM_USERNAME"), "Instagram username")
	password                = flag.String("password", os.Getenv("INSTAGRAM_PASSWORD"), "Instagram password")
	twoFactorAuthentication = flag.String("twofactorauthentication", os.Getenv("INSTAGRAM_2FA"), "Instagram 2FA code")
	cert                    = flag.String("cert", "./keys/cert.pem", "Specify the path of TLS cert")
	key                     = flag.String("key", "./keys/key.pem", "Specify the path of TLS key")
	httpsPort               = flag.Int("httpsport", 443, "HTTPS port")
	httpPort                = flag.Int("httpport", 80, "HTTP port")
	insta                   *goinsta.Instagram
	searchbar               *goinsta.Search
	botconfigs              *BotConfigs
	server                  *http.Server
	router                  *gin.Engine
)

func init() {
	botconfigs.SetTimeBetweenRequestsMax(*flag.Int("timebetweenrequestsmax", 15, "Maximum time between requests"))
	botconfigs.SetTimeBetweenRequestsMin(*flag.Int("timebetweenrequestsmin", 5, "Minimum time between requests"))
	botconfigs.SetLikeProbability(*flag.Int("likeprobability", 2, "Probability of liking a post,default is 1 in 2"))
	botconfigs.SetCommentsProbability(*flag.Int("commentsprobability", 4, "Probability of commenting on a post,default is 1 in 4"))
	gin.SetMode(gin.ReleaseMode)
	// Logging to a file.
	gin.DisableConsoleColor() // Disable Console Color, you don't need console color when writing the logs to file.
	path := fmt.Sprintf("private/logs/%02dy_%02dm_%02dd_%02dh_%02dm_%02ds.log", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	logerror1 := os.MkdirAll("private/logs/", 0755)
	f, logerror2 := os.Create(path)
	if logerror1 != nil || logerror2 != nil {
		log.Println("Cant log to file")
	} else {
		//gin.DefaultWriter = io.MultiWriter(f)
		// Use the following code if you need to write the logs to file and console at the same time.
		gin.DefaultWriter = io.MultiWriter(f, os.Stderr)
	}
}

func login() {
	if _, err := os.Stat("~/.goinsta"); err == nil {
		if insta, err = goinsta.Import("~/.goinsta"); err != nil {
			log.Println(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		insta := goinsta.New(*username, *password, *twoFactorAuthentication)
		if err := insta.Login(); err != nil {
			log.Println(err)
		}
	}
}

func iteration() {
	for _, v := range botconfigs.GetHashtagsToSearch() {
		time.Sleep(time.Duration(rand.Intn(botconfigs.GetTimeBetweenRequestsMax()-botconfigs.GetTimeBetweenRequestsMin())+botconfigs.GetTimeBetweenRequestsMin()) * time.Second)
		if result, errors := searchbar.SearchHashtag(v); errors != nil {
			var hashtag *goinsta.Hashtag
			for _, tag := range result.Tags {
				if tag.Name == v {
					result.RegisterHashtagClick(tag)
					hashtag = tag
					break
				}
			}
			if !hashtag.Next() {
				panic(hashtag.Error())
			}
			for _, item := range hashtag.Items {
				time.Sleep(time.Duration(rand.Intn(botconfigs.GetTimeBetweenRequestsMax()-botconfigs.GetTimeBetweenRequestsMin())+botconfigs.GetTimeBetweenRequestsMin()) * time.Second)
				rand.Seed(time.Now().UnixNano())
				if rand.Intn(botconfigs.GetLikeProbability()) == 1 {
					item.Like()
				}
				rand.Seed(time.Now().UnixNano())
				if rand.Intn(botconfigs.GetCommentsProbability()) == 1 {
					item.Comment(botconfigs.GetCommentsToUse()[rand.Intn(len(botconfigs.GetCommentsToUse()))])
				}
			}
		}
	}
}

func main() {
	flag.Parse()
	login()
	iteration()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	defer insta.Export("~/.goinsta")
}
