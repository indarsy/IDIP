package models

// WebRTCRequest 前端发送的请求结构
type WebRTCRequest struct {
	RTSPURL string `json:"rtspUrl"`
	SDP     string `json:"sdp"`
}

// WebRTCResponse 返回给前端的响应结构
type WebRTCResponse struct {
	SDP     string `json:"sdp"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
