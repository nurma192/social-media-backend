package controllers

import "social-media-back/services"

type AppController struct {
	AppService *services.AppService
}

func NewController(appService *services.AppService) *AppController {
	return &AppController{AppService: appService}
}
