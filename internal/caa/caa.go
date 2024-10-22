package caa

import "time"

type Queue struct {
	RoomID   int `gorm:"primaryKey`
	AgentID  int
	CreateAt time.Time
}

type QueueDB struct {
	RoomID    int `gorm:"primaryKey`
	AgentID   int `gorm:"default:null"`
	CreateAt  time.Time
	ResolveAt time.Time `gorm:"default:null"`
}

type QismoCAAWebhook struct {
	AppID         string `json:"app_id"`
	Source        string `json:"source"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	AvatarURL     string `json:"avatar_url"`
	Extras        string `json:"extras"`
	IsResolved    bool   `json:"is_resolved"`
	LatestService struct {
		ID                    int         `json:"id"`
		UserID                int         `json:"user_id"`
		RoomLogID             int         `json:"room_log_id"`
		AppID                 int         `json:"app_id"`
		RoomID                string      `json:"room_id"`
		Notes                 interface{} `json:"notes"`
		ResolvedAt            string      `json:"resolved_at"`
		IsResolved            bool        `json:"is_resolved"`
		CreatedAt             string      `json:"created_at"`
		UpdatedAt             string      `json:"updated_at"`
		FirstCommentID        string      `json:"first_comment_id"`
		LastCommentID         string      `json:"last_comment_id"`
		RetrievedAt           string      `json:"retrieved_at"`
		FirstCommentTimestamp interface{} `json:"first_comment_timestamp"`
	} `json:"latest_service"`
	RoomID         string `json:"room_id"`
	CandidateAgent struct {
		ID                  int         `json:"id"`
		Name                string      `json:"name"`
		Email               string      `json:"email"`
		AuthenticationToken string      `json:"authentication_token"`
		CreatedAt           string      `json:"created_at"`
		UpdatedAt           string      `json:"updated_at"`
		SdkEmail            string      `json:"sdk_email"`
		SdkKey              string      `json:"sdk_key"`
		IsAvailable         bool        `json:"is_available"`
		Type                int         `json:"type"`
		AvatarURL           string      `json:"avatar_url"`
		AppID               int         `json:"app_id"`
		IsVerified          bool        `json:"is_verified"`
		NotificationsRoomID string      `json:"notifications_room_id"`
		BubbleColor         string      `json:"bubble_color"`
		QismoKey            string      `json:"qismo_key"`
		DirectLoginToken    interface{} `json:"direct_login_token"`
		TypeAsString        string      `json:"type_as_string"`
		AssignedRules       []string    `json:"assigned_rules"`
	} `json:"candidate_agent"`
}

type QismoMarkAsResolvedWebhook struct {
	Customer struct {
		AdditionalInfo []struct {
			AnyKey string `json:"anykey"`
		} `json:"additional_info"`
		Avatar string `json:"avatar"`
		Name   string `json:"name"`
		UserID string `json:"user_id"`
	} `json:"customer"`
	ResolvedBy struct {
		Email       string `json:"email"`
		ID          int    `json:"id"`
		IsAvailable bool   `json:"is_available"`
		Name        string `json:"name"`
		Type        string `json:"type"`
	} `json:"resolved_by"`
	Service struct {
		FirstCommentID string `json:"first_comment_id"`
		ID             int    `json:"id"`
		IsResolved     bool   `json:"is_resolved"`
		LastCommentID  string `json:"last_comment_id"`
		Notes          string `json:"notes"`
		RoomID         string `json:"room_id"`
		Source         string `json:"source"`
	} `json:"service"`
}
