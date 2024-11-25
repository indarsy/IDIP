import asyncio
from aiohttp import web
from aiortc import RTCPeerConnection, VideoStreamTrack, RTCSessionDescription
from aiortc.contrib.media import MediaPlayer
from av import VideoFrame
import aiohttp_cors
import av
import gc

av.logging.set_level(av.logging.ERROR)

pcs = set()
players = {}


class VideoTrack(VideoStreamTrack):
    def __init__(self, player):
        super().__init__()
        self.player = player
        self._last_frame = None

    async def recv(self):
        MAX_RETRIES = 3
        retry_count = 0
        
        while retry_count < MAX_RETRIES:
            try:
                frame = await self.player.video.recv()
                
                if frame.format.name not in ['nv12', 'yuv420p']:
                    frame = frame.reformat(format='yuv420p')
                
                self._last_frame = frame
                return frame
                
            except Exception as e:
                print(f"接收视频帧错误: {str(e)}")
                retry_count += 1
                if retry_count >= MAX_RETRIES:
                    if self._last_frame:
                        return self._last_frame
                    raise
                await asyncio.sleep(0.1)

    def __del__(self):
        self._last_frame = None


async def close_peer_connection(pc):
    """更完整的清理逻辑"""
    if pc in players:
        try:
            player = players[pc]
            # 停止所有轨道
            if hasattr(player, 'video'):
                player.video.close()
            if hasattr(player, 'audio'):
                player.audio.close()
            player.close()
        except Exception as e:
            print(f"清理 player 失败: {str(e)}")
        finally:
            del players[pc]
    
    try:
        # 停止所有轨道
        for sender in pc.getSenders():
            if sender.track:
                sender.track.stop()
        
        for receiver in pc.getReceivers():
            if receiver.track:
                receiver.track.stop()
                
        await pc.close()
    except Exception as e:
        print(f"清理 PC 失败: {str(e)}")
    finally:
        pcs.discard(pc)
        
    # 强制垃圾回收
    gc.collect()


async def offer(request):
    pc = None
    try:
        params = await request.json()
        rtsp_url = params.get("rtsp_url")
        if not rtsp_url:
            return web.json_response({"error": "RTSP URL is required"}, status=400)

        pc = RTCPeerConnection()
        pcs.add(pc)

        if pc in players:
            await close_peer_connection(pc)

        try:
            print(f"尝试连接 RTSP 流: {rtsp_url}")
            player = MediaPlayer(
                rtsp_url,
                format="rtsp",
                options={
                    'rtsp_transport': 'tcp',
                    'stimeout': '5000000',
                    'buffer_size': '1024000',
                    'max_delay': '500000',
                    'fflags': 'nobuffer',
                    'flags': 'low_delay',
                    'timeout': '5000000',
                    'reconnect': '1',
                    'reconnect_at_eof': '1',
                    'reconnect_streamed': '1',
                    'reconnect_delay_max': '2',
                    'pix_fmt': 'yuv420p',
                    'vsync': '0',
                    'framerate': '25',
                    'video_size': '1280x720',
                    'thread_queue_size': '512',
                    'max_error_rate': '0.99',
                }
            )
            
            if not player or not player.video:
                raise Exception("无法获取视频流")
                
            players[pc] = player
                
        except Exception as e:
            await close_peer_connection(pc)
            return web.json_response({"error": f"Failed to create MediaPlayer: {str(e)}"}, status=500)

        @pc.on("connectionstatechange")
        async def on_connectionstatechange():
            print(f"Connection state is {pc.connectionState}")
            if pc.connectionState == "failed" or pc.connectionState == "closed":
                await close_peer_connection(pc)

        try:
            video_track = VideoTrack(player)
            pc.addTrack(video_track)
            
            await pc.setRemoteDescription(RTCSessionDescription(params["sdp"], params["type"]))
            answer = await pc.createAnswer()
            await pc.setLocalDescription(answer)
            
            return web.json_response({
                "sdp": pc.localDescription.sdp, 
                "type": pc.localDescription.type
            })

        except Exception as e:
            await close_peer_connection(pc)
            return web.json_response({"error": str(e)}, status=500)

    except Exception as e:
        if pc:
            await close_peer_connection(pc)
        return web.json_response({"error": str(e)}, status=500)


async def cleanup(app):
    cleanup_tasks = []
    for pc in pcs.copy():
        cleanup_tasks.append(close_peer_connection(pc))
    
    if cleanup_tasks:
        await asyncio.gather(*cleanup_tasks)
    
    players.clear()
    pcs.clear()
    gc.collect()


app = web.Application()
app.router.add_post("/offer", offer)

# 配置 CORS 支持
cors = aiohttp_cors.setup(app, defaults={
    "*": aiohttp_cors.ResourceOptions(
        allow_credentials=True,
        expose_headers="*",
        allow_headers="*",
    )
})

for route in list(app.router.routes()):
    cors.add(route)

app.on_shutdown.append(cleanup)

if __name__ == "__main__":
    web.run_app(app, host="0.0.0.0", port=5000)