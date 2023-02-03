package main

import (
	"net/http"
	"testing"
)

func TestConfig_AddDefaultData(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")
	td := testApp.AddDefaultData(&TemplateData{}, req)
	if td.Flash != "flash" {
		t.Error("failed to get flash data")
	}
	if td.Error != "warning" {
		t.Error("failed to get error data")
	}
	if td.Warning != "warning" {
		t.Error("failed to get warning data")
	}

}
