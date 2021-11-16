package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test from fabric</title>
</head>
<body>
  <h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

func main() {
	r := gin.Default()
	r.SetHTMLTemplate(html)

	r.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "https", gin.H{
			"status": "success",
		})
	})

	r.StaticFile("/web", "C:\\Users\\morei\\Desktop\\tese\\hyperledger_lab\\appServer\\guiServer\\goapp\\webpageSrc\\test.html")

	/*r.GET("/web", func(c *gin.Context) {
		c.HTML(http.StatusOK, "https", gin.H{
			"status": "success",
		})
	})*/

	// Listen and Server in https://127.0.0.1:8080
	r.RunTLS(":8070", "../../certshare/webserver-module1/cert.pem", "../../certshare/webserver-module1/key.pem")
	//r.Run(":8070")
}
