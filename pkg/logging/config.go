package logging

type Config struct {
	FilePath string `koanf:"file_path"`
	Encoding string `koanf:"encoding"`
	Level    string `koanf:"level"`
	Logger   string `koanf:"logger"`
}
