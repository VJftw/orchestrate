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
				"persister.gorm":     false,
				"user.manager":       false,
				"user.validator":     false,
				"user.provider":      false,
				"user.resolver":      false,
				"auth.provider":      false,
				"project.manager":    false,
				"project.provider":   false,
				"project.resolver":   false,
				"project.validator":  false,
				"user.controller":    false,
				"auth.controller":    false,
				"project.controller": false,
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
