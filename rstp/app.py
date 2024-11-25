from flask import Flask, Response, jsonify
from flask_cors import CORS
import cv2
import threading
import time
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)  # 启用CORS支持

class VideoStream:
    def __init__(self):
        self.stream = None
        self.frame = None
        self.is_running = False
        self.lock = threading.Lock()
        self.video_width = 0
        self.video_height = 0
        self.aspect_ratio = 16/9  # 默认宽高比

    def start(self, rtsp_url):
        if self.is_running:
            return
        
        try:
            stream = cv2.VideoCapture(rtsp_url)
            if not stream.isOpened():
                raise Exception("无法连接到RTSP流")
                
            self.stream = stream
            self.is_running = True
            threading.Thread(target=self._read_stream, daemon=True).start()
            
            # 获取视频尺寸
            self.video_width = int(self.stream.get(cv2.CAP_PROP_FRAME_WIDTH))
            self.video_height = int(self.stream.get(cv2.CAP_PROP_FRAME_HEIGHT))
            if self.video_width and self.video_height:
                self.aspect_ratio = self.video_width / self.video_height
            
            logger.info(f"Video dimensions: {self.video_width}x{self.video_height}")
            
        except Exception as e:
            logger.error(f"启动流失败: {str(e)}")
            raise Exception(f"连接RTSP流失败: {str(e)}")

    def stop(self):
        self.is_running = False
        if self.stream:
            self.stream.release()
        self.stream = None
        self.frame = None

    def _read_stream(self):
        while self.is_running:
            if self.stream and self.stream.isOpened():
                ret, frame = self.stream.read()
                if ret:
                    with self.lock:
                        self.frame = frame
                else:
                    self.stop()
                    break
            time.sleep(0.01)

    def get_frame(self):
        with self.lock:
            return self.frame.copy() if self.frame is not None else None

    def is_active(self):
        return self.is_running and self.frame is not None

video_stream = VideoStream()

def generate_frames():
    while True:
        frame = video_stream.get_frame()
        if frame is None:
            time.sleep(0.1)
            continue
            
        ret, buffer = cv2.imencode('.jpg', frame, [cv2.IMWRITE_JPEG_QUALITY, 80])
        if not ret:
            continue
            
        yield (b'--frame\r\n'
               b'Content-Type: image/jpeg\r\n\r\n' + buffer.tobytes() + b'\r\n')

@app.route('/api/stream/start/<path:rtsp_url>')
def start_stream(rtsp_url):
    try:
        video_stream.stop()  # 停止现有流
        video_stream.start(f'rtsp://{rtsp_url}')
        return jsonify({
            'status': 'success',
            'message': '视频流已启动'
        })
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        }), 400

@app.route('/api/stream/stop')
def stop_stream():
    video_stream.stop()
    return jsonify({
        'status': 'success',
        'message': '视频流已停止'
    })

@app.route('/api/stream/status')
def stream_status():
    return jsonify({
        'status': 'success',
        'is_active': video_stream.is_active()
    })

@app.route('/api/stream/feed')
def video_feed():
    return Response(
        generate_frames(),
        mimetype='multipart/x-mixed-replace; boundary=frame'
    )

@app.route('/api/stream/info')
def stream_info():
    return jsonify({
        'status': 'success',
        'width': video_stream.video_width,
        'height': video_stream.video_height,
        'aspect_ratio': video_stream.aspect_ratio
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, threaded=True) 