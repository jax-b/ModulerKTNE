module github.com/jax-b/ModulerKTNE

go 1.23.0

toolchain go1.23.11

require (
	github.com/asticode/go-astikit v0.56.0
	github.com/asticode/go-astilectron v0.30.0
	github.com/d2r2/go-i2c v0.0.0-20191123181816-73a8a799d6bc
	github.com/faiface/beep v1.1.0
	github.com/james-barrow/golang-ipc v1.2.4
	github.com/jax-b/go-i2c7Seg v1.0.0
	github.com/spf13/viper v1.20.1
	github.com/stianeikeland/go-rpio/v4 v4.6.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.42.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/d2r2/go-logger v0.0.0-20210606094344-60e9d1233e22 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/hajimehoshi/oto v1.0.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20250711185948-6ae5c78190dc // indirect
	golang.org/x/exp/shiny v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/image v0.29.0 // indirect
	golang.org/x/mobile v0.0.0-20250711185624-d5bb5ecc55c0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	"github.com/jax-b/ModulerKTNE/WireTypes" => ./WireTypes
)