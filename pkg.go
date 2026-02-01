package sqlbuilder

import "encoding/json"

func toMap(data, dst any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, &dst)
	if err != nil {
		return err
	}

	return nil
}
