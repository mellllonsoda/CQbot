package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	quotes   map[string]string
	keywords map[string][]string
)

func main() {
	// .envがあれば読み込むが、なくてもエラー（Fatal）にしない
	_ = godotenv.Load(".env")

	token := os.Getenv("DISCORD_TOKEN")
	guildID := os.Getenv("DEV_GUILD_ID")

	if token == "" {
		// トークンが「環境変数」としても存在しない場合のみ終了する
		log.Fatalf("DISCORD_TOKEN not set in environment variables")
	}

	// 語録とキーワードを読み込み
	quotesFile, err := ioutil.ReadFile("quotes.json")
	if err != nil {
		log.Fatalf("Error reading quotes.json: %v", err)
	}
	json.Unmarshal(quotesFile, &quotes)

	keywordsFile, err := ioutil.ReadFile("keywords.json")
	if err != nil {
		log.Fatalf("Error reading keywords.json: %v", err)
	}
	json.Unmarshal(keywordsFile, &keywords)

	// Discordセッションを作成
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// イベントハンドラを登録
	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(reactionAdd)
	dg.AddHandler(interactionCreate)

	// Intentsを設定
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	// 接続を開く
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	// スラッシュコマンドを登録
	registerCommands(dg, guildID)

	// Botが終了するまで待機
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// セッションを閉じる
	dg.Close()
}

// Botの準備ができたら呼ばれる
func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	log.Println("スラッシュコマンドを同期しました")
}

// メッセージが作成されたら呼ばれる
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 自分のメッセージは無視
	if m.Author.ID == s.State.User.ID {
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		return
	}

	// キーワードを探す
	var matchedIDs []string
	for kw, ids := range keywords {
		if strings.Contains(m.Content, kw) {
			matchedIDs = append(matchedIDs, ids...)
		}
	}

	// 10%の確率で反応
	if len(matchedIDs) > 0 && rand.Float32() < 0.1 {
		rand.Seed(time.Now().UnixNano())
		selectedID := matchedIDs[rand.Intn(len(matchedIDs))]
		quote, ok := quotes[selectedID]
		if !ok {
			quote = "（該当する語録が見つかりませんでした）"
		}
		embed := &discordgo.MessageEmbed{
			Description: quote,
			Color:       0x1abc9c,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}

// リアクションが追加されたら呼ばれる
func reactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Botのリアクションは無視
	if r.UserID == s.State.User.ID {
		return
	}
	// ❌リアクションでない場合は無視
	if r.Emoji.Name == "❌" {
		// メッセージの投稿者がBot自身であるかを確認
		msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
		if err != nil {
			log.Printf("Failed to get message: %v", err)
			return
		}
		if msg.Author.ID == s.State.User.ID {
			err := s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			if err != nil {
				log.Printf("Failed to delete message: %v", err)
			}
		}
	}
}

// スラッシュコマンドのハンドラ
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "random_quote":
			handleRandomQuote(s, i)
		case "revolutionized":
			handleRevolutionized(s, i)
		case "ping":
			handlePing(s, i)
		}
	}
}

func handleRandomQuote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rand.Seed(time.Now().UnixNano())
	keys := make([]string, 0, len(quotes))
	for k := range quotes {
		keys = append(keys, k)
	}
	randomID := keys[rand.Intn(len(keys))]
	quote := quotes[randomID]
	embed := &discordgo.MessageEmbed{
		Description: quote,
		Color:       0x1abc9c,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func handleRevolutionized(s *discordgo.Session, i *discordgo.InteractionCreate) {
	message := i.ApplicationCommandData().Options[0].StringValue()
	transformed := strings.Join(strings.Split(message, ""), "☆")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🔴☭%s☭🔴", transformed),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	start := time.Now()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pinging...",
		},
	})
	if err != nil {
		log.Printf("Failed to send ping response: %v", err)
		return
	}
	latency := time.Since(start)

	// followup messageを送信
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("Pong! レイテンシ: %s", latency),
	})
}

// スラッシュコマンドの定義と登録
func registerCommands(s *discordgo.Session, guildID string) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "random_quote",
			Description: "ランダムに名言を出す",
		},
		{
			Name:        "revolutionized",
			Description: "入力を変換する",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "変換するメッセージ",
					Required:    true,
				},
			},
		},
		{
			Name:        "ping",
			Description: "Bot の応答速度を測定",
		},
	}

	if guildID != "" {
		log.Printf("ローカルモード: ギルド %s にコマンドを登録します", guildID)
	} else {
		log.Println("本番モード: グローバルにコマンドを登録します")
	}

	for _, cmd := range commands {
		// guildID が空ならグローバル、値があればそのサーバー専用になる
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", cmd.Name, err)
		}
	}
}
