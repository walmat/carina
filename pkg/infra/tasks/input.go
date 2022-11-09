package tasks

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const inputTag = "input"

var (
	InvalidInputTypeErr = errors.New("invalid input type")
)

type InputData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type TextInput struct {
	Name        string `json:"name"`
	Placeholder string `json:"placeholder"`
}

type NumericInput struct {
	Name         string `json:"name"`
	InitialValue int    `json:"initialValue"`
	Minimum      int    `json:"min"`
	Maximum      int    `json:"max"`
}

type SelectInput struct {
	Name         string   `json:"name"`
	Values       []string `json:"values"`
	InitialValue int      `json:"initialValue"`
}

func getInputDataFromFieldTag(tag string) (InputData, error) {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "text":
		textInput := TextInput{}
		_, err := fmt.Sscanf(strings.Join(args[1:], ","), "name: %s placeholder: %s", &textInput.Name, &textInput.Placeholder)
		if err != nil {
			return InputData{}, err
		}
		return InputData{
			Type: "text",
			Data: textInput,
		}, nil
	case "numeric":
		numericInput := NumericInput{}
		_, err := fmt.Sscanf(strings.Join(args[1:], ","), "name: %s initial: %d min: %d max: %d", &numericInput.Name, &numericInput.InitialValue, &numericInput.Minimum, &numericInput.Maximum)
		if err != nil {
			return InputData{}, err
		}
		return InputData{
			Type: "numeric",
			Data: numericInput,
		}, nil
	case "select":
		selectInput := SelectInput{}
		rawValues := ""
		_, err := fmt.Sscanf(strings.Join(args[1:], ","), "name: %s values: %s initial: %d", &selectInput.Name, &rawValues, &selectInput.InitialValue)
		selectInput.Values = strings.Split(rawValues, ",")
		if err != nil {
			return InputData{}, err
		}
		return InputData{
			Type: "select",
			Data: selectInput,
		}, nil
	}

	return InputData{}, InvalidInputTypeErr
}

func getInputDataFromType(ctxType reflect.Type) ([]InputData, error) {
	var arrInputData []InputData
	for i := 0; i < ctxType.NumField(); i++ {
		tag := ctxType.Field(i).Tag.Get(inputTag)

		if tag == "" || tag == "-" {
			continue
		}

		inputData, err := getInputDataFromFieldTag(tag)
		if err != nil {
			return nil, err
		}
		arrInputData = append(arrInputData, inputData)
	}
	return arrInputData, nil
}

func GetInputDataForTaskType(typeId string, modeId string) ([]InputData, error) {
	handler, err := getTaskHandler(typeId, modeId)
	if err != nil {
		return nil, err
	}
	return getInputDataFromType(handler.staticCtx.Elem())
}
