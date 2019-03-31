package run

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"
)

func Run(file, src string) (string, error) {

	out := ""
	state := 0
	script := ""
	for _, line := range strings.Split(string(src), "\n") {
		switch state {
		case 0:
			if line == "``` javascript" {
				state = 1
			} else {
				out += line + "\n"
			}
		case 1:
			if line == "```" {
				state = 0
				result, err := runScript(file, script)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					out += err.Error() + "\n"
				} else {
					out += result + "\n"
				}

				script = ""
			} else {
				script += line + "\n"
			}
		}
	}
	return out[:len(out)-1], nil
}

func RunFile(file string) (string, error) {

	src, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	file, err = filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return Run(file, string(src))
}

func runScript(file, src string) (string, error) {

	vm := startVM(file)
	val, err := vm.Run(fmt.Sprintf("(function(){%s})()", src))
	if err != nil {
		return "", err
	}
	if val.IsUndefined() {
		return "", nil
	}
	return val.String(), nil
}

func throw(vm *otto.Otto, err error) otto.Value {
	vm.Eval(`throw new Error("test")`)
	return otto.Value{}
}

func startVM(file string) *otto.Otto {

	vm := otto.New()

	vm.Set("readFile", func(call otto.FunctionCall) otto.Value {
		inFile, err := call.Argument(0).ToString()
		if err != nil {
			return throw(call.Otto, err)
		}

		out, err := ioutil.ReadFile(inFile)
		if err != nil {
			return throw(call.Otto, err)
		}

		val, err := call.Otto.ToValue(string(out))
		if err != nil {
			return throw(call.Otto, err)
		}
		return val
	})
	vm.Set("run", func(call otto.FunctionCall) otto.Value {
		file := filepath.Join(filepath.Dir(file), call.Argument(0).String())

		out, err := RunFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return otto.Value{}
		}
		val, _ := call.Otto.ToValue(out)
		return val
	})
	vm.Set("$", func(call otto.FunctionCall) otto.Value {
		section := call.Argument(0).String()

		out, err := ExtractFile(file, section)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return otto.Value{}
		}
		val, _ := call.Otto.ToValue(out)
		return val
	})
	vm.Set("extract", func(call otto.FunctionCall) otto.Value {
		inFile := filepath.Join(filepath.Dir(file), call.Argument(0).String())
		section := call.Argument(1).String()

		out, err := ExtractFile(inFile, section)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return otto.Value{}
		}
		val, _ := call.Otto.ToValue(out)
		return val
	})
	vm.Set("log", func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			fmt.Println(arg.String())
		}
		return otto.Value{}
	})

	vm.Set("__FILE__", file)
	vm.Set("__DIR__", filepath.Dir(file))

	vm.Run(`
	function require(file) {
		var src = readFile(file);
		var module = {};
		var exports = {};
		var __FILE__ = file
		var __DIR__ = file.split('/').slice(0. -1).join("/")
		eval(src);
		return exports || module.exports;
	}
	`)
	return vm
}
