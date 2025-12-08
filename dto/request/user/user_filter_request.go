package user

type UserFilter struct {
	Search       string   `query:"q" form:"q"` // Tambahin form tag juga!
	Page         int      `query:"page" form:"page"`
	Limit        int      `query:"limit" form:"limit"`
	ExcludeRole  string   `query:"-"` // Ini dari usecase
	ExcludeRoles []string `query:"-"`
}