module github.com/caeret/logging

go 1.19

require (
	go.uber.org/zap v1.23.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace go.uber.org/zap => github.com/caeret/zap v0.0.0-20220910091553-975ef1636a9a
