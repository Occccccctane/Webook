package Handler

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Date any    `json:"date"`
}
