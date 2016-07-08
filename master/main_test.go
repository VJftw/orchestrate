package main

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestOrchestrateApp(t *testing.T) {
	convey.Convey("Given an Orchestrate App", t, func() {
		app := NewOrchestrateApp()

		convey.Convey("There should be a default.router", func() {
			hasService := false
			for _, element := range app.graph.Objects() {
				if element.Name == "default.router" {
					hasService = true
				}
			}
			convey.So(hasService, convey.ShouldBeTrue)
		})
	})
}
