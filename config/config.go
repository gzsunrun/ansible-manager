package config

import (
	"errors"
	"os"
	"path/filepath"

	"git.gzsunrun.cn/sunrunlib/log"
	"github.com/go-ini/ini"
)

type Config struct {
	AnsibleManager AnsibleManager `ini:"ansible_manager"`
}
type AnsibleManager struct {
	Port       		int    `ini:"port"`
	Concurrent 		int    `ini:"concurrent"`
	WorkPath   		string `ini:"work_path"`
	MysqlURL   		string `ini:"mysql_url"`
	MysqlName   	string `ini:"mysql_name"`
	MysqlUser   	string `ini:"mysql_user"`
	MysqlPassword   string `ini:"mysql_password"`
	S3Status   		bool   `ini:"s3"`
	S3URL      		string `ini:"s3_endpoint"`
	S3Key      		string `ini:"s3_key"`
	S3Secret   		string `ini:"s3_secret"`
	BucketName 		string `ini:"bucket_name"`
	JwtSecret  		string `ini:"jwt_secret"`
}

var Cfg = &Config{}

func NewConfig(file string) error {
	err := LoadConfig(file, Cfg)
	return err
}

func LoadConfig(file string, settings interface{}) error {

	if file != "" {

		absConfPath, err := filepath.Abs(file)
		if err != nil {
			log.Debug(err)
			return err
		}

		if err := ini.MapTo(settings, absConfPath); err != nil {
			log.Debug(err)
			return err
		}

		return nil
	}

	return errors.New("file is nil")
}

func WriteConfig(file string, settings interface{}) error {

	if file != "" {

		cfg := ini.Empty()
		err := ini.ReflectFrom(cfg, settings)
		if err != nil {
			return err
		}

		if file == "-" {
			cfg.WriteTo(os.Stdout)
		} else {
			err = cfg.SaveTo(file)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("file is nil")
}
