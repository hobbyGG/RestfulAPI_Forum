package models

type ParamSignUp struct {
	UserName string `json:"username" binding:"required"`
	Pwd      string `json:"pwd" binding:"required"`
	RePwd    string `json:"re_pwd" binding:"required,eqfield=Pwd"`
	Email    string `json:"email" binding:"required"`
}

type ParamLogin struct {
	UserName string `json:"username" binding:"required"`
	Pwd      string `json:"pwd" binding:"required"`
}

type ParamCreatePost struct {
	CommID  int16  `json:"commID" binding:"omitempty,min=1,max=5" db:"commID"` //作为可选项，不加就默认为通用区
	Title   string `json:"title" binding:"required" db:"title"`
	Content string `json:"content" binding:"required" db:"content"`
}
