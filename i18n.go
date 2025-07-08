package vortex

func init() {
	i18n = make(map[string]string)
}

// 初始化i18n配置
func initI18n(tmp map[string]string) {
	for k, v := range tmp {
		i18n[k] = v
	}
}

var i18n map[string]string

// 获取i18n
func getI18n(key string) string {
	return i18n[key]
}
