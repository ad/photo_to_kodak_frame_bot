package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

const ConfigFileName = "/data/options.json"

// Config ...
type Config struct {
	TelegramToken        string  `json:"TELEGRAM_TOKEN"`
	TelegramAdminIDs     string  `json:"TELEGRAM_ADMIN_IDS"`
	TelegramAdminIDsList []int64 `json:"-"`

	TargetEmail string `json:"TARGET_EMAIL"`
	FromEmail   string `json:"FROM_EMAIL"`

	SMTP_HOST string `json:"SMTP_HOST"`
	SMTP_PORT int    `json:"SMTP_PORT"`
	SMTP_USER string `json:"SMTP_USER"`
	SMTP_PASS string `json:"SMTP_PASS"`

	Debug bool `json:"DEBUG"`
}

func InitConfig(args []string) (*Config, error) {
	var config = &Config{
		TelegramToken:        "",
		TelegramAdminIDs:     "",
		TelegramAdminIDsList: []int64{},

		TargetEmail: "",
		FromEmail:   "",

		SMTP_HOST: "smtp.gmail.com",
		SMTP_PORT: 587,
		SMTP_USER: "",
		SMTP_PASS: "",

		Debug: false,
	}

	var initFromFile = false

	if _, err := os.Stat(ConfigFileName); err == nil {
		jsonFile, err := os.Open(ConfigFileName)
		if err == nil {
			byteValue, _ := io.ReadAll(jsonFile)
			if err = json.Unmarshal(byteValue, &config); err == nil {
				initFromFile = true
			} else {
				fmt.Printf("error on unmarshal config from file %s\n", err.Error())
			}
		}
	}

	if !initFromFile {
		flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
		flags.StringVar(&config.TelegramToken, "telegramToken", lookupEnvOrString("TELEGRAM_TOKEN", config.TelegramToken), "TELEGRAM_TOKEN")
		flags.StringVar(&config.TelegramAdminIDs, "telegramAdminIDs", lookupEnvOrString("TELEGRAM_ADMIN_IDS", config.TelegramAdminIDs), "TELEGRAM_ADMIN_IDS")

		flags.StringVar(&config.TargetEmail, "targetEmail", lookupEnvOrString("TARGET_EMAIL", config.TargetEmail), "TARGET_EMAIL")
		flags.StringVar(&config.FromEmail, "fromEmail", lookupEnvOrString("FROM_EMAIL", config.FromEmail), "FROM_EMAIL")

		flags.StringVar(&config.SMTP_HOST, "smtpHost", lookupEnvOrString("SMTP_HOST", config.SMTP_HOST), "SMTP_HOST")
		flags.IntVar(&config.SMTP_PORT, "smtpPort", lookupEnvOrInt("SMTP_PORT", config.SMTP_PORT), "SMTP_PORT")
		flags.StringVar(&config.SMTP_USER, "smtpUser", lookupEnvOrString("SMTP_USER", config.SMTP_USER), "SMTP_USER")
		flags.StringVar(&config.SMTP_PASS, "smtpPassword", lookupEnvOrString("SMTP_PASS", config.SMTP_PASS), "SMTP_PASS")

		flags.BoolVar(&config.Debug, "debug", lookupEnvOrBool("DEBUG", config.Debug), "Debug")

		if err := flags.Parse(args[1:]); err != nil {
			return nil, err
		}
	}

	if config.TelegramAdminIDs != "" {
		chatIDS := strings.Split(config.TelegramAdminIDs, ",")
		for _, chatID := range chatIDS {
			if chatIDInt, err := strconv.ParseInt(strings.Trim(chatID, "\n\t "), 10, 64); err == nil {
				config.TelegramAdminIDsList = append(config.TelegramAdminIDsList, chatIDInt)
			}
		}
	}

	return config, nil
}
