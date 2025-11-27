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
	m.SetHeader("Subject", "Your Verification Code")

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OTP Verification</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f5;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse; background-color: #f5f5f5;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" style="max-width: 600px; width: 100%%; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600; letter-spacing: -0.5px;">
                                🔐 Verification Code
                            </h1>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 24px; color: #333333; font-size: 16px; line-height: 1.6;">
                                Hello,
                            </p>
                            <p style="margin: 0 0 32px; color: #666666; font-size: 16px; line-height: 1.6;">
                                Use the verification code below to complete your login:
                            </p>

                            <table role="presentation" style="width: 100%%; border-collapse: collapse; margin: 0 0 32px;">
                                <tr>
                                    <td align="center" style="padding: 30px; background-color: #f8f9fa; border: 2px dashed #e0e0e0; border-radius: 8px;">
                                        <div style="font-size: 36px; font-weight: 700; letter-spacing: 8px; color: #667eea; font-family: 'Courier New', monospace;">
                                            %s
                                        </div>
                                    </td>
                                </tr>
                            </table>

                            <table role="presentation" style="width: 100%%; border-collapse: collapse; margin: 0 0 24px;">
                                <tr>
                                    <td style="padding: 20px; background-color: #fff3cd; border-left: 4px solid #ffc107; border-radius: 4px;">
                                        <p style="margin: 0; color: #856404; font-size: 14px; line-height: 1.5;">
                                            ⏱️ <strong>This code will expire in %d minutes.</strong>
                                        </p>
                                    </td>
                                </tr>
                            </table>

                            <p style="margin: 0 0 16px; color: #666666; font-size: 14px; line-height: 1.6;">
                                If you didn't request this code, ignore this email.
                            </p>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 30px 40px; background-color: #f8f9fa; border-radius: 0 0 12px 12px; text-align: center; border-top: 1px solid #e0e0e0;">
                            <p style="margin: 0 0 8px; color: #999999; font-size: 13px;">
                                © 2025 %s. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
	`, otp, s.OTPExpiryMin, s.SMTPConfig.FromName)

	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		s.SMTPConfig.Host,
		s.SMTPConfig.Port,
		s.SMTPConfig.User,
		s.SMTPConfig.Password,
	)

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
