package uciadapter

type EngineInfo struct {
	Name    string
	Author  string
	Options []OptionInfo
}

type OptionInfo struct {
	Name string
}

type Option struct {
	Name  string
	Value string
}
