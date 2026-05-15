package httpgateway

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type logoutRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type createMovieRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Genre       string  `json:"genre" binding:"required"`
	Duration    int32   `json:"duration" binding:"required,gt=0"`
	Rating      float64 `json:"rating" binding:"gte=0,lte=10"`
}

type createHallRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int32  `json:"capacity" binding:"required,gt=0"`
}

type createSessionRequest struct {
	MovieID   string  `json:"movie_id" binding:"required"`
	HallID    string  `json:"hall_id" binding:"required"`
	StartTime string  `json:"start_time" binding:"required"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

type createBookingRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	SeatID    string `json:"seat_id" binding:"required"`
}

type confirmPaymentRequest struct {
	BookingID string  `json:"booking_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}
