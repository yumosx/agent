package domain

type Plan struct {
	Id    string
	Title string
	// 表示当前任务的执行情况
	Steps []Step
}

type Step struct {
	State   string
	Content string
}
