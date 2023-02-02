package config

// config yaml.Unmarshal config.yml
type config struct {
	Server  *server  `yaml:"server"`
	Log     *log     `yaml:"log"`
	Babel   *babel   `yaml:"babel"`
	VmPool  *vmpool  `yaml:"vmpool"`
	Astpool *astpool `yaml:"astpool"`
	Db      *db      `yaml:"db"`
}

type server struct {
	Debug   bool   `yaml:"debug"`
	Port    int    `yaml:"port"`
	Dir     string `yaml:"dir"`
	Timeout int    `yaml:"timeout"`
	Enable  struct {
		Dir  bool `yaml:"dir"`
		File bool `yaml:"file"`
		Jssp bool `yaml:"jssp"`
		Jsjs bool `yaml:"jsjs"`
	} `yaml:"enable"`
}

type log struct {
	Info   string `yaml:"info"`
	Error  string `yaml:"error"`
	Access string `yaml:"access"`
}

type babel struct {
	Enable bool   `yaml:"enable"`
	Path   string `yaml:"path"`
	Ts     bool   `yaml:"ts"`
}

type vmpool struct {
	Enable bool `yaml:"enable"`
	Size   int  `yaml:"size"`
	Retry  int  `yaml:"retry"`
}

type astpool struct {
	Enable bool   `yaml:"enable"`
	Size   int    `yaml:"size"`
	Mode   string `yaml:"mode"`
}

type db struct {
	Enable bool `yaml:"enable"`
}
