package config

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/astaxie/beego/logs"
	"github.com/go-ini/ini"
)

// Config config strcut
type Config struct {
	Common       Common       `ini:"common"`
	Mysql        Mysql        `ini:"mysql"`
	LocalStorage LocalStorage `ini:"local_storage"`
	S3           S3           `ini:"s3_storage"`
	Git          Git          `ini:"git_storage"`
	Etcd         Etcd         `ini:"etcd"`
	FileLog      FileLog      `ini:"file_log"`
}

// Common common config
type Common struct {
	Port       int    `ini:"port"`
	Concurrent int    `ini:"concurrent"`
	WorkPath   string `ini:"work_path"`
	Master     bool   `ini:"master_enable"`
	Worker     bool   `ini:"worker_enable"`
	UAPI       bool   `ini:"uapi_enable"`
	Timeout    int64  `ini:"node_timeout"`
}

// Mysql mysql config
type Mysql struct {
	MysqlURL      string `ini:"mysql_url"`
	MysqlName     string `ini:"mysql_name"`
	MysqlUser     string `ini:"mysql_user"`
	MysqlPassword string `ini:"mysql_password"`
}

// LocalStorage local storage for repo
type LocalStorage struct {
	Enable bool   `ini:"enable"`
	Path   string `ini:"storage_dir"`
}

// S3 s3 storage for repo
type S3 struct {
	Enable     bool   `ini:"enable"`
	S3URL      string `ini:"s3_endpoint"`
	S3Key      string `ini:"s3_key"`
	S3Secret   string `ini:"s3_secret"`
	BucketName string `ini:"bucket_name"`
}

// Git git storage for repo
type Git struct {
	Enable bool `ini:"enable"`
}

// Etcd etcd config
type Etcd struct {
	Enable    bool     `ini:"enable"`
	Endpoints []string `ini:"endpoints"`
}

// FileLog output log to file
type FileLog struct {
	Enable bool   `ini:"enable"`
	Path   string `ini:"log_dir"`
}

// Cfg defualt config
var Cfg = &Config{}

// NewConfig new config by config file
func NewConfig(file string) error {
	err := LoadConfig(file, Cfg)
	return err
}

// SetConfig set Cfg by Config struct
func SetConfig(c *Config) {
	Cfg = c
}

// LoadConfig load config from config file
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

// WriteConfig write the config info to file
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
