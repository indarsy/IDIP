# IDIP (Industrial Digital Intelligence Platform)

工业数智平台 - 集成视频采集与工业实时数据采集的智能监控平台

## 功能特点

- 📹 视频采集与管理
  - RTSP 视频流实时预览
  - 定时视频采集
  - 视频文件管理与回放

- 📊 工业实时数据
  - OPC UA 数据源管理
  - 实时数据采集
  - 历史数据查询

## 环境要求

### 后端环境
- Go >= 1.18
- Mysql >= 8.0
- FFmpeg >= 4.2
- Python >= 3.8

### 前端环境
- Node.js >= 16.0
- npm >= 8.0
- Vue.js 3.x

## 快速开始

### 1. 克隆项目
bash
git clone git@github.com:indarsy/IDIP.git
cd IDIP


### 2. 后端设置

```bash
cd be
go mod tidy
go run main.go
```
### 3. 前端设置

```bash
cd fe
npm install
npm run serve
```
### 4. RTSP服务
```bash
cd rtsp
python main.py
```
### 端口说明
- 后端 API：8882
- RTSP 服务：5000
- web服务端口：8080
