package nagflux

import (
	"fmt"
	"github.com/griesbacher/nagflux/helper"
)

//Printable converts from nagfluxfile format to X
type Printable struct {
	Table     string
	Timestamp string
	Value     string
	tags      map[string]string
	fields    map[string]string
}

//PrintForInfluxDB prints the data in influxdb lineformat
func (p Printable) PrintForInfluxDB(version float32) string {
	line := helper.SanitizeInfluxInput(p.Table)
	p.tags = helper.SanitizeMap(p.tags)
	if len(p.tags) > 0 {
		line += fmt.Sprintf(`,%s`, helper.PrintMapAsString(helper.SanitizeMap(p.tags), ",", "="))
	}
	p.fields = helper.SanitizeMap(p.fields)
	line += fmt.Sprintf(` value=%s`, p.Value)
	if len(p.fields) > 0 {
		line += fmt.Sprintf(`,%s`, helper.PrintMapAsString(helper.SanitizeMap(p.fields), ",", "="))
	}
	return fmt.Sprintf("%s %s", line, p.Timestamp)
}

//PrintForElasticsearch prints in the elasticsearch json format
func (p Printable) PrintForElasticsearch(version float32, index string) string {
	if version >= 2 {
		head := fmt.Sprintf(`{"index":{"_index":"%s","_type":"%s"}}`, index, p.Table) + "\n"
		data := fmt.Sprintf(`{"timestamp":%s,"value":%s`, p.Timestamp, helper.GenJSONValueString(p.Value))
		data += helper.CreateJSONFromStringMap(p.tags)
		data += helper.CreateJSONFromStringMap(p.fields)
		data += "}\n"
		return head + data
	}
	return ""
}
