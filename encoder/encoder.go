package encoder

import (
	"encoding/base64"
	"encoding/json"
)

func Encode(obj interface{}) (string, error) {
	state, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.WithPadding(base64.StdPadding).EncodeToString(state), nil
}

func Decode(encState string, obj interface{}) error {

	decodedState, err := base64.StdEncoding.WithPadding(base64.StdPadding).DecodeString(encState)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decodedState, obj)
	return err
}
