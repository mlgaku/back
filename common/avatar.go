package common

// 头像文件名
func AvatarFile(name string) string {
	return "avatar/" + name
}

// 头像URL连接
func AvatarURL(name string, url string) string {
	return url + AvatarFile(name)
}
