package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var Config config

type config struct {
	Ip                string    `json:"ip" yaml:"ip"`
	Port              string    `json:"port" yaml:"port"`
	AdminTest         Configure `json:"admin_test" yaml:"admin_test"`
	AdminRelease      Configure `json:"admin_release" yaml:"admin_release"`
	EnterpriseTest    Configure `json:"im_enterprise_test" yaml:"im_enterprise_test"`
	EnterpriseRelease Configure `json:"im_enterprise_release" yaml:"im_enterprise_release"`
	ServerTest        Configure `json:"im_server_test" yaml:"im_server_test"`
	ServerRelease     Configure `json:"im_server_release" yaml:"im_server_release"`
	Accounts          []Account `json:"accounts" yaml:"accounts"`
	Sessions          map[string]string
}

type Configure struct {
	ProjectConfig `json:"project_config" yaml:"project_config"`
	ClientConfig  `json:"client_config" yaml:"client_config"`
	GitConfig     `json:"git_config" yaml:"git_config"`
	BuildConfigs  []BuildConfig `json:"build_configs" yaml:"build_configs"`
	ZipFilePath   string        `json:"zip_file_path" yaml:"zip_file_path"`
	ZipName       string        `json:"zip_name" yaml:"zip_name"`
	ServerPath    string        `json:"server_path" yaml:"server_path"`
}

type Account struct {
	Username string `json:"username" form:"username" yaml:"username"`
	Password string `json:"password" form:"password" yaml:"password"`
}

type BuildConfig struct {
	Env     string `json:"env" yaml:"env"`
	ModPath string `json:"mod_path" yaml:"mod_path"`
	BinName string `json:"bin_name" yaml:"bin_name"`
	Name    string `json:"name" yaml:"name"`
}

type ProjectConfig struct {
	ProjectPath string `json:"project_path" yaml:"project_path"`
	ProjectName string `json:"project_name" yaml:"project_name"`
	BinPath     string `json:"bin_path" yaml:"bin_path"`
	GitUrl      string `json:"git_url" yaml:"git_url"`
}

type ClientConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
}

type GitConfig struct {
	UserName string `json:"user_name" yaml:"user_name"`
	PassWord string `json:"pass_word" yaml:"pass_word"`
}

func Init(name string) {
	file, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(file, &Config); err != nil {
		panic(err)
	}
	Config.Sessions = make(map[string]string)
}
