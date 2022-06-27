package main

type searchEngine interface {
	sendData(d interface{})
	closeInputChannel()
	outputChannel() <-chan searchEngineResult
}

type searchEngineResult struct {
	err  error
	data map[string]interface{}
}
