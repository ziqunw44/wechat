const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false,
  // 本地开发：HTTP + 常见端口。云上 HTTPS 可改回下方注释块（证书路径按环境调整）。
  devServer: {
    host: '0.0.0.0',
    port: 8080,
  },
  // devServer: {
  //   host: '0.0.0.0',
  //   https: {
  //     cert: fs.readFileSync(path.join(__dirname, 'src/assets/cert/127.0.0.1+2.pem')),
  //     key: fs.readFileSync(path.join(__dirname, 'src/assets/cert/127.0.0.1+2-key.pem')),
  //     // 或 Ubuntu: cert/key 指向 /etc/ssl/certs/server.crt 等
  //   },
  //   port: 443,
  // },
})
