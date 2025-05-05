module github.com/mohdfareed/gptx-cli

go 1.23.8

require ( // MARK: Dependencies
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.6.0-pre.2
	github.com/openai/openai-go v0.1.0-beta.10
	github.com/urfave/cli/v3 v3.3.2
)

require ( // MARK: Configuration
	github.com/knadh/koanf/parsers/dotenv v1.1.0
	github.com/knadh/koanf/providers/env v1.1.0
	github.com/knadh/koanf/providers/file v1.2.0
	github.com/knadh/koanf/providers/structs v1.0.0
	github.com/knadh/koanf/v2 v2.2.0
)

require ( // MARK: Transitive
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
