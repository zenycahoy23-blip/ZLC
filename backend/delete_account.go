package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", 405)
		return
	}

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
		respondJSON(w, 401, map[string]string{"error": "Invalid token"})
		return
	}

	tokenUserID := parts[1]
	var requestedUserID int
	fmt.Sscanf(tokenUserID, "%d", &requestedUserID)

	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)

	requestBodyUserID := int(req["user_id"].(float64))

	// Verify that user can only delete their own account
	if requestedUserID != requestBodyUserID {
		logAudit(requestedUserID, "DELETE_ACCOUNT_FAILED_UNAUTHORIZED", getIP(r), "")
		respondJSON(w, 403, map[string]string{"error": "Cannot delete other user accounts"})
		return
	}

	// Delete user from database (cascading deletes should handle related records)
	result, err := db.ExecContext(context.Background(),
		`DELETE FROM users WHERE id = $1 AND role_id != 1`,
		requestBodyUserID,
	)

	if err != nil {
		log.Printf("Error deleting user %d: %v", requestBodyUserID, err)
		logAudit(requestedUserID, "DELETE_ACCOUNT_FAILED", getIP(r), "Database error")
		respondJSON(w, 500, map[string]string{"error": "Failed to delete account"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logAudit(requestedUserID, "DELETE_ACCOUNT_FAILED", getIP(r), "User not found or admin account")
		respondJSON(w, 404, map[string]string{"error": "User not found or cannot delete admin account"})
		return
	}

	logAudit(requestedUserID, "ACCOUNT_DELETED", getIP(r), "Account self-deleted")
	log.Printf("User %d account deleted", requestBodyUserID)

	respondJSON(w, 200, map[string]string{"message": "Account deleted successfully"})
}
