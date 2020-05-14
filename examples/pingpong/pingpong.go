package main

import (
	"flag"
	"log"

	"github.com/warmind-io/dgrouter"

	"github.com/warmind-io/dgrouter/exrouter"
	"github.com/warmind-io/discordgo"
)

// Command line flags
var (
	fToken  = flag.String("t", "", "bot token")
	fPrefix = flag.String("p", "!", "bot prefix")
)

func main() {
	flag.Parse()

	s, err := discordgo.New("Bot " + *fToken)
	if err != nil {
		log.Fatal(err)
	}

	router := exrouter.New()

	// Add some commands
	router.On("ping", func(ctx *exrouter.Context) {
		ctx.Reply("pong")
	}).Desc("responds with pong")

	router.On("avatar", func(ctx *exrouter.Context) {
		ctx.Reply(ctx.Msg.Author.AvatarURL("2048"))
	}).Desc("returns the user's avatar")

	// Match the regular expression user(name)?
	router.OnMatch("username", dgrouter.NewRegexMatcher("user(name)?"), func(ctx *exrouter.Context) {
		ctx.Reply("Your username is " + ctx.Msg.Author.Username)
	})

	router.Default = router.On("help", func(ctx *exrouter.Context) {
		var text = ""
		for _, v := range router.Routes {
			text += v.Name + " : \t" + v.Description + "\n"
		}
		ctx.Reply("```" + text + "```")
	}).Desc("prints this help menu")

	// Add message handler
	s.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(s, *fPrefix, s.State.User.ID, m.Message)
	})

	err = s.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("bot is running...")
	// Prevent the bot from exiting
	<-make(chan struct{})
}
