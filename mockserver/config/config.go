package config

var (
	// ServerAddress configures the server prefix in url generations.
	// 一般来说，我们网页中对其它网页的链接只需使用相对路径即可。
	// 但目标网站的所有链接都是用了绝对路径，为了模拟，我们也需要生成绝对路径。
	// 所以增加ServerAddress配置，所有的链接都使用形式：
	// http://<ServerAddress>/mock/album.zhenai.com/<相对路径>
	// 若将服务器部署在云，我们需要把这里替换成外网ip/域名:8080
	ServerAddress = "localhost:8080"

	// ListenAddress configures where the server listens at.
	ListenAddress = ":8080"
)
