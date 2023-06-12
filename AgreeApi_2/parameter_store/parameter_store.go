package parameter_store

import (
	"encoding/json"
	"fmt"
	"os"
)

type ParameterStoreIface interface {
	ReadFile() error
	ReadConfig() error
	KeyVals() KeyValFile
	Parameter(string) (map[string]interface{}, error)
	ParameterToBytes(string) ([]byte, error)
}

// {
// 	keys": [{
// 			"key":   "<some key in AWS Parameter Store>",
// 			"value": {
// 						"Name":  "<the key in AWS Parameter Store>",
// 						"Value": "<the value of the key in AWS Parameter Store"
// 					 }
// 		   }]
// }
type KeyValFile struct {
	KeyVals []struct {
		Key        string                 `json:"key"`
		Parameters map[string]interface{} `json:"value"`
	} `json:"keys"`
}

type KeyClient struct {
	KeyValFile
	Filename  string
	FileBytes []byte
	// paramBytes []byte
}

type MockKeyClient struct {
	KeyValFile
	FileBytes  []byte
	paramBytes []byte
}

func (kc *KeyClient) ReadFile() error {
	var err error
	kc.FileBytes, err = os.ReadFile(kc.Filename)
	if err != nil {
		return err
	}
	return nil
}

func (kc *KeyClient) ReadConfig() error {
	if kc.KeyValFile.KeyVals == nil {
		if err := kc.ReadFile(); err != nil {
			return err
		}
	}
	if err := json.Unmarshal(kc.FileBytes, &kc.KeyValFile); err != nil {
		return err
	}
	// var err error
	// kc.paramBytes, err = json.Marshal(kc.fileBytes)
	// if err != nil {
	// 	return err
	// }
	return nil
}
func (kc *KeyClient) KeyVals() KeyValFile { return kc.KeyValFile }

func (kc *KeyClient) Parameter(key string) (map[string]interface{}, error) {
	return getParameter(kc, key)
}

func (kc *KeyClient) ParameterToBytes(key string) ([]byte, error) {
	params, err := getParameter(kc, key)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (mkc *MockKeyClient) ReadFile() error { return nil }
func (mkc *MockKeyClient) ReadConfig() error {
	if err := json.Unmarshal(mkc.FileBytes, &mkc.KeyValFile); err != nil {
		return err
	}
	// var err error
	// mkc.paramBytes, err = json.Marshal(mkc.fileBytes)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (kc *MockKeyClient) KeyVals() KeyValFile { return kc.KeyValFile }

func (kc *MockKeyClient) Parameter(key string) (map[string]interface{}, error) {
	return getParameter(kc, key)
}
func (kc *MockKeyClient) ParameterToBytes(key string) ([]byte, error) {
	params, err := getParameter(kc, key)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func getParameter(secretsStore ParameterStoreIface, key string) (map[string]interface{}, error) {
	if secretsStore.KeyVals().KeyVals == nil {
		return nil, fmt.Errorf("Store empty")
	}

	for _, v := range secretsStore.KeyVals().KeyVals {
		if v.Key == key {
			return v.Parameters, nil
		}
	}

	return nil, fmt.Errorf("Key not found: %s", key)
}
