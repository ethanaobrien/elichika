package table

import (
	"reflect"
)

var (
	handlers    map[uintptr](func())
	prequisites map[uintptr][]uintptr
	updateOrder [](func())

	git    string
	output string
)

func addHandler(f func()) {
	if handlers == nil {
		handlers = make(map[uintptr](func()))
		prequisites = make(map[uintptr][]uintptr)
	}
	handlers[reflect.ValueOf(f).Pointer()] = f
}

func addPrequisite(function, prequisite func()) {
	addHandler(function)
	addHandler(prequisite)
	prequisites[reflect.ValueOf(function).Pointer()] = append(prequisites[reflect.ValueOf(function).Pointer()],
		reflect.ValueOf(prequisite).Pointer())
}

func generateUpdateOrder(fid uintptr) {
	_, exist := handlers[fid]
	if !exist {
		return // done
	}
	for _, prequisite := range prequisites[fid] {
		generateUpdateOrder(prequisite)
	}

	updateOrder = append(updateOrder, handlers[fid])
	delete(handlers, fid)
}

func Run(gitPlace string) string {
	git = gitPlace
	output = ""

	for len(handlers) > 0 {
		var fid uintptr
		for key := range handlers {
			fid = key
			break
		}
		generateUpdateOrder(fid)
	}
	for _, f := range updateOrder {
		f()
	}
	return output
}
