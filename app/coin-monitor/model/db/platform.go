package db

type Platform struct {
	Id           uint64 `json:"id"`
	PlatformName string `json:"platform_name"`
	PlatformUrl  string `json:"platform_url"`
}
