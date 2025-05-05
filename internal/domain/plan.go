package domain

// Plan 大模型给我们返回的
type Plan struct {
	Id    string
	Title string
	Steps []Step
}

type Step struct {
	State   string
	Content string
}
