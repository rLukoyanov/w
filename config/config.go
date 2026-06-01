package config

type Config struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	DataDir string `yaml:"data_dir"`
	DBPath  string `yaml:"db_path"`
}
