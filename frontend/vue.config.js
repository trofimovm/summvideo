const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  // Disable ESLint during build
  lintOnSave: false,
  // Output to the backend's static folder
  outputDir: '../static/vue',
  // Set the base URL to root path
  publicPath: '/',
  // Configure the dev server to proxy API requests to the backend
  devServer: {
    proxy: {
      '/upload_video': {
        target: 'http://localhost:8000',
        changeOrigin: true
      },
      '/static': {
        target: 'http://localhost:8000',
        changeOrigin: true
      }
    }
  }
})