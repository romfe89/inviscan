package config

import (
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	OutputDirectoryBase string `mapstructure:"outputDirectoryBase"`
	Wordlists           struct {
		Ffuf string `mapstructure:"ffuf"`
	} `mapstructure:"wordlists"`
	Tools struct {
		Subfinder struct {
			Path string `mapstructure:"path"`
		} `mapstructure:"subfinder"`
		Assetfinder struct {
			Path string `mapstructure:"path"`
		} `mapstructure:"assetfinder"`
		Curl struct {
			Path string `mapstructure:"path"`
		} `mapstructure:"curl"`
		Jq struct {
			Path string `mapstructure:"path"`
		} `mapstructure:"jq"`
		Sed struct {
			Path string `mapstructure:"path"`
		} `mapstructure:"sed"`
		Ffuf struct {
			Path       string `mapstructure:"path"`
			Threads    int    `mapstructure:"threads"`
			MatchCodes string `mapstructure:"matchCodes"`
		} `mapstructure:"ffuf"`
		Httpx struct {
			Path       string `mapstructure:"path"`
			Threads    int    `mapstructure:"threads"`
			MatchCodes string `mapstructure:"matchCodes"`
			Timeout    int    `mapstructure:"timeout"`
		} `mapstructure:"httpx"`
		Gowitness struct {
			Path    string `mapstructure:"path"`
			Threads int    `mapstructure:"threads"`
		} `mapstructure:"gowitness"`
	} `mapstructure:"tools"`
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/inviscan/")
	viper.AddConfigPath("$HOME/.inviscan")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("outputDirectoryBase", "resultados")
	viper.SetDefault("wordlists.ffuf", "/usr/share/seclists/Discovery/DNS/subdomains-top1million-5000.txt")
	viper.SetDefault("tools.subfinder.path", "subfinder")
	viper.SetDefault("tools.assetfinder.path", "assetfinder")
	viper.SetDefault("tools.curl.path", "curl")
	viper.SetDefault("tools.jq.path", "jq")
	viper.SetDefault("tools.sed.path", "sed")
	viper.SetDefault("tools.ffuf.path", "ffuf")
	viper.SetDefault("tools.ffuf.threads", 40)
	viper.SetDefault("tools.ffuf.matchCodes", "200,301,302,307,403")
	viper.SetDefault("tools.httpx.path", "httpx")
	viper.SetDefault("tools.httpx.threads", 50)
	viper.SetDefault("tools.httpx.matchCodes", "200,301,302,307")
	viper.SetDefault("tools.httpx.timeout", 5000)
	viper.SetDefault("tools.gowitness.path", "gowitness")
	viper.SetDefault("tools.gowitness.threads", 4)

	viper.SetEnvPrefix("INVISCAN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			utils.LogWarn("Arquivo de configuração não encontrado. Usando padrões e variáveis de ambiente.")
		} else {
			utils.LogError("Erro ao ler arquivo de configuração: " + err.Error())
			return err
		}
	} else {
		utils.LogInfo("Usando arquivo de configuração: " + viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		utils.LogError("Erro ao fazer unmarshal da configuração: " + err.Error())
		return err
	}
	return nil
}

func GetConfig() *Config {
	return &AppConfig
}
