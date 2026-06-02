package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DBPath string
	Auth115 struct {
		DownloadPath string
		AccessToken  string
		RefreshToken string
	}
}
