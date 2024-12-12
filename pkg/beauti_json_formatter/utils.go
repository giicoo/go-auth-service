package beauti_json_formatter

import (
	"encoding/json"
	"log"
)

func marshal(o interface{}) string {
	b, err := json.MarshalIndent(o, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
