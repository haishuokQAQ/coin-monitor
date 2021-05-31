package util

import (
    `github.com/mitchellh/mapstructure`
)

func DecodeFromMap(data interface{}, targetObject interface{}) error {
    decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
        Result:  targetObject,
        TagName: "json",
    })
    if err != nil {
        return err
    }
    return decoder.Decode(data)
}

func EncodeToMap(data interface{}, targetMap map[string]interface{}) error {
    encoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
        Result:           &targetMap,
        TagName:          "json",
    })
    if err != nil {
        return err
    }
    return encoder.Decode(data)
}
