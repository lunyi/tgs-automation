package util

import (
	"fmt"
	"os"
	"tgs-automation/internal/log"

	"gopkg.in/yaml.v3"
)

type TgsConfig struct {
	CloudflareToken  string             `yaml:"cloudflare_token"`
	ChromeDriverPath string             `yaml:"chrome_driver_path"`
	GoogleSheet      GoogleSheetConfig  `yaml:"google_sheet"`
	CdnNetwork       CdnConfig          `yaml:"cdn_network`
	Namecheap        NamecheapConfig    `yaml:"namecheap"`
	Postgresql       PostgresqlConfig   `yaml:"postgresql"`
	CreateSiteDb     PostgresqlConfig   `yaml:"create_site_db"`
	Telegram         TelegramConfig     `yaml:"telegram"`
	AwsS3            AwsS3Config        `yaml:"aws_s3"`
	MomoTelegram     MomoTelegramConfig `yaml:"momo_telegram"`
	LetsTalk         LetsTalkConfig     `yaml:"letstalk"`
	Dockerhub        DockerhubConfig    `yaml:"dockerhub"`
}

type DockerhubConfig struct {
	BaseUrl  string `yaml:"base_url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MomoTelegramConfig struct {
	Token       string `yaml:"token"`
	Movn2ChatId int64  `yaml:"movn2_chat_id"`
	MophChatId  int64  `yaml:"moph_chat_id"`
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

type LetsTalkConfig struct {
	AccountId string `yaml:"account_id"`
	ApiKey    string `yaml:"api_key"`
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
