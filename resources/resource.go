package resources

import "embed"

//go:embed images/*
//go:embed audio/*
//go:embed music/*
//go:embed fonts/*
var resources embed.FS
