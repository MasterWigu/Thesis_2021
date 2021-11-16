package main


var tool *ToolHandler

type ToolHandler struct {

}


func createToolHandler() {
	tool = new(ToolHandler)
}

func (t *ToolHandler) ExecAction(actionId string) (string, error) {

	return "", nil
}

func (t *ToolHandler) DryRunAction(actionId string) (string, error) {

	return "", nil
}
