package dto

type Event struct {
	Id             string `json:"id" gorm:"column:id"`
	CreatedAt      string `json:"createdAt" gorm:"column:createdAt"`
	Type           string `json:"type" gorm:"column:type"`
	Project        string `json:"project,omitempty" gorm:"column:project"`
	ProjectBaseDir string `json:"projectBaseDir,omitempty" gorm:"column:projectBaseDir"`
	Language       string `json:"language,omitempty" gorm:"column:language"`
	Target         string `json:"target,omitempty" gorm:"column:target"`
	Branch         string `json:"branch,omitempty" gorm:"column:branch"`
	Timezone       string `json:"timezone,omitempty" gorm:"column:timezone"`
	Params         string `json:"params,omitempty" gorm:"column:params"`
}

type Events struct {
	Events []Event `json:"events"`
}

type Response struct {
	Events []Event `json:"events"`
}
