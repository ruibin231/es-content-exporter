package settings

type esClient struct {
	Host string
	Port string
	User string
	Pass string
}

type searchData struct {
	IndexPrefix string `comment:"索引前缀"` // 会以 xxx*方式搜索
	Field       string
	Must        []string
	MustNot     []string
	Cycle       int
}

type config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	IndexPrefix string `comment:"索引前缀"` // 会以 xxx*方式搜索
	Field       string
	Content     string
	Must        []string
	MustNot     []string
	Cycle       int
}
