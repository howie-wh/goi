package model

type ImportConfig struct {
	Type       string `json:"type"`
	FileName   string `json:"file"`
	InjectCode string `json:"code"`
}

type FunctionConfig struct {
	Type         string `json:"type"`
	FileName     string `json:"file"`
	FuncName     string `json:"func"`
	RelativeLine int    `json:"relative_line"`
	InjectCode   string `json:"code"`
}

type InjectConfig struct {
	ImportConfigs   []*ImportConfig   `json:"import_configs"`
	FunctionConfigs []*FunctionConfig `json:"function_configs"`
}
