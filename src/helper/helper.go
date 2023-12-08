package helper

import "encoding/json"

func AddRandomAttribute(data []byte) ([]byte, error) {
	var jsonData interface{}

	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	if jsonObject, ok := jsonData.(map[string]interface{}); ok {
		jsonObject["foo"] = "bar"

	} else if jsonArray, ok := jsonData.([]interface{}); ok {
		for _, item := range jsonArray {
			if obj, ok := item.(map[string]interface{}); ok {
				obj["foo"] = "bar"
			}
		}
	}

	return json.Marshal(jsonData)
}
