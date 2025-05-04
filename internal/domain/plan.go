package domain

// Plan 大模型给我们返回的
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
