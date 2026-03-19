package handler

//============用户认证请求结构体=============
type RegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInfo struct {
	Username string `json:"username"`
	Telenum  string `json:"telenum"`
	Gender   string `json:"gender" binding:"omitempty,oneof=男 女 草履虫"`
}

type UpdateEmail struct {
	NewEmail string `json:"newemail"`
	Password string `json:"password"`
}

type UpdatePassword struct {
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

type CancelUser struct {
	Password string `json:"password"`
}

//============待办事项请求结构体=============
type CreateTodoRequest struct {
	Event string `json:"event" binding:"required"`
}

// UpdateTodoRequest 更新待办事项请求
type UpdateTodoRequest struct {
	Event string `json:"event"`
}

//============学习数据请求结构体=============
type AddStudyDataRequest struct {
	StudyTime int `json:"studytime" binding:"required"`
	Tomatoes  int `json:"tomatoes" binding:"required"`
}

//============音乐请求结构体=============
type UploadMusicRequest struct {
	Author string `form:"author" binding:"required"`
	Title  string `form:"title" binding:"required"`
}

// ============环境音请求结构体=============
type CreateAmbientSoundRequest struct {
	Name string `form:"name" binding:"required"`
}
