package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Admin Audit Dashboard Stats
func handleAdminDashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// Verify admin access
	if !isAdmin(r) {
		respondJSON(w, 403, map[string]string{"error": "Admin access required"})
		return
	}

	var totalUsers, newUsersToday, activeUsers, failedLogins, lockedAccounts, verificationPending int64

	// Total users
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)

	// New users today
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURRENT_DATE`).Scan(&newUsersToday)

	// Active users (logged in last 24 hours)
	db.QueryRow(`SELECT COUNT(DISTINCT user_id) FROM login_history WHERE login_time > NOW() - INTERVAL '24 hours' AND status = 'success'`).Scan(&activeUsers)

	// Failed login attempts today
	db.QueryRow(`SELECT COUNT(*) FROM login_history WHERE DATE(login_time) = CURRENT_DATE AND status = 'failed'`).Scan(&failedLogins)

	// Locked accounts
	db.QueryRow("SELECT COUNT(*) FROM users WHERE account_locked = true").Scan(&lockedAccounts)

	// Verification pending
	db.QueryRow("SELECT COUNT(*) FROM users WHERE email_verified = false OR phone_verified = false").Scan(&verificationPending)

	respondJSON(w, 200, map[string]interface{}{
		"total_users":            totalUsers,
		"new_users_today":        newUsersToday,
		"active_users":           activeUsers,
		"failed_login_attempts":  failedLogins,
		"locked_accounts":        lockedAccounts,
		"verification_pending":   verificationPending,
	})
}

// Login trends (last 7 days)
func handleLoginTrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if !isAdmin(r) {
		respondJSON(w, 403, map[string]string{"error": "Admin access required"})
		return
	}

	rows, _ := db.QueryContext(context.Background(), `
		SELECT DATE(login_time) as date, COUNT(*) as count, 
		       SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as successful,
		       SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed
		FROM login_history 
		WHERE login_time > NOW() - INTERVAL '7 days'
		GROUP BY DATE(login_time)
		ORDER BY DATE(login_time)
	`)
	defer rows.Close()

	var trends []map[string]interface{}
	for rows.Next() {
		var date time.Time
		var count, successful, failed int64
		rows.Scan(&date, &count, &successful, &failed)
		trends = append(trends, map[string]interface{}{
			"date":        date.Format("2006-01-02"),
			"total":       count,
			"successful":  successful,
			"failed":      failed,
		})
	}

	respondJSON(w, 200, trends)
}

// User activity audit for specific user
func handleUserActivityAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondJSON(w, 400, map[string]string{"error": "user_id required"})
		return
	}

	// Verify access (admin or self)
	if !isAdminOrSelf(r, userID) {
		respondJSON(w, 403, map[string]string{"error": "Access denied"})
		return
	}

	// Get login history
	loginHistoryRows, _ := db.QueryContext(context.Background(), `
		SELECT id, login_time, logout_time, ip_address, device_info, browser_info, 
		       location, status, failure_reason
		FROM login_history 
		WHERE user_id = $1 
		ORDER BY login_time DESC 
		LIMIT 50
	`, userID)
	defer loginHistoryRows.Close()

	var loginHistory []map[string]interface{}
	for loginHistoryRows.Next() {
		var id int64
		var loginTime, logoutTime time.Time
		var ipAddress, deviceInfo, browserInfo, location, status, failureReason string
		loginHistoryRows.Scan(&id, &loginTime, &logoutTime, &ipAddress, &deviceInfo, &browserInfo, &location, &status, &failureReason)

		var logoutTimeStr string
		if logoutTime.IsZero() {
			logoutTimeStr = ""
		} else {
			logoutTimeStr = logoutTime.Format(time.RFC3339)
		}

		loginHistory = append(loginHistory, map[string]interface{}{
			"id":              id,
			"login_time":      loginTime.Format(time.RFC3339),
			"logout_time":     logoutTimeStr,
			"ip_address":      ipAddress,
			"device_info":     deviceInfo,
			"browser_info":    browserInfo,
			"location":        location,
			"status":          status,
			"failure_reason":  failureReason,
		})
	}

	// Get password change logs
	passwordLogsRows, _ := db.QueryContext(context.Background(), `
		SELECT id, changed_at, ip_address, device_info, status
		FROM password_logs 
		WHERE user_id = $1 
		ORDER BY changed_at DESC 
		LIMIT 20
	`, userID)
	defer passwordLogsRows.Close()

	var passwordLogs []map[string]interface{}
	for passwordLogsRows.Next() {
		var id int64
		var changedAt time.Time
		var ipAddress, deviceInfo, status string
		passwordLogsRows.Scan(&id, &changedAt, &ipAddress, &deviceInfo, &status)
		passwordLogs = append(passwordLogs, map[string]interface{}{
			"id":          id,
			"changed_at":  changedAt.Format(time.RFC3339),
			"ip_address":  ipAddress,
			"device_info": deviceInfo,
			"status":      status,
		})
	}

	// Get contact change logs
	contactLogsRows, _ := db.QueryContext(context.Background(), `
		SELECT id, changed_at, change_type, old_value, new_value, ip_address, status
		FROM contact_change_logs 
		WHERE user_id = $1 
		ORDER BY changed_at DESC 
		LIMIT 20
	`, userID)
	defer contactLogsRows.Close()

	var contactLogs []map[string]interface{}
	for contactLogsRows.Next() {
		var id int64
		var changedAt time.Time
		var changeType, oldValue, newValue, ipAddress, status string
		contactLogsRows.Scan(&id, &changedAt, &changeType, &oldValue, &newValue, &ipAddress, &status)
		contactLogs = append(contactLogs, map[string]interface{}{
			"id":          id,
			"changed_at":  changedAt.Format(time.RFC3339),
			"type":        changeType,
			"old_value":   oldValue,
			"new_value":   newValue,
			"ip_address":  ipAddress,
			"status":      status,
		})
	}

	// Get failed login attempts
	failedLoginsRows, _ := db.QueryContext(context.Background(), `
		SELECT id, login_time, ip_address, device_info, browser_info, status, failure_reason
		FROM login_history 
		WHERE user_id = $1 AND status = 'failed'
		ORDER BY login_time DESC 
		LIMIT 20
	`, userID)
	defer failedLoginsRows.Close()

	var failedLogins []map[string]interface{}
	for failedLoginsRows.Next() {
		var id int64
		var loginTime time.Time
		var ipAddress, deviceInfo, browserInfo, status, failureReason string
		failedLoginsRows.Scan(&id, &loginTime, &ipAddress, &deviceInfo, &browserInfo, &status, &failureReason)
		failedLogins = append(failedLogins, map[string]interface{}{
			"id":              id,
			"login_time":      loginTime.Format(time.RFC3339),
			"ip_address":      ipAddress,
			"device_info":     deviceInfo,
			"browser_info":    browserInfo,
			"status":          status,
			"failure_reason":  failureReason,
		})
	}

	// Get profile update logs
	profileLogsRows, _ := db.QueryContext(context.Background(), `
		SELECT id, updated_at, field_name, old_value, new_value, ip_address
		FROM profile_update_logs 
		WHERE user_id = $1 
		ORDER BY updated_at DESC 
		LIMIT 20
	`, userID)
	defer profileLogsRows.Close()

	var profileLogs []map[string]interface{}
	for profileLogsRows.Next() {
		var id int64
		var updatedAt time.Time
		var fieldName, oldValue, newValue, ipAddress string
		profileLogsRows.Scan(&id, &updatedAt, &fieldName, &oldValue, &newValue, &ipAddress)
		profileLogs = append(profileLogs, map[string]interface{}{
			"id":         id,
			"updated_at": updatedAt.Format(time.RFC3339),
			"field":      fieldName,
			"old_value":  oldValue,
			"new_value":  newValue,
			"ip_address": ipAddress,
		})
	}

	respondJSON(w, 200, map[string]interface{}{
		"login_history":     loginHistory,
		"password_changes":  passwordLogs,
		"contact_changes":   contactLogs,
		"failed_logins":     failedLogins,
		"profile_updates":   profileLogs,
	})
}

// Get all users for admin
func handleAdminGetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if !isAdmin(r) {
		respondJSON(w, 403, map[string]string{"error": "Admin access required"})
		return
	}

	rows, _ := db.QueryContext(context.Background(), `
		SELECT id, email, phone_number, role_id, email_verified, phone_verified, 
		       account_locked, failed_login_attempts, created_at, last_login
		FROM users
		ORDER BY created_at DESC
	`)
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id, roleID int
		var email, phone string
		var emailVerified, phoneVerified, accountLocked bool
		var failedAttempts int
		var createdAt, lastLogin time.Time
		rows.Scan(&id, &email, &phone, &roleID, &emailVerified, &phoneVerified, &accountLocked, &failedAttempts, &createdAt, &lastLogin)

		users = append(users, map[string]interface{}{
			"id":                   id,
			"email":                email,
			"phone_number":         phone,
			"role_id":              roleID,
			"email_verified":       emailVerified,
			"phone_verified":       phoneVerified,
			"account_locked":       accountLocked,
			"failed_login_attempts": failedAttempts,
			"created_at":           createdAt.Format(time.RFC3339),
			"last_login":           lastLogin.Format(time.RFC3339),
		})
	}

	respondJSON(w, 200, users)
}

// Unlock account (admin only)
func handleUnlockAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if !isAdmin(r) {
		respondJSON(w, 403, map[string]string{"error": "Admin access required"})
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondJSON(w, 400, map[string]string{"error": "user_id required"})
		return
	}

	db.ExecContext(context.Background(),
		`UPDATE users SET account_locked = false, failed_login_attempts = 0, locked_until = NULL WHERE id = $1`,
		userID,
	)

	logAudit(0, "ACCOUNT_UNLOCKED", getIP(r), "Unlocked user "+userID)
	respondJSON(w, 200, map[string]string{"message": "Account unlocked successfully"})
}

// Helper functions
func isAdmin(r *http.Request) bool {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return false
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false
	}

	roleIDStr := parts[2]
	roleID, _ := strconv.Atoi(roleIDStr)
	return roleID == 1 // 1 = Admin
}

func isAdminOrSelf(r *http.Request, targetUserID string) bool {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return false
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false
	}

	userIDStr := parts[1]
	roleIDStr := parts[2]
	roleID, _ := strconv.Atoi(roleIDStr)

	// Admin can access all
	if roleID == 1 {
		return true
	}

	// User can only access their own
	return userIDStr == targetUserID
}
