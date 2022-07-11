package wechat

type Service struct{}

// Resolve resolves authorization code to wechat open id.
// 伪代码 测试使用
func (s *Service) Resolve(code string) string {
	resp := "this_is_openid" + code
	return resp
}
