package util

import (
	"cdnetwork/internal/log"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type TgsConfig struct {
	CloudflareToken  string             `yaml:"cloudflare_token"`
	ChromeDriverPath string             `yaml:"chrome_driver_path"`
	GoogleSheet      GoogleSheetConfig  `yaml:"google_sheet"`
	Cdnetwork        CdnConfig          `yaml:"cdnetwork"`
	Namecheap        NamecheapConfig    `yaml:"namecheap"`
	Postgresql       PostgresqlConfig   `yaml:"postgresql"`
	Telegram         TelegramConfig     `yaml:"telegram"`
	AwsS3            AwsS3Config        `yaml:"aws_s3"`
	MomoTelegram     MomoTelegramConfig `yaml:"momo_telegram"`
}

type MomoTelegramConfig struct {
	Token  string `yaml:"token"`
	ChatId int64  `yaml:"chat_id"`
}

type AwsS3Config struct {
	Region       string `yaml:"region"`
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
	Bucket       string `yaml:"bucket"`
}

type GoogleSheetConfig struct {
	GoogleApiKey string `yaml:"google_api_key"`
	SheetId      string `yaml:"sheet_id"`
}

type CdnConfig struct {
	CdnUserName               string `yaml:"cdn_user_name"`
	CdnApiKey                 string `yaml:"cdn_api_key"`
	CdnEndPoint               string `yaml:"cdn_end_point"`
	CdnLoginUrl               string `yaml:"cdn_login_url"`
	CdnCertificateCreationUrl string `yaml:"cdn_certificate_creation_url"`
	CdnPassword               string `yaml:"cdn_password"`
	DnsContent                string `yaml:"dns_content"`
}

type NamecheapConfig struct {
	NamecheapApiKey   string `yaml:"namecheap_api_key"`
	NamecheapUsername string `yaml:"namecheap_username"`
	NamecheapPassword string `yaml:"namecheap_password"`
	NamecheapClientIp string `yaml:"namecheap_client_ip"`
}

type PostgresqlConfig struct {
	PgHost     string `yaml:"pg_host"`
	PgPort     string `yaml:"pg_port"`
	PgDb       string `yaml:"pg_database"`
	PgUsername string `yaml:"pg_username"`
	PgPassword string `yaml:"pg_password"`
}

type TelegramConfig struct {
	TelegramBotToken string `yaml:"telegram_bot_token"`
	TelegramChatId   string `yaml:"telegram_chat_id"`
	TelegramWebhook  string `yaml:"telegram_webhook"`
}

func GetConfig() TgsConfig {
	data, err := os.ReadFile(os.Getenv("CONFIGPATH"))

	if err != nil {
		log.LogFatal(fmt.Sprintf("Fail to load config: %s  %v", os.Getenv("CONFIGPATH"), err))
		panic(err)
	}

	// create a person struct and deserialize the data into that struct
	var config TgsConfig

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.LogFatal(fmt.Sprintf("Fail to parse config: %v", err))
		panic(err)
	}
	return config
}
