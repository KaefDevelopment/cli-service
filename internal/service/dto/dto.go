package dto

import "github.com/jaroslav1991/cli-service/internal/model"

type Event struct {
	Id             string       `json:"id" gorm:"column:id"`
	CreatedAt      string       `json:"createdAt" gorm:"column:createdAt"`
	Type           string       `json:"type" gorm:"column:type"`
	Project        string       `json:"project,omitempty" gorm:"column:project"`
	ProjectBaseDir string       `json:"projectBaseDir,omitempty" gorm:"column:projectBaseDir"`
	Language       string       `json:"language,omitempty" gorm:"column:language"`
	Target         string       `json:"target,omitempty" gorm:"column:target"`
	Branch         string       `json:"branch,omitempty" gorm:"column:branch"`
	Timezone       string       `json:"timezone,omitempty" gorm:"column:timezone"`
	Params         model.Params `json:"params,omitempty" gorm:"column:params"`
	PluginId       string       `json:"pluginId" gorm:"column:pluginId"`
}

type SendEvents struct {
	OsName     string  `json:"osName"`
	DeviceName string  `json:"deviceName"`
	CliVersion string  `json:"cliVersion"`
	Events     []Event `json:"events"`
}
