type Site struct {
	Name      string    `json:"name" validate:"required"`
	SiteURL   string    `json:"site_url" validate:"required,url"`
	ImageURL  string    `json:"image_url"`
	Language  string    `json:"language" validate:"required,oneof=fr en es"`
	Ads       int       `json:"ads" validate:"oneof=0 1"`
	Type      string    `json:"type" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
