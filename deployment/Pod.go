package deployment

type Pod struct {
	Name       string
	NameSpace  string
	Images     string
	NodeName   string
	IsReady    string
	Message    string
	CreateTime string
	IPs        []string //第一个是pod ip,第二个是host ip
}
