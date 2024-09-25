package services

import (
	"bookstore_api/internal/repositories"
	"bookstore_api/models"
	"bookstore_api/tools"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
	"unicode"
)

type UserService struct {
	*Service
	userRepo repositories.IUserRepository
}

func NewUserService(service *Service, userRepo repositories.IUserRepository) *UserService {
	return &UserService{
		Service:  service,
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, user *models.UserRegister) (*models.UserResponse, error) {
	err := s.validateEmail(user.Email)
	if err != nil {
		return nil, err
	}

	err = s.validatePassword(user.Password)
	if err != nil {
		return nil, err
	}

	newPass, err := s.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = newPass

	createdUser, err := s.userRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	userResponse := s.convertToResponse(createdUser)
	return userResponse, nil
}

func (s *UserService) LoginUser(ctx context.Context, user *models.UserLogin) (*models.UserResponse, error) {
	// Just in case of SQL injection or something
	err := s.validateEmail(user.Email)
	if err != nil {
		return nil, err
	}

	checkUser, err := s.userRepo.Get(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	err = s.checkPassword(user.Password, checkUser.Password)
	if err != nil {
		return nil, err
	}

	userResponse := s.convertToResponse(checkUser)
	return userResponse, nil
}

func (s *UserService) UpdateUserData(ctx context.Context, userData *models.UserUpdateData) (*models.UserResponse, error) {
	token := ctx.Value("token").(string)
	claims, err := tools.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	email := claims.Subject

	checkUser, err := s.userRepo.Get(ctx, email)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	checkUser.UpdatedAt = &t
	checkUser.Name = userData.Name

	updatedUser, err := s.userRepo.Update(ctx, checkUser)
	if err != nil {
		return nil, err
	}

	// Do I must re-authenticate? Invalidate session token? Maybe nah, don't have to

	userResponse := s.convertToResponse(updatedUser)
	return userResponse, nil
}

func (s *UserService) UpdateUserEmail(ctx context.Context, updatedEmail string) (*models.UserResponse, error) {
	token := ctx.Value("token").(string)
	claims, err := tools.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Will be nice if we send OTP when validating email
	err = s.validateEmail(updatedEmail)
	if err != nil {
		return nil, err
	}

	email := claims.Subject
	checkUser, err := s.userRepo.Get(ctx, email)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	checkUser.UpdatedAt = &t
	checkUser.Email = updatedEmail

	updatedUser, err := s.userRepo.Update(ctx, checkUser)
	if err != nil {
		return nil, err
	}

	userResponse := s.convertToResponse(updatedUser)

	return userResponse, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, updatedPassword string) (*models.UserResponse, error) {
	token := ctx.Value("token").(string)
	claims, err := tools.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	err = s.validatePassword(updatedPassword)
	if err != nil {
		return nil, err
	}

	email := claims.Subject
	checkUser, err := s.userRepo.Get(ctx, email)
	if err != nil {
		return nil, err
	}

	newPassword, err := s.hashPassword(updatedPassword)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	checkUser.UpdatedAt = &t
	checkUser.Password = newPassword

	updatedUser, err := s.userRepo.Update(ctx, checkUser)
	if err != nil {
		return nil, err
	}

	// Okay, maybe for this and email, I must invalidate sessions.
	// Hey, you know, instead of injecting sessions service into user service, I should just inject it to user handler

	userResponse := s.convertToResponse(updatedUser)

	return userResponse, nil
}

func (s *UserService) validateEmail(email string) error {
	// Define a regular expression for a valid email format
	var re = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Maybe elaborate further by sending an OTP
	return nil
}

func (s *UserService) validatePassword(password string) error {
	hasMinLen := false
	hasUpper := false
	hasLower := false
	hasNumber := false

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasMinLen {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}

func (s *UserService) convertToResponse(user *models.User) *models.UserResponse {
	userResponse := &models.UserResponse{}

	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	userResponse.IsAdmin = user.IsAdmin

	return userResponse
}

//func (s *UserService) convertToLoginResponse(user *models.User, accessToken string, accessClaims *tools.CustomClaims, refreshToken string, refreshClaims *tools.CustomClaims) *models.UserLoginResponse {
//	userResponse := &models.UserLoginResponse{}
//
//	userResponse.User.ID = user.ID
//	userResponse.User.Name = user.Name
//	userResponse.User.Email = user.Email
//	userResponse.User.IsAdmin = user.IsAdmin
//
//	userResponse.AccessToken = accessToken
//	userResponse.AccessTokenExpiresAt = accessClaims.ExpiresAt.Time
//
//	userResponse.RefreshToken = refreshToken
//	userResponse.RefreshTokenExpiresAt = refreshClaims.ExpiresAt.Time
//
//	return userResponse
//}

func (s *UserService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *UserService) checkPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Not worth it, would require me to do an extra query. Just hash the pass, and we're good. Only commenting, just in case.
/*func (s *UserService) encryptUser(user *models.UserRegister) error {
	encryptedName, err := tools.Encrypt([]byte(user.Name), s.AESKey)
	if err != nil {
		return err
	}

	encryptedEmail, err := tools.Encrypt([]byte(user.Email), s.AESKey)
	if err != nil {
		return err
	}

	user.Name = encryptedName
	user.Email = encryptedEmail

	return nil
}

func (s *UserService) decryptUser(user *models.UserResponse) error {
	decryptedName, err := tools.Decrypt(user.Name, s.AESKey)
	if err != nil {
		return err
	}

	decryptedEmail, err := tools.Decrypt(user.Email, s.AESKey)
	if err != nil {
		return err
	}

	user.Name = string(decryptedName)
	user.Email = string(decryptedEmail)

	return nil
}
*/
