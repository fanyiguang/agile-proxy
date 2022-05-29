package config

var (
	getIpUrl = []string{
		"https://ipinfo.io/json",
		"http://lumtest.com/myip.json",
		"https://checkip.amazonaws.com",
		"http://myipip.net/",
	}
)

func GetIpUrls() []string {
	return getIpUrl
}
