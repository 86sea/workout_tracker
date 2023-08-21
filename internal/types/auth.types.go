package types

type LoginDTO struct {
    Email string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"password"`
}

type UserResponse struct {
    ID uint `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
    Password string `json:"Password"`
}

type AccessResponse struct {
    Token string `json:"token"`
}

type AuthResponse struct {
    User *UserResponse `json:"user"`
    Auth *AccessResponse `json:"auth"`
}
