package params

func NewPlanParams() *Parameters {
	parameters := newParams()

	parameters.Params["command"] = NewValue(
		"string",
		"The command to execute. Available commands: create, update, list, get, set_active, mark_step, delete.",
		WithEnum([]string{
			"create",
			"update",
			"list",
			"get",
			"set_active",
			"mark_step",
			"delete",
		}))

	parameters.Params["plan_id"] = NewValue(
		"string",
		"Unique identifier for the plan. Required for create, update, set_active, and delete commands. Optional for get and mark_step (uses active plan if not specified).")

	parameters.Params["title"] = NewValue(
		"string",
		"Title for the plan. Required for create command, optional for update command.")

	parameters.Params["steps"] = NewValue(
		"array",
		"List of plan steps. Required for create command, optional for update command.",
		WithItem(map[string]string{
			"type": "string",
		}),
	)

	parameters.Params["step_index"] = NewValue(
		"integer",
		"Index of the step to update (0-based). Required for mark_step command.")

	parameters.Params["step_status"] = NewValue(
		"string",
		"Status to set for a step. Used with mark_step command.",
		WithEnum([]string{
			"not_started",
			"in_progress",
			"completed",
			"blocked",
		}))

	parameters.Params["step_notes"] = NewValue(
		"string",
		"Additional notes for a step. Optional for mark_step command.")

	return parameters
}

func NewTrimParams() *Parameters {
	parameters := newParams()
	parameters.Params["status"] = NewValue(
		"string",
		"the finish status of the interaction.",
		WithEnum([]string{"success", "failure"}),
	)
	return parameters
}

func NewGoParams() *Parameters {
	parameters := newParams()
	parameters.Params["code"] = NewValue(
		"string",
		"The Golang code to execute")
	return parameters
}

func NewChatParams() *Parameters {
	parameters := newParams()
	parameters.Params["response"] = NewValue(
		"string",
		"The response text that should be delivered to the user.")
	return parameters
}
