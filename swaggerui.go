package main

import (
	"embed"
)

//go:embed swagger-ui/*
var swaggerUiFolder embed.FS
