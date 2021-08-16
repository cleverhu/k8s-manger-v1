package deployment

type Deployment struct {
	Name       string
	NameSpace  string
	Replicas   [3]int32
	Images     string
	CreateTime string
	Pods       []*Pod
	IsComplete bool
	Message    string
}
