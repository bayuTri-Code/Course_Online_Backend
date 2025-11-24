package otpemail

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"course_online_backend/internal/config"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

type OTPService struct {
	Redis        *redis.Client
	Config       *config.Config
	SMTPConfig   *config.SMTPConfig
	OTPExpiryMin int
	OTPLength    int
}

func NewOTPService(rdb *redis.Client, cfg *config.Config, smtpCfg *config.SMTPConfig) *OTPService {
	expiry := cfg.OTPExpiryMinutes
	if expiry == 0 {
		expiry = 5
	}
	length := cfg.OTPLength
	if length == 0 {
		length = 6
	}
	return &OTPService{
		Redis:        rdb,
		Config:       cfg,
		SMTPConfig:   smtpCfg,
		OTPExpiryMin: expiry,
		OTPLength:    length,
	}
}

func (s *OTPService) GenerateOTP() (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(s.OTPLength)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	format := fmt.Sprintf("%%0%dd", s.OTPLength)
	return fmt.Sprintf(format, n.Int64()), nil
}

func (s *OTPService) StoreOTP(ctx context.Context, email, otp string) error {
	key := fmt.Sprintf("otp:%s", email)
	exp := time.Duration(s.OTPExpiryMin) * time.Minute
	return s.Redis.Set(ctx, key, otp, exp).Err()
}

func (s *OTPService) GetOTP(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)
	return s.Redis.Get(ctx, key).Result()
}

func (s *OTPService) DeleteOTP(ctx context.Context, email string) {
	key := fmt.Sprintf("otp:%s", email)
	_ = s.Redis.Del(ctx, key).Err()
}

func (s *OTPService) CheckRateLimit(ctx context.Context, email string) error {
	key := fmt.Sprintf("otp_rate:%s", email)
	maxAttempts := 5
	window := time.Hour
	if s.Config.ServerEnv == "development" {
		maxAttempts = 100
		window = time.Minute
	}
	val, err := s.Redis.Get(ctx, key).Result()
	if err == redis.Nil {
		_ = s.Redis.Set(ctx, key, "1", window).Err()
		return nil
	}
	if err != nil {
		return err
	}
	var count int
	fmt.Sscanf(val, "%d", &count)
	if count >= maxAttempts {
		return fmt.Errorf("too many OTP requests, please try again later")
	}
	_, _ = s.Redis.Incr(ctx, key).Result()
	_, _ = s.Redis.Expire(ctx, key, window).Result()
	return nil
}

func (s *OTPService) SendOTPEmail(ctx context.Context, toEmail, otp string) error {
	if s.SMTPConfig == nil {
		return fmt.Errorf("smtp config not loaded")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.SMTPConfig.FromName, s.SMTPConfig.FromEmail))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your OTP Code")
	html := fmt.Sprintf(`
		<html>
		<body>
		<p>Your OTP code: <strong>%s</strong></p>
		<p>This code will expire in %d minutes.</p>
		</body>
		</html>
	`, otp, s.OTPExpiryMin)
	m.SetBody("text/html", html)
	d := gomail.NewDialer(
		s.SMTPConfig.Host,
		s.SMTPConfig.Port,
		s.SMTPConfig.User,
		s.SMTPConfig.Password,
	)
	d.TLSConfig = nil
	if err := d.DialAndSend(m); err != nil {
		log.Printf("[EMAIL] send error: %v", err)
		return err
	}
	return nil
}

func (s *OTPService) SendOTP(ctx context.Context, email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if err := s.CheckRateLimit(ctx, email); err != nil {
		return err
	}
	otp, err := s.GenerateOTP()
	if err != nil {
		return err
	}
	if err := s.StoreOTP(ctx, email, otp); err != nil {
		return err
	}
	if err := s.SendOTPEmail(ctx, email, otp); err != nil {
		s.DeleteOTP(ctx, email)
		return err
	}
	return nil
}

func (s *OTPService) VerifyOTP(ctx context.Context, email, otp string) error {
	stored, err := s.GetOTP(ctx, email)
	if err != nil {
		return fmt.Errorf("otp expired or not found")
	}
	if stored != otp {
		return fmt.Errorf("invalid otp")
	}
	_ = s.Redis.Del(ctx, fmt.Sprintf("otp_rate:%s", email)).Err()
	_ = s.Redis.Del(ctx, fmt.Sprintf("otp:%s", email)).Err()
	return nil
}
