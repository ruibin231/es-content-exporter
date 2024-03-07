package settings

import (
	"es-content-export/utils/async"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const common = "settings"

var Config config
var ESMap = make(map[string]*EsClient)

type set struct {
	viper *viper.Viper
}

func (s *set) unmarshalKey() {
	if err := s.viper.UnmarshalKey(common, &Config); err != nil {
		panic(err)
	}
	for _, c := range Config.EsList {
		ESMap[c.Alise] = c
	}
}

// watch - 监听yml文件修改后，重新解析 没啥意义
func (s *set) watch() {
	async.Go(func() {
		s.viper.WatchConfig()
		s.viper.OnConfigChange(func(in fsnotify.Event) {
			s.unmarshalKey()
		})
	})
}

func SetupConfig(configFile string) {
	s := set{
		viper: viper.New(),
	}
	s.viper.SetConfigFile(configFile)
	s.viper.SetConfigType("yml")
	if err := s.viper.ReadInConfig(); err != nil {
		panic(err)
	}
	s.unmarshalKey()
	s.watch()
}
