const proxy = require("http-proxy-middleware")

// F'ing back because despite the documentation saying the create-react-app proxy feature
// supports websockets, it does not.

module.exports = app => {
  app.use('/api', proxy.createProxyMiddleware({target: "http://localhost:8100", ws: true}))
}