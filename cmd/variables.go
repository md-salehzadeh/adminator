package cmd

const CONFIG_NAME string = "adminator"

type Package struct {
	Name string
	Type string
	Args string
}

type File struct {
	Type   string
	Source string
	Path   string
}

type Service struct {
	Type    string
	Name    string
	Enabled bool
}
