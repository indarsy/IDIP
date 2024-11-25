const { defineConfig } = require('@vue/cli-service')
module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:8882',
        changeOrigin: true,
        pathRewrite: {
          '^/api': ''
        }
      },
      '/rtsp-api': {
        target: 'http://localhost:5000',
        changeOrigin: true,
        pathRewrite: {
          '^/rtsp-api': '/api'
        }
      }
    }
  }
}
