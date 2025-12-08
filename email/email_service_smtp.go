package service

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
)

type smtpEmailService struct {
	from     string
	password string
	host     string
	port     string
}

func NewSMTPEmailService() EmailService {
	return &smtpEmailService{
		from:     os.Getenv("SMTP_EMAIL"),
		password: os.Getenv("SMTP_PASSWORD"),
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
	}
}

// Method Send yang lama (tetap ada)
func (s *smtpEmailService) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.from, s.password, s.host)

	msg := []byte(
		"From: Arshaka Bimantara <" + s.from + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n\r\n" +
			body,
	)

	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}


// Method SendOTP yang baru (dengan template HTML keren)
func (s *smtpEmailService) SendOTP(to, userName, otp string) error {
	subject := "Kode OTP Reset Password - Arshaka"
	currentDate := time.Now().Format("02/01/2006")

	// GANTI BASE64 DENGAN URL LANGSUNG
	bgURL := "https://imgur.com/VcjpbRw.png" // Upload ke imgur/cloudinary
	logoURL := "https://imgur.com/gicQliA.png"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin: 0; padding: 0; font-family: 'Poppins', Arial, sans-serif; background-color: #f4f4f4;">
    <table width="100%%" cellpadding="0" cellspacing="0" style="background-color: #f4f4f4; padding: 20px;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 20px; overflow: hidden;">
                    
                    <!-- Header dengan Background -->
                    <tr>
                        <td style="background-image: url('%s'); background-size: cover; background-position: center; padding: 30px 40px; height: 200px;">
                            <table width="100%%">
                                <tr>
                                    <td style="padding: 0;">
                                        <img src="%s" alt="Logo Mapala" style="height: 80px; width: auto;" />
                                    </td>
                                    <td align="right">
                                        <p style="color: #ffffff; font-size: 14px; margin: 0;">tanggal: %s</p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <table width="100%%" cellpadding="0" cellspacing="0" style="background-color: #f8f8f8; border-radius: 20px; padding: 50px 40px;">
                                <tr>
                                    <td align="center">
                                        <h1 style="font-size: 48px; font-weight: 800; color: #1a1a1a; margin: 0 0 30px 0;">YOUR OTP</h1>
                                        
                                        <p style="font-size: 18px; color: #333333; margin: 0 0 40px 0;">Hey, %s</p>
                                        
                                        <p style="font-size: 16px; color: #666666; margin: 0 0 40px 0;">
                                            Berikut adalah kode OTP untuk melakukan reset<br>password akun Anda:
                                        </p>
                                        
                                        <!-- OTP -->
                                        <div style="margin: 40px 0;">
                                            <span style="font-size: 72px; font-weight: 800; color: #DC143C; letter-spacing: 20px;">%s</span>
                                        </div>
                                        
                                        <p style="font-size: 18px; color: #333333; margin: 40px 0 20px 0;">Masa berlaku: 10 menit</p>
                                        
                                        <p style="font-size: 14px; color: #666666; margin: 0;">
                                            Jika Anda tidak merasa meminta reset password, abaikan pesan ini.<br>
                                            Jangan bagikan kode ini kepada siapa pun demi keamanan akun Anda.
                                        </p>
                                        
                                        <p style="font-size: 18px; color: #333333; margin: 40px 0 0 0;">Terima kasih.</p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #f8f9fa; padding: 30px 40px; text-align: center;">
                            <p style="color: #999999; font-size: 12px; margin: 0 0 5px 0;">
                                Email ini dikirim otomatis, mohon tidak membalas.
                            </p>
                            <p style="color: #999999; font-size: 12px; margin: 0;">
                                &copy; 2024 Mapala UIN Syarif Hidayatullah Jakarta
                            </p>
                        </td>
                    </tr>
                    
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
	`, bgURL, logoURL, currentDate, userName, otp)

	return s.Send(to, subject, htmlBody)
}