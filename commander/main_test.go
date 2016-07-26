package main

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestOrchestrateApp(t *testing.T) {
	convey.Convey("Given an Orchestrate App", t, func() {
		app := NewOrchestrateApp()

		convey.Convey("All the services should be present", func() {
			services := map[string]bool{
				"persister.gorm":      false,
				"manager.default":     false,
				"controller.user":     false,
				"validator.user":      false,
				"provider.user":       false,
				"provider.auth_token": false,
				"resolver.user":       false,
			}

			for _, element := range app.graph.Objects() {
				for serviceName := range services {
					if element.Name == serviceName {
						services[serviceName] = true
						break
					}
				}
			}

			for _, value := range services {
				convey.So(value, convey.ShouldBeTrue)
			}

		})
	})
}
