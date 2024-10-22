package agent

import "time"

type Agent struct {
	ID            int
	Name          string
	Email         string
	IsOnline      bool
	CustomerCount int
}

type AgentDB struct {
	ID    int `gorm:"primaryKey"`
	Name  string
	Email string
}

type AgentQismoResponse struct {
	Data struct {
		Agents []struct {
			AvatarURL            string      `json:"avatar_url"`
			CreatedAt            string      `json:"created_at"`
			CurrentCustomerCount int         `json:"current_customer_count"`
			Email                string      `json:"email"`
			ForceOffline         bool        `json:"force_offline"`
			ID                   int         `json:"id"`
			IsAvailable          bool        `json:"is_available"`
			IsReqOtpReset        interface{} `json:"is_req_otp_reset"`
			LastLogin            time.Time   `json:"last_login"`
			Name                 string      `json:"name"`
			SdkEmail             string      `json:"sdk_email"`
			SdkKey               string      `json:"sdk_key"`
			Type                 int         `json:"type"`
			TypeAsString         string      `json:"type_as_string"`
			UserChannels         []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"user_channels"`
			UserRoles []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"user_roles"`
		} `json:"agents"`
	} `json:"data"`
	Meta struct {
		PerPage    int `json:"per_page"`
		TotalCount int `json:"total_count"`
	} `json:"meta"`
	Status int `json:"status"`
}

type AssignAgentQismoResponse struct {
	Data struct {
		AddedAgent struct {
			AvatarURL    interface{} `json:"avatar_url"`
			CreatedAt    string      `json:"created_at"`
			Email        string      `json:"email"`
			ForceOffline bool        `json:"force_offline"`
			ID           int         `json:"id"`
			IsAvailable  bool        `json:"is_available"`
			IsVerified   bool        `json:"is_verified"`
			LastLogin    time.Time   `json:"last_login"`
			Name         string      `json:"name"`
			SdkEmail     string      `json:"sdk_email"`
			SdkKey       string      `json:"sdk_key"`
			Type         int         `json:"type"`
			TypeAsString string      `json:"type_as_string"`
			UpdatedAt    string      `json:"updated_at"`
		} `json:"added_agent"`
		Service struct {
			CreatedAt             string      `json:"created_at"`
			FirstCommentID        string      `json:"first_comment_id"`
			FirstCommentTimestamp interface{} `json:"first_comment_timestamp"`
			IsResolved            bool        `json:"is_resolved"`
			LastCommentID         string      `json:"last_comment_id"`
			Notes                 interface{} `json:"notes"`
			ResolvedAt            interface{} `json:"resolved_at"`
			RetrievedAt           time.Time   `json:"retrieved_at"`
			RoomID                string      `json:"room_id"`
			UpdatedAt             string      `json:"updated_at"`
			UserID                int         `json:"user_id"`
		} `json:"service"`
	} `json:"data"`
}