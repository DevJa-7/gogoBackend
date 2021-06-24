package config

import "time"

// EmailServer is email server address
var EmailServer = "smtp.live.com"

// EmailPort is 465 by default
var EmailPort = 587

// EmailUsername is username
var EmailUsername = "tiandage719@outlook.com"

// EmailPassword is password
var EmailPassword = "DDavid19731233"

// Constants for Email service
const (
	EmailTimeout = 10 * time.Second
	EmailMethod  = "email"
)
