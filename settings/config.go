package settings

type EsClient struct {
	Alise string
	Host  string
	Port  string
	User  string
	Pass  string
}

type QueryData struct {
	ES          string
	IndexPrefix string `comment:"索引前缀"` // 会以 xxx*方式搜索
	Field       string
	Content     string
	Cycle       int
}

type config struct {
	EsList    []*EsClient
	QueryList []*QueryData
}

/*
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
*/
