package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/MasterWigu/Thesis/appServer/APIs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func declareAPIEndpoint() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"https://webserver-module:8070"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"access-control-allow-origin, access-control-allow-headers"}
	config.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(config))


	r.POST("/login", loginHandler)
	r.GET("/assets/types", getAssetTypesHandler)
	r.POST("/assets/types", addAssetTypeHandler)

	r.GET("/assets/:type", getAssetsByTypeHandler)
	
	r.GET("/asset/:id", getAssetHandler)
	r.POST("/asset", registerAssetHandler)
	r.POST("/asset/modify", modifyAssetHandler)
	r.DELETE("/asset/:assetId", removeAssetHandler)


	r.POST("tools/:toolId/plan", planActionHandler)
	r.POST("tools/:toolId/execute", executeActionHandler)
	r.POST("tools/:toolId/:actionId/confirm", confirmActionHandler)

	return r
}

func startListening(r *gin.Engine) {
	r.RunTLS(":8080", "/certs/cert.pem", "/certs/key.pem")
}

func sessionManager(c *gin.Context, changeContext bool) (*User, error) {
	sessionId, err := c.Cookie("SESSION_ID")
	if err != nil {
		if changeContext {
			c.String(http.StatusUnauthorized, "Session id not found, please login")
		}
		return nil, err
	}


	user, notFound := checkLoggedIn(sessionId)
	if notFound {
		if changeContext {
			c.String(http.StatusUnauthorized, "Invalid session id, please (re)login")
			c.SetCookie("SESSION_ID", sessionId, -1, "/", "broker-module", true, true) //Ask to delete cookie (max age = -1)
		}
		return nil, fmt.Errorf("invalid session id")
	}

	if changeContext {
		c.SetCookie("SESSION_ID", user.sessionId, 86400, "/", "broker-module", true, true)
	}

	return user, nil
}

func loginHandler(c *gin.Context) {
	// receives a zip with the user's credentials
	// stores the user in the logged in users list
	// creates cookie for the user
	// sends zip and cookie to the ledger module for effective verification
	// sends login status to user (and cookie)

	//_, err := sessionManager(c, false)
	//if err == nil { //If user ir already logged in
	//	c.String(http.StatusOK, fmt.Sprintf("Logged in!"))
	//	return
	//}


	formFile, err := c.FormFile("identityFile")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	file, err := formFile.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Generate 40 byte session id
	b := make([]byte, 40)
	_, err = rand.Read(b)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	sessId := base64.URLEncoding.EncodeToString(b)

	success, err := loginUser(file, formFile.Size, sessId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if !success {
		c.String(http.StatusUnauthorized, "Invalid credentials")
		return
	}

	c.SetCookie("SESSION_ID", sessId, 86400, "/", "broker-module", true, true)
	c.String(http.StatusOK, fmt.Sprintf("Logged in!"))
}


// getAssetTypesHandler godoc
// @Summary Gets the existing asset types
// @Description Makes a request to the ledger to get the existing asset types. Needs session id set as cookie
// @Tags root
// @Security ApiKeyAuth
// @Accept */*
// @Produce json
// @Success 200 {object} APIs.AssetTypesResp
// @Router /assets/types [get]
func getAssetTypesHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	responseStruct, err := processor.GetAssetTypes(user)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}

// addAssetTypeHandler godoc
// @Summary Adds a new asset type to the ledger
// @Description Makes a request to the ledger to add a new asset type. Needs session id set as cookie
// @Tags AssetTypes
// @Security ApiKeyAuth
// @Accept */*
// @Produce json
// @Success 200 {object} APIs.AssetTypesResp
// @Router /assets/types [post]
func addAssetTypeHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	if user.userPerm != "admin" {
		c.Status(http.StatusUnauthorized)
		return
	}

	newType := c.Query("newType")
	//newType := c.GetString("newType")
	if newType == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	responseStruct, err := processor.AddAssetType(user, newType)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}

func getAssetsByTypeHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}


	assetType := c.Param("type")
	//newType := c.GetString("newType")
	if assetType == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	responseStruct, err := processor.GetAssetsByType(user, assetType)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}

func getAssetHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	assetId := c.Param("id")
	if assetId == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	responseStruct, err := processor.GetAssetById(user, assetId)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}

func registerAssetHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	var newAsset APIs.AssetResp

	err = c.BindJSON(&newAsset)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	responseStruct, err := processor.RegisterAsset(user, &newAsset)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}

func modifyAssetHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	var newAsset APIs.AssetResp

	err = c.BindJSON(&newAsset)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	responseStruct, err := processor.ModifyAsset(user, &newAsset)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}


func removeAssetHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	assetId := c.Param("assetId")
	if assetId == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	err = processor.RemoveAsset(user, assetId)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}




func planActionHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	tool := c.Param("toolId")
	if tool == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	files, err := c.FormFile("toolResources")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	responseStruct, err := processor.PlanAction(user, tool, files)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}


func executeActionHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	tool := c.Param("toolId")
	if tool == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	files, err := c.FormFile("toolResources")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	responseStruct, err := processor.ExecuteAction(user, tool, files)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}

func confirmActionHandler(c *gin.Context) {
	user, err := sessionManager(c, true)
	if err != nil {
		log.Println(err)
		return
	}

	tool := c.Param("toolId")
	if tool == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	actionId := c.Param("actionId")
	if tool == "" {
		c.Status(http.StatusBadRequest)
		return
	}


	responseStruct, err := processor.ConfirmAction(user, tool, actionId)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, responseStruct)
}
