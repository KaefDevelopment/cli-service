package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

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
	Params         Params `json:"params,omitempty" gorm:"column:params"`
	AuthKey        string `json:"authKey" gorm:"column:authKey"`
	Send           bool   `json:"send" gorm:"column:send"`
}

type Events struct {
	Events []Event `json:"events"`
}

type Response struct {
	Events []Events `json:"events"`
}

type Params map[string]interface{}

func (m Params) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := m.MarshalJSON()
	return string(ba), err
}

func (m *Params) Scan(val interface{}) error {
	if val == nil {
		*m = make(Params)
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	t := map[string]interface{}{}
	rd := bytes.NewReader(ba)
	decoder := json.NewDecoder(rd)
	decoder.UseNumber()
	err := decoder.Decode(&t)
	*m = t
	return err
}

func (m Params) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	t := (map[string]interface{})(m)
	return json.Marshal(t)
}

func (m *Params) UnmarshalJSON(b []byte) error {
	t := map[string]interface{}{}
	err := json.Unmarshal(b, &t)
	*m = Params(t)
	return err
}

func (m Params) GormDataType() string {
	return "jsonmap"
}

func (Params) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	}

	return ""
}
