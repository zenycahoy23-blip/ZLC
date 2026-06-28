package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var emailService *EmailService

func main() {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connected")

	emailService = NewEmailService()
	log.Println("Email service initialized")

	go updateUserMetrics()

	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/api/register", handleRegister)
	http.HandleFunc("/api/verify-email", handleVerifyEmail)
	http.HandleFunc("/api/resend-email-verification", handleResendEmailVerification)
	http.HandleFunc("/api/send-phone-otp", handleSendPhoneOTP)
	http.HandleFunc("/api/verify-phone-otp", handleVerifyPhoneOTP)
	http.HandleFunc("/api/resend-phone-otp", handleResendPhoneOTP)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/verify-2fa", handleVerify2FA)
	http.HandleFunc("/api/delete-account", handleDeleteAccount)
	http.HandleFunc("/api/user/dashboard", authMiddleware(handleUserDashboard))
	http.HandleFunc("/api/admin/logs", authMiddleware(handleGetLogs))
	http.HandleFunc("/api/admin/users", authMiddleware(handleGetUsers))
	http.HandleFunc("/api/admin/dashboard-stats", authMiddleware(handleAdminDashboardStats))
	http.HandleFunc("/api/admin/login-trends", authMiddleware(handleLoginTrends))
	http.HandleFunc("/api/admin/user-activity", authMiddleware(handleUserActivityAudit))
	http.HandleFunc("/api/admin/all-users", authMiddleware(handleAdminGetUsers))
	http.HandleFunc("/api/admin/unlock-account", authMiddleware(handleUnlockAccount))
	http.HandleFunc("/api/users", handleGetUsers)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, map[string]string{"status": "ok"})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req map[string]string
	json.NewDecoder(r.Body).Decode(&req)

	email, password, phone := req["email"], req["password"], req["phone_number"]
	if email == "" || password == "" || phone == "" {
		respondJSON(w, 400, map[string]string{"error": "Email, password, and phone number required"})
		return
	}

	if len(phone) < 10 {
		respondJSON(w, 400, map[string]string{"error": "Invalid phone number"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	max := big.NewInt(10000000)
	n, _ := rand.Int(rand.Reader, max)
	emailToken := fmt.Sprintf("%07d", n.Int64())

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO users (email, phone_number, password, role_id, email_verified, phone_verified, 
			verification_token, verification_token_expires)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		email, phone, string(hash), 3, false, false, emailToken, time.Now().Add(10*time.Minute),
	)

	if err != nil {
		respondJSON(w, 409, map[string]string{"error": "Email already registered"})
		return
	}

	incrementRegistrations()
	recordAuditLog("USER_REGISTERED")
	logAudit(0, "USER_REGISTERED", getIP(r), email)

	emailService.SendVerificationEmail(email, email, emailToken)

	respondJSON(w, 201, map[string]interface{}{
		"message": "Registration successful. Verification email sent.",
		"step":    "email_verification",
	})
}

func handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req map[string]string
	json.NewDecoder(r.Body).Decode(&req)

	var storedToken string
	var expiresAt time.Time
	var userID int
	var email string

	err := db.QueryRowContext(context.Background(),
		`SELECT id, verification_token, verification_token_expires, email FROM users WHERE email = $1`,
		req["email"],
	).Scan(&userID, &storedToken, &expiresAt, &email)

	if err != nil || storedToken != req["token"] || time.Now().After(expiresAt) {
		respondJSON(w, 401, map[string]string{"error": "Invalid or expired token"})
		return
	}

	db.ExecContext(context.Background(),
		`UPDATE users SET email_verified = true, verification_token = NULL WHERE id = $1`,
		userID,
	)

	incrementEmailVerifications()
	recordAuditLog("EMAIL_VERIFIED")
	logAudit(userID, "EMAIL_VERIFIED", getIP(r), req["email"])

	otp := generateOTP()
	log.Printf("[VERIFY_EMAIL] Auto-sending phone OTP: '%s' (len=%d, bytes=%v) for user_id=%d", otp, len(otp), []byte(otp), userID)
	otpExpiry, _ := strconv.Atoi(os.Getenv("OTP_EXPIRY"))
	if otpExpiry == 0 {
		otpExpiry = 600
	}

	result, err := db.ExecContext(context.Background(),
		`UPDATE users SET phone_otp = $1, phone_otp_expires = $2, otp_resend_count = 0 WHERE id = $3`,
		otp, time.Now().Add(time.Duration(otpExpiry)*time.Second), userID,
	)

	if err != nil {
		log.Printf("[VERIFY_EMAIL] Failed to store OTP: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("[VERIFY_EMAIL] Stored OTP, rows affected: %d", rowsAffected)

	var storedOTP string
	var storedExpiry time.Time
	verifyErr := db.QueryRowContext(context.Background(),
		`SELECT phone_otp, phone_otp_expires FROM users WHERE id = $1`,
		userID,
	).Scan(&storedOTP, &storedExpiry)

	if verifyErr == nil {
		log.Printf("[VERIFY_EMAIL] Verification: stored_otp='%s' (len=%d, bytes=%v)", storedOTP, len(storedOTP), []byte(storedOTP))
	} else {
		log.Printf("[VERIFY_EMAIL] Verification query failed: %v", verifyErr)
	}

	emailService.SendPhoneOTPEmail(email, otp)
	recordAuditLog("PHONE_OTP_SENT")
	logAudit(userID, "PHONE_OTP_SENT", getIP(r), email)

	respondJSON(w, 200, map[string]interface{}{
		"message": "Email verified. Phone OTP sent to your email.",
		"user_id": userID,
		"step":    "phone_verification",
	})
}

func handleResendEmailVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req map[string]string
	json.NewDecoder(r.Body).Decode(&req)

	email := req["email"]
	if email == "" {
		respondJSON(w, 400, map[string]string{"error": "Email required"})
		return
	}

	max := big.NewInt(10000000)
	n, _ := rand.Int(rand.Reader, max)
	token := fmt.Sprintf("%07d", n.Int64())
	expiresAt := time.Now().Add(10 * time.Minute)

	result, err := db.ExecContext(context.Background(),
		`UPDATE users SET verification_token = $1, verification_token_expires = $2 
		 WHERE email = $3 AND email_verified = false`,
		token, expiresAt, email,
	)

	if err != nil {
		respondJSON(w, 400, map[string]string{"error": "Failed to resend token"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondJSON(w, 404, map[string]string{"error": "User not found or already verified"})
		return
	}

	emailService.SendVerificationEmail(email, email, token)
	recordAuditLog("RESEND_EMAIL_VERIFICATION")
	logAudit(0, "RESEND_EMAIL_VERIFICATION", getIP(r), email)

	respondJSON(w, 200, map[string]string{"message": "Verification email resent"})
}

func handleSendPhoneOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)

	userIDVal, exists := req["user_id"]
	if !exists {
		respondJSON(w, 400, map[string]string{"error": "user_id required"})
		return
	}

	userID := int(userIDVal.(float64))
	var email string

	err := db.QueryRowContext(context.Background(),
		`SELECT email FROM users WHERE id = $1 AND email_verified = true`,
		userID,
	).Scan(&email)

	if err != nil {
		log.Printf("send-phone-otp: user not found or email not verified: user_id=%d, err=%v", userID, err)
		respondJSON(w, 404, map[string]string{"error": "User not found or email not verified"})
		return
	}

	otp := generateOTP()
	log.Printf("[SEND_PHONE_OTP] Generated OTP: '%s' (len=%d, bytes=%v) for user_id=%d", otp, len(otp), []byte(otp), userID)
	otpExpiry, _ := strconv.Atoi(os.Getenv("OTP_EXPIRY"))
	if otpExpiry == 0 {
		otpExpiry = 600
	}

	result, err := db.ExecContext(context.Background(),
		`UPDATE users SET phone_otp = $1, phone_otp_expires = $2, otp_resend_count = 0 WHERE id = $3`,
		otp, time.Now().Add(time.Duration(otpExpiry)*time.Second), userID,
	)

	if err != nil {
		log.Printf("[SEND_PHONE_OTP] Failed to update OTP in DB: %v", err)
		respondJSON(w, 500, map[string]string{"error": "Failed to send OTP"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("[SEND_PHONE_OTP] Updated database, rows affected: %d", rowsAffected)

	var storedOTP string
	var storedExpiry time.Time
	verifyErr := db.QueryRowContext(context.Background(),
		`SELECT phone_otp, phone_otp_expires FROM users WHERE id = $1`,
		userID,
	).Scan(&storedOTP, &storedExpiry)

	if verifyErr == nil {
		log.Printf("[SEND_PHONE_OTP] Verification: stored_otp='%s' (len=%d, bytes=%v), expires=%v", storedOTP, len(storedOTP), []byte(storedOTP), storedExpiry)
	} else {
		log.Printf("[SEND_PHONE_OTP] Verification query failed: %v", verifyErr)
	}

	emailService.SendPhoneOTPEmail(email, otp)
	recordAuditLog("PHONE_OTP_SENT")
	logAudit(userID, "PHONE_OTP_SENT", getIP(r), email)

	respondJSON(w, 200, map[string]interface{}{
		"message": "OTP sent to email",
		"expires": otpExpiry,
	})
}

func handleVerifyPhoneOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)

	userID := int(req["user_id"].(float64))
	log.Printf("[VERIFY_PHONE_OTP] Starting verification for user_id=%d", userID)

	var otp string
	switch v := req["otp"].(type) {
	case string:
		otp = v
		log.Printf("[VERIFY_PHONE_OTP] OTP received as string: '%s' (len=%d)", otp, len(otp))
	case float64:
		otp = fmt.Sprintf("%d", int(v))
		log.Printf("[VERIFY_PHONE_OTP] OTP received as float64: %.0f, converted to string: '%s' (len=%d)", v, otp, len(otp))
	default:
		otp = fmt.Sprintf("%v", v)
		log.Printf("[VERIFY_PHONE_OTP] OTP received as other type, converted to string: '%s' (len=%d)", otp, len(otp))
	}

	var storedOTP string
	var expiresAt time.Time
	var email string

	err := db.QueryRowContext(context.Background(),
		`SELECT phone_otp, phone_otp_expires, email FROM users WHERE id = $1`,
		userID,
	).Scan(&storedOTP, &expiresAt, &email)

	if err != nil {
		log.Printf("[VERIFY_PHONE_OTP] Database query failed: %v", err)
		recordLoginAttempt("phone_otp_failed")
		recordAuditLog("PHONE_OTP_FAILED")
		respondJSON(w, 401, map[string]string{"error": "Invalid or expired OTP"})
		return
	}

	log.Printf("[VERIFY_PHONE_OTP] Database retrieved - stored_otp: '%s' (len=%d, type=%T), expires_at: %v, email: %s", storedOTP, len(storedOTP), storedOTP, expiresAt, email)
	log.Printf("[VERIFY_PHONE_OTP] Provided OTP: '%s' (len=%d, type=%T)", otp, len(otp), otp)
	log.Printf("[VERIFY_PHONE_OTP] OTP comparison: storedOTP == otp: %v", storedOTP == otp)
	log.Printf("[VERIFY_PHONE_OTP] Stored OTP bytes: %v", []byte(storedOTP))
	log.Printf("[VERIFY_PHONE_OTP] Provided OTP bytes: %v", []byte(otp))
	log.Printf("[VERIFY_PHONE_OTP] Stored OTP (trimmed): '%s'", strings.TrimSpace(storedOTP))
	log.Printf("[VERIFY_PHONE_OTP] Provided OTP (trimmed): '%s'", strings.TrimSpace(otp))

	normalizedStored := strings.TrimSpace(storedOTP)
	normalizedProvided := strings.TrimSpace(otp)

	if normalizedStored != normalizedProvided {
		log.Printf("[VERIFY_PHONE_OTP] OTP MISMATCH DETECTED - stored: '%s' (len=%d), provided: '%s' (len=%d)", normalizedStored, len(normalizedStored), normalizedProvided, len(normalizedProvided))
		log.Printf("[VERIFY_PHONE_OTP] Stored OTP is NULL or empty: %v, Provided OTP is NULL or empty: %v", normalizedStored == "", normalizedProvided == "")
		recordLoginAttempt("phone_otp_failed")
		recordAuditLog("PHONE_OTP_FAILED")
		respondJSON(w, 401, map[string]string{"error": "Invalid OTP"})
		return
	}

	if time.Now().After(expiresAt) {
		log.Printf("[VERIFY_PHONE_OTP] OTP EXPIRED - expiry: %v, now: %v", expiresAt, time.Now())
		recordLoginAttempt("phone_otp_failed")
		recordAuditLog("PHONE_OTP_FAILED")
		respondJSON(w, 401, map[string]string{"error": "OTP expired"})
		return
	}

	log.Printf("[VERIFY_PHONE_OTP] OTP verified successfully for user_id=%d", userID)
	db.ExecContext(context.Background(),
		`UPDATE users SET phone_verified = true, phone_otp = NULL, phone_otp_expires = NULL WHERE id = $1`,
		userID,
	)

	recordLoginAttempt("phone_otp_success")
	recordAuditLog("PHONE_VERIFIED")
	logAudit(userID, "PHONE_VERIFIED", getIP(r), email)

	respondJSON(w, 200, map[string]interface{}{
		"message": "Phone number verified. Registration complete!",
		"user_id": userID,
	})
}

func handleResendPhoneOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)

	userID := int(req["user_id"].(float64))
	var email string
	var resendCount int
	var lastResendTime sql.NullTime

	err := db.QueryRowContext(context.Background(),
		`SELECT email, otp_resend_count, otp_resend_last_time FROM users WHERE id = $1`,
		userID,
	).Scan(&email, &resendCount, &lastResendTime)

	if err != nil {
		respondJSON(w, 404, map[string]string{"error": "User not found"})
		return
	}

	if resendCount >= 3 && lastResendTime.Valid {
		if time.Since(lastResendTime.Time) < 10*time.Minute {
			respondJSON(w, 429, map[string]string{"error": "Too many resend attempts. Try again later."})
			return
		}
	}

	otp := generateOTP()
	otpExpiry, _ := strconv.Atoi(os.Getenv("OTP_EXPIRY"))
	if otpExpiry == 0 {
		otpExpiry = 600
	}

	newCount := resendCount + 1
	if lastResendTime.Valid && time.Now().Sub(lastResendTime.Time) > 10*time.Minute {
		newCount = 1
	}

	db.ExecContext(context.Background(),
		`UPDATE users SET phone_otp = $1, phone_otp_expires = $2, otp_resend_count = $3, otp_resend_last_time = $4 WHERE id = $5`,
		otp, time.Now().Add(time.Duration(otpExpiry)*time.Second), newCount, time.Now(), userID,
	)

	emailService.SendPhoneOTPEmail(email, otp)
	recordAuditLog("PHONE_OTP_RESENT")
	logAudit(userID, "PHONE_OTP_RESENT", getIP(r), email)

	respondJSON(w, 200, map[string]interface{}{
		"message": "OTP resent to email",
		"expires": otpExpiry,
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req map[string]string
	json.NewDecoder(r.Body).Decode(&req)

	var userID int
	var hash string
	var emailVerified, phoneVerified, accountLocked bool
	var lockedUntil sql.NullTime
	var failedAttempts int

	err := db.QueryRowContext(context.Background(),
		`SELECT id, password, email_verified, phone_verified, account_locked, locked_until, failed_login_attempts 
		 FROM users WHERE email = $1`,
		req["email"],
	).Scan(&userID, &hash, &emailVerified, &phoneVerified, &accountLocked, &lockedUntil, &failedAttempts)

	if accountLocked {
		if lockedUntil.Valid && time.Now().Before(lockedUntil.Time) {
			recordLoginAttempt("failed")
			logAudit(0, "LOGIN_FAILED_ACCOUNT_LOCKED", getIP(r), req["email"])
			db.ExecContext(context.Background(),
				`INSERT INTO login_history (user_id, ip_address, device_info, browser_info, status, failure_reason) 
				 VALUES ($1, $2, $3, $4, $5, $6)`,
				userID, getIP(r), getDeviceInfo(r), getBrowserInfo(r), "failed", "Account locked",
			)
			respondJSON(w, 403, map[string]string{"error": "Account is locked. Please try again later."})
			return
		}
		db.ExecContext(context.Background(),
			`UPDATE users SET account_locked = false, locked_until = NULL WHERE id = $1`,
			userID,
		)
	}

	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req["password"])) != nil {
		newFailedAttempts := failedAttempts + 1
		var lockTime sql.NullTime

		if newFailedAttempts >= 5 {
			lockTime.Valid = true
			lockTime.Time = time.Now().Add(30 * time.Minute)
			db.ExecContext(context.Background(),
				`UPDATE users SET failed_login_attempts = $1, last_failed_login = $2, account_locked = true, locked_until = $3 WHERE id = $4`,
				newFailedAttempts, time.Now(), lockTime.Time, userID,
			)
		} else {
			db.ExecContext(context.Background(),
				`UPDATE users SET failed_login_attempts = $1, last_failed_login = $2 WHERE id = $3`,
				newFailedAttempts, time.Now(), userID,
			)
		}

		recordLoginAttempt("failed")
		logAudit(0, "LOGIN_FAILED", getIP(r), req["email"])
		db.ExecContext(context.Background(),
			`INSERT INTO login_history (user_id, ip_address, device_info, browser_info, status, failure_reason) 
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			userID, getIP(r), getDeviceInfo(r), getBrowserInfo(r), "failed", "Invalid credentials",
		)
		respondJSON(w, 401, map[string]string{"error": "Invalid credentials"})
		return
	}

	if !emailVerified || !phoneVerified {
		respondJSON(w, 403, map[string]string{"error": "Email and phone must be verified"})
		return
	}

	otp := generateOTP()
	db.ExecContext(context.Background(),
		`UPDATE users SET otp_token = $1, otp_expires = $2 WHERE id = $3`,
		otp, time.Now().Add(5*time.Minute), userID,
	)

	db.ExecContext(context.Background(),
		`INSERT INTO login_history (user_id, ip_address, device_info, browser_info, status) 
		 VALUES ($1, $2, $3, $4, $5)`,
		userID, getIP(r), getDeviceInfo(r), getBrowserInfo(r), "pending",
	)

	db.ExecContext(context.Background(),
		`UPDATE users SET failed_login_attempts = 0 WHERE id = $1`,
		userID,
	)

	recordLoginAttempt("2fa_sent")
	recordAuditLog("LOGIN_INITIATED")
	logAudit(userID, "LOGIN_INITIATED", getIP(r), req["email"])

	emailService.SendLoginOTPEmail(req["email"], otp)

	respondJSON(w, 200, map[string]interface{}{
		"message": "Login OTP sent to email",
		"user_id": userID,
	})
}

func handleVerify2FA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)

	userID := int(req["user_id"].(float64))

	var storedOTP string
	var expiresAt time.Time
	var roleID int
	var email string

	err := db.QueryRowContext(context.Background(),
		`SELECT otp_token, otp_expires, role_id, email FROM users WHERE id = $1`,
		userID,
	).Scan(&storedOTP, &expiresAt, &roleID, &email)

	if err != nil || storedOTP != fmt.Sprintf("%v", req["token"]) || time.Now().After(expiresAt) {
		recordLoginAttempt("2fa_failed")
		recordAuditLog("LOGIN_2FA_FAILED")
		respondJSON(w, 401, map[string]string{"error": "Invalid or expired 2FA token"})
		return
	}

	db.ExecContext(context.Background(),
		`UPDATE users SET otp_token = NULL, last_login = $1 WHERE id = $2`,
		time.Now(), userID,
	)

	db.ExecContext(context.Background(),
		`UPDATE login_history SET status = 'success' WHERE user_id = $1 AND status = 'pending' 
		 ORDER BY id DESC LIMIT 1`,
		userID,
	)

	token := fmt.Sprintf("jwt.%d.%d", userID, roleID)
	recordLoginAttempt("success")
	recordAuditLog("LOGIN_SUCCESSFUL")
	logAudit(userID, "LOGIN_SUCCESSFUL", getIP(r), email)

	respondJSON(w, 200, map[string]interface{}{
		"token":   token,
		"user_id": userID,
		"role_id": roleID,
	})
}

func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.QueryContext(context.Background(),
		`SELECT id, user_id, action, timestamp, ip_address, details FROM audit_logs ORDER BY timestamp DESC LIMIT 100`,
	)
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var action, ip, details string
		var ts time.Time
		rows.Scan(&id, &userID, &action, &ts, &ip, &details)
		logs = append(logs, map[string]interface{}{
			"id": id, "user_id": userID, "action": action, "timestamp": ts, "ip_address": ip, "details": details,
		})
	}

	respondJSON(w, 200, logs)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.QueryContext(context.Background(),
		`SELECT id, email, phone_number, role_id, email_verified, phone_verified, created_at FROM users`,
	)
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id, roleID int
		var email, phone string
		var emailVerified, phoneVerified bool
		var createdAt time.Time
		rows.Scan(&id, &email, &phone, &roleID, &emailVerified, &phoneVerified, &createdAt)
		users = append(users, map[string]interface{}{
			"id":             id,
			"email":          email,
			"phone_number":   phone,
			"role_id":        roleID,
			"email_verified": emailVerified,
			"phone_verified": phoneVerified,
			"created_at":     createdAt,
		})
	}

	respondJSON(w, 200, users)
}

func handleUserDashboard(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		respondJSON(w, 401, map[string]string{"error": "Missing token"})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		respondJSON(w, 401, map[string]string{"error": "Invalid token format"})
		return
	}

	userIDStr := parts[1]
	var userID int
	fmt.Sscanf(userIDStr, "%d", &userID)

	var email, phone string
	var roleID int
	var lastLogin sql.NullTime

	err := db.QueryRowContext(context.Background(),
		`SELECT email, phone_number, role_id, last_login FROM users WHERE id = $1`,
		userID,
	).Scan(&email, &phone, &roleID, &lastLogin)

	if err != nil {
		respondJSON(w, 404, map[string]string{"error": "User not found"})
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"user_id":      userID,
		"email":        email,
		"phone_number": phone,
		"role_id":      roleID,
		"last_login":   lastLogin.Time,
	})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

// *** ONLY THIS FUNCTION CHANGED ***
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := []string{
			"http://localhost:9031",
			"http://127.0.0.1:9031",
		}

		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func logAudit(userID int, action, ip, details string) {
	db.ExecContext(context.Background(),
		`INSERT INTO audit_logs (user_id, action, ip_address, details, timestamp) VALUES ($1, $2, $3, $4, $5)`,
		userID, action, ip, details, time.Now(),
	)
}

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func updateUserMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var verified, unverified int64
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email_verified = true AND phone_verified = true").Scan(&verified)
		if err == nil {
			updateVerifiedUsers(float64(verified))
		}
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email_verified = false OR phone_verified = false").Scan(&unverified)
		if err == nil {
			updateUnverifiedUsers(float64(unverified))
		}
	}
}

func generateOTP() string {
	max := big.NewInt(1000000)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64())
}

func getDeviceInfo(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return "Unknown"
	}
	if strings.Contains(userAgent, "Mobile") {
		return "Mobile"
	} else if strings.Contains(userAgent, "Tablet") {
		return "Tablet"
	}
	return "Desktop"
}

func getBrowserInfo(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return "Unknown"
	}
	switch {
	case strings.Contains(userAgent, "Chrome"):
		return "Chrome"
	case strings.Contains(userAgent, "Firefox"):
		return "Firefox"
	case strings.Contains(userAgent, "Safari"):
		return "Safari"
	case strings.Contains(userAgent, "Edge"):
		return "Edge"
	default:
		return "Other"
	}
}