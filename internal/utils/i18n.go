package vortexu

// 初始化i18n配置
func InitI18n(tmp map[string]string) {
	for k, v := range tmp {
		I18n[k] = v
	}
}

var I18n map[string]string

// 获取i18n
func GetI18n(key string) string {
	return I18n[key]
}
