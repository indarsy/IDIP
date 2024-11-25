package services

import (
	"fmt"
	"sync"
	"time"
	"videodb/be/config"

	"github.com/aler9/gortsplib"
	"github.com/aler9/gortsplib/pkg/url"
	"github.com/pion/webrtc/v3"
)

type WebRTCService struct {
	config    *config.Config
	connMutex sync.RWMutex
	connMap   map[string]*webrtc.PeerConnection
	rtspMap   map[string]*gortsplib.Client
}

func NewWebRTCService(config *config.Config) *WebRTCService {
	return &WebRTCService{
		config:  config,
		connMap: make(map[string]*webrtc.PeerConnection),
		rtspMap: make(map[string]*gortsplib.Client),
	}
}

func (s *WebRTCService) HandleRTSP(rtspURL string, offerSDP string) (*webrtc.SessionDescription, error) {
	// 创建 WebRTC 连接配置
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	// 设置远程描述（前端发来的 offer）
	err = peerConnection.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerSDP,
	})
	if err != nil {
		peerConnection.Close()
		return nil, fmt.Errorf("failed to set remote description: %v", err)
	}

	// 解析 RTSP URL
	u, err := url.Parse(rtspURL)
	if err != nil {
		peerConnection.Close()
		return nil, fmt.Errorf("failed to parse RTSP URL: %v", err)
	}

	// 创建 RTSP 客户端
	rtspClient := gortsplib.Client{}

	// 连接到 RTSP 服务器
	err = rtspClient.Start(u.Scheme, u.Host)
	if err != nil {
		peerConnection.Close()
		return nil, fmt.Errorf("failed to connect to RTSP: %v", err)
	}

	// 创建视频轨道
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264},
		"video",
		"pion",
	)
	if err != nil {
		rtspClient.Close()
		peerConnection.Close()
		return nil, fmt.Errorf("failed to create video track: %v", err)
	}

	// 添加轨道到连接
	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		rtspClient.Close()
		peerConnection.Close()
		return nil, fmt.Errorf("failed to add track: %v", err)
	}

	// 创建应答
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		rtspClient.Close()
		peerConnection.Close()
		return nil, fmt.Errorf("failed to create answer: %v", err)
	}

	// 设置本地描述
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		rtspClient.Close()
		peerConnection.Close()
		return nil, fmt.Errorf("failed to set local description: %v", err)
	}

	// 保存连接信息
	connID := fmt.Sprintf("%s-%d", rtspURL, time.Now().UnixNano())
	s.connMutex.Lock()
	s.connMap[connID] = peerConnection
	s.rtspMap[connID] = &rtspClient
	s.connMutex.Unlock()

	// 启动流转发
	go s.streamRTSPToWebRTC(connID, &rtspClient, videoTrack, peerConnection)

	return &answer, nil
}

func (s *WebRTCService) streamRTSPToWebRTC(
	connID string,
	rtspClient *gortsplib.Client,
	videoTrack *webrtc.TrackLocalStaticRTP,
	peerConnection *webrtc.PeerConnection,
) {
	defer func() {
		s.connMutex.Lock()
		delete(s.connMap, connID)
		delete(s.rtspMap, connID)
		s.connMutex.Unlock()

		rtspClient.Close()
		peerConnection.Close()
	}()

	// 设置分段长度（从配置中获取）
	segmentLength := s.config.RTSP.SegmentLength

	// 创建定时器用于分段处理
	ticker := time.NewTicker(time.Duration(segmentLength) * time.Second)
	defer ticker.Stop()

	// 设置 RTP 数据的回调
	rtspClient.OnPacketRTP = func(ctx *gortsplib.ClientOnPacketRTPCtx) {
		// 使用 Marshal() 替代 Raw
		data, err := ctx.Packet.Marshal()
		if err != nil {
			fmt.Printf("Error marshaling RTP packet: %v\n", err)
			return
		}
		_, err = videoTrack.Write(data)
		if err != nil {
			fmt.Printf("Error writing RTP: %v\n", err)
		}
	}

	// 分段处理逻辑
	for range ticker.C {
		// 在这里可以添加分段处理逻辑，例如日志记录或数据切片
		fmt.Println("Segment processed")
	}
}

func (s *WebRTCService) CloseConnection(connID string) {
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	if pc, ok := s.connMap[connID]; ok {
		pc.Close()
		delete(s.connMap, connID)
	}

	if client, ok := s.rtspMap[connID]; ok {
		client.Close()
		delete(s.rtspMap, connID)
	}
}
