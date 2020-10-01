package models

import (
	"testing"

	"github.com/magneticio/vamp-policies/policy/interface/persistence/vault/dto"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDTOValidator(t *testing.T) {
	Convey("Given DTO Validator", t, func() {
		validate := NewValidateDTO()

		Convey("and command DTO", func() {
			command := dto.Command{}

			Convey("When validating command DTO without type", func() {
				command.Value = []struct{ string }{{"test"}}
				err := validate(command)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "type field is required")
				})
			})

			Convey("When validating command DTO without value", func() {
				command.Type = "http"
				err := validate(command)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "value field is required")
				})
			})

			Convey("When validating correct command DTO", func() {
				command.Type = "http"
				command.Value = []struct{ string }{{"test"}}
				err := validate(command)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and condition DTO", func() {
			condition := dto.Condition{}

			Convey("When validating condition DTO without value", func() {
				err := validate(condition)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "value field is required")
				})
			})

			Convey("When validating condition DTO with value", func() {
				condition.Value = "condition value"
				err := validate(condition)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and elasticsearchMetricValue DTO", func() {
			elasticsearchMetricValue := dto.ElasticsearchMetricValue{}

			Convey("When validating elasticsearchMetricValue without source", func() {
				elasticsearchMetricValue.Index = "test-index"
				elasticsearchMetricValue.Tags = []string{"test-tag"}
				err := validate(elasticsearchMetricValue)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "source field is required")
				})
			})

			Convey("When validating elasticsearchMetricValue without index", func() {
				elasticsearchMetricValue.Source = "es"
				elasticsearchMetricValue.Tags = []string{"test-tag"}
				err := validate(elasticsearchMetricValue)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "index field is required")
				})
			})

			Convey("When validating elasticsearchMetricValue without tags", func() {
				elasticsearchMetricValue.Source = "es"
				elasticsearchMetricValue.Index = "test-index"
				err := validate(elasticsearchMetricValue)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "tags field is required")
				})
			})

			Convey("When validating elasticsearchMetricValue with empty tags", func() {
				elasticsearchMetricValue.Source = "es"
				elasticsearchMetricValue.Tags = []string{}
				elasticsearchMetricValue.Index = "test-index"
				err := validate(elasticsearchMetricValue)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "tags field must contain at least 1 element(s)")
				})
			})

			Convey("When validating correct elasticsearchMetricValue", func() {
				elasticsearchMetricValue.Source = "es"
				elasticsearchMetricValue.Tags = []string{"test-tag"}
				elasticsearchMetricValue.Index = "test-index"
				err := validate(elasticsearchMetricValue)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and endAfterCondition DTO", func() {
			endAfterCondition := dto.EndAfterCondition{}

			Convey("When validating endAfterCondition DTO without value", func() {
				err := validate(endAfterCondition)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "value field is required")
				})
			})

			Convey("When validating correct endAfterCondition DTO", func() {
				endAfterCondition.Value = "duration == 4m"
				err := validate(endAfterCondition)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and httpCommand DTO", func() {
			httpCommand := dto.HttpCommand{}

			Convey("When validating endAfterCondition DTO without URL", func() {
				err := validate(httpCommand)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "url field is required")
				})
			})

			Convey("When validating correct endAfterCondition DTO", func() {
				httpCommand.URL = "http://test.local"
				err := validate(httpCommand)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and k8sDeploymentHealthMetricValue DTO", func() {
			k8sDeploymentHealthMetricValue := dto.K8sDeploymentHealthMetricValue{}

			Convey("When validating k8sDeploymentHealthMetricValue DTO without source", func() {
				err := validate(k8sDeploymentHealthMetricValue)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "source field is required")
				})

			})

			Convey("When validating correct k8sDeploymentHealthMetricValue DTO", func() {
				k8sDeploymentHealthMetricValue.Source = "k8s-deployment-health"
				err := validate(k8sDeploymentHealthMetricValue)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and k8sStatefulSetHealthMetricValue DTO", func() {
			k8sStatefulSetHealthMetricValue := dto.K8sStatefulSetHealthMetricValue{}

			Convey("When validating k8sStatefulSetHealthMetricValue DTO without source", func() {
				err := validate(k8sStatefulSetHealthMetricValue)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "source field is required")
				})

			})

			Convey("When validating correct k8sStatefulSetHealthMetricValue DTO", func() {
				k8sStatefulSetHealthMetricValue.Source = "k8s-statefulset-health"
				err := validate(k8sStatefulSetHealthMetricValue)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and keyValueBaseline DTO", func() {
			keyValueBaseline := dto.KeyValueBaseline{}

			Convey("When validating keyValueBaseline DTO without name", func() {
				keyValueBaseline.Value = "2m30s"
				err := validate(keyValueBaseline)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "name field is required")
				})
			})

			Convey("When validating keyValueBaseline DTO without value", func() {
				keyValueBaseline.Name = "max-duration"
				err := validate(keyValueBaseline)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "value field is required")
				})

			})

			Convey("When validating correct keyValueBaseline DTO", func() {
				keyValueBaseline.Name = "max-duration"
				keyValueBaseline.Value = "2m30s"
				err := validate(keyValueBaseline)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and metric DTO", func() {
			metric := dto.Metric{}

			Convey("When validaitng metric DTO without name", func() {
				metric.Value = struct{ string }{"test-value"}
				err := validate(metric)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "name field is required")
				})
			})

			Convey("When validaitng metric DTO without value", func() {
				metric.Name = "min-requests"
				err := validate(metric)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "value field is required")
				})
			})

			Convey("When validating correct metric DTO", func() {
				metric.Name = "min-requests"
				metric.Value = struct{ string }{"test-value"}
				err := validate(metric)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and metricBaseline DTO", func() {
			metricBaseline := dto.MetricBaseline{}

			Convey("When validating metricBaseline DTO without name", func() {
				baselineValue := 1.23
				metricBaseline.Metric = "test-metric"
				metricBaseline.Value = &baselineValue
				err := validate(metricBaseline)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "name field is required")
				})
			})

			Convey("When validating metricBaseline DTO without metric", func() {
				baselineValue := 1.23
				metricBaseline.Name = "max-duration"
				metricBaseline.Value = &baselineValue
				err := validate(metricBaseline)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "metric field is required")
				})
			})

			Convey("When validating metricBaseline DTO without value", func() {
				metricBaseline.Name = "max-duration"
				metricBaseline.Metric = "test-metric"
				err := validate(metricBaseline)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "value field is required")
				})
			})

			Convey("When validating metricBaseline DTO with value equal to default float64 value", func() {
				var baselineValue float64
				metricBaseline.Name = "max-duration"
				metricBaseline.Metric = "test-metric"
				metricBaseline.Value = &baselineValue
				err := validate(metricBaseline)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When validating correct metricBaseline DTO", func() {
				baselineValue := 1.23
				metricBaseline.Name = "max-duration"
				metricBaseline.Metric = "test-metric"
				metricBaseline.Value = &baselineValue
				err := validate(metricBaseline)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and policy DTO", func() {
			policy := dto.Policy{}

			Convey("When validating policy DTO without type", func() {
				policy.Name = "test-policy"
				policy.Steps = []dto.Step{
					dto.Step{
						Source: dto.Route{
							Weight: 10,
						},
						Target: dto.Route{
							Weight: 90,
						},
						EndAfter: dto.EndAfterCondition{
							Value: "duration == 5m",
						},
					},
				}

				err := validate(policy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "type field is required")
				})
			})

			Convey("When validating policy DTO with wrong type", func() {
				policy.Name = "test-policy"
				policy.Type = "unknown"
				policy.Steps = []dto.Step{
					dto.Step{
						Source: dto.Route{
							Weight: 10,
						},
						Target: dto.Route{
							Weight: 90,
						},
						EndAfter: dto.EndAfterCondition{
							Value: "duration == 5m",
						},
					},
				}

				err := validate(policy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "type field is not valid")
				})
			})

			Convey("When validating policy DTO without name", func() {
				policy.Steps = []dto.Step{
					dto.Step{
						Source: dto.Route{
							Weight: 10,
						},
						Target: dto.Route{
							Weight: 90,
						},
						EndAfter: dto.EndAfterCondition{
							Value: "duration == 5m",
						},
					},
				}
				policy.Type = "release"

				err := validate(policy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "name field is required")
				})
			})

			Convey("When validating policy DTO without steps field", func() {
				policy.Name = "test-policy"
				policy.Type = "release"

				err := validate(policy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "steps field is required")
				})
			})

			Convey("When validating policy DTO with empty steps field", func() {
				policy.Name = "test-policy"
				policy.Type = "release"
				policy.Steps = []dto.Step{}

				err := validate(policy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "steps field must contain at least 1 element(s)")
				})
			})

			Convey("When validating policy DTO with step without endAfterCondition value", func() {
				policy.Name = "test-policy"
				policy.Type = "release"
				policy.Steps = []dto.Step{
					dto.Step{
						Source: dto.Route{
							Weight: 10,
						},
						Target: dto.Route{
							Weight: 90,
						},
						EndAfter: dto.EndAfterCondition{
							Name: "test",
						},
					},
				}

				err := validate(policy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "steps[0].endAfter.value field is required")
				})
			})

			Convey("When validating correct policy DTO", func() {
				policy.Name = "test-policy"
				policy.Type = "release"
				policy.Steps = []dto.Step{
					dto.Step{
						EndAfter: dto.EndAfterCondition{
							Value: "value",
						},
					},
				}

				err := validate(policy)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and prometheusMetricValue DTO", func() {
			prometheusMetricValue := dto.PrometheusMetricValue{}

			Convey("When validating prometheusMetricValue DTO without source", func() {
				prometheusMetricValue.Query = "test query"
				err := validate(prometheusMetricValue)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "source field is required")
				})
			})

			Convey("When validating prometheusMetricValue DTO without query", func() {
				prometheusMetricValue.Source = "prometheus"
				err := validate(prometheusMetricValue)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "query field is required")
				})
			})

			Convey("When validating correct prometheusMetricValue DTO", func() {
				prometheusMetricValue.Query = "test query"
				prometheusMetricValue.Source = "prometheus"
				err := validate(prometheusMetricValue)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and route DTO", func() {
			route := dto.Route{}

			Convey("When validating route DTO with weight below 0", func() {
				route.Weight = -1
				err := validate(route)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "weight field must be at least 0")
				})
			})

			Convey("When validating route DTO with weight above 100", func() {
				route.Weight = 101
				err := validate(route)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "weight field must be at most 100")
				})
			})

			Convey("When validating route DTO with conditionStrength below 0", func() {
				route.ConditionStrength = -1
				err := validate(route)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "conditionStrength field must be at least 0")
				})
			})

			Convey("When validating route DTO with conditionStrength above 100", func() {
				route.ConditionStrength = 101
				err := validate(route)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "conditionStrength field must be at most 100")
				})
			})

			Convey("When validating correct conditionStrength DTO", func() {
				route.Weight = 10
				route.ConditionStrength = 20
				err := validate(route)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and step DTO", func() {
			step := dto.Step{}

			Convey("When validating step DTO without endAfter condition value", func() {
				err := validate(step)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "endAfter.value field is required")
				})
			})

			Convey("When validating correct step DTO", func() {
				step.EndAfter = dto.EndAfterCondition{
					Value: "duration == 4m",
				}
				err := validate(step)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and validation policy DTO", func() {
			validationPolicy := dto.ValidationPolicy{}

			Convey("When validating validation policy DTO without type", func() {
				validationPolicy.Steps = []dto.ValidationStep{
					dto.ValidationStep{
						EndAfter: dto.EndAfterCondition{
							Value: "duration == 5m",
						},
					},
				}
				validationPolicy.Name = "test"

				err := validate(validationPolicy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "type field is required")
				})
			})

			Convey("When validating validation policy DTO with wrong type", func() {
				validationPolicy.Steps = []dto.ValidationStep{
					dto.ValidationStep{
						EndAfter: dto.EndAfterCondition{
							Value: "duration == 5m",
						},
					},
				}
				validationPolicy.Name = "test"
				validationPolicy.Type = "unknown"

				err := validate(validationPolicy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "type field is not valid")
				})
			})

			Convey("When validating validation policy DTO without name", func() {
				validationPolicy.Steps = []dto.ValidationStep{
					dto.ValidationStep{
						EndAfter: dto.EndAfterCondition{
							Value: "duration == 5m",
						},
					},
				}
				validationPolicy.Type = "validation"

				err := validate(validationPolicy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "name field is required")
				})
			})

			Convey("When validating validation policy DTO without steps field", func() {
				validationPolicy.Name = "test-policy"
				validationPolicy.Type = "validation"

				err := validate(validationPolicy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "steps field is required")
				})
			})

			Convey("When validating validation policy DTO with empty steps field", func() {
				validationPolicy.Name = "test-policy"
				validationPolicy.Type = "validation"
				validationPolicy.Steps = []dto.ValidationStep{}

				err := validate(validationPolicy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "steps field must contain at least 1 element(s)")
				})
			})

			Convey("When validating validation policy DTO with step without endAfterCondition value", func() {
				validationPolicy.Name = "test-policy"
				validationPolicy.Type = "validation"
				validationPolicy.Steps = []dto.ValidationStep{
					dto.ValidationStep{
						EndAfter: dto.EndAfterCondition{
							Name: "test",
						},
					},
				}

				err := validate(validationPolicy)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "steps[0].endAfter.value field is required")
				})
			})

			Convey("When validating correct validation policy DTO", func() {
				validationPolicy.Name = "test-policy"
				validationPolicy.Type = "validation"
				validationPolicy.Steps = []dto.ValidationStep{
					dto.ValidationStep{
						EndAfter: dto.EndAfterCondition{
							Value: "value",
						},
					},
				}

				err := validate(validationPolicy)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("and validation step DTO", func() {
			validationStep := dto.ValidationStep{}

			Convey("When validating validation step DTO without endAfter condition value", func() {
				err := validate(validationStep)

				Convey("error should be thrown", func() {
					So(err.Error(), ShouldEqual, "endAfter.value field is required")
				})
			})

			Convey("When validating correct validation step DTO", func() {
				validationStep.EndAfter = dto.EndAfterCondition{
					Value: "duration == 4m",
				}
				err := validate(validationStep)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
