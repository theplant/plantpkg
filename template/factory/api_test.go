package factory_test

import (
	"testing"

	"github.com/theplant/testingutils"

	"github.com/theplant/plantpkg/template"
	"github.com/theplant/plantpkg/template/factory"
)

var sayHelloCases = []struct {
	name     string
	input    *template.SayHelloInput
	expected *template.SayHelloResult
}{
	{
		name: "case 1",
		input: &template.SayHelloInput{
			Hello: "hello1",
		},
		expected: &template.SayHelloResult{
			Result: "result",
		},
	},
	{
		name: "case 2",
		input: &template.SayHelloInput{
			Hello: "hello2",
		},
		expected: &template.SayHelloResult{
			Result: "result",
		},
	},
}

func TestSayHello(t *testing.T) {
	serv := factory.New(nil, "2").OptionalParam1("param1").OptionalParam2("param2")
	for _, c := range sayHelloCases {
		result, err := serv.SayHello(c.input)
		if err != nil {
			panic(err)
		}

		diff := testingutils.PrettyJsonDiff(c.expected, result)
		if len(diff) > 0 {
			t.Error(c.name, diff)
		}
	}
}
