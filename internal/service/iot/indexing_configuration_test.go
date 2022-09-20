// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iot_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccIoTIndexingConfiguration_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]func(t *testing.T){
		"basic":         testAccIndexingConfiguration_basic,
		"allAttributes": testAccIndexingConfiguration_allAttributes,
	}

	acctest.RunSerialTests1Level(t, testCases, 0)
}

func testAccIndexingConfiguration_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_iot_indexing_configuration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexingConfigurationConfig_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "thing_group_indexing_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "thing_group_indexing_configuration.0.custom_field.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "thing_group_indexing_configuration.0.managed_field.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "thing_group_indexing_configuration.0.thing_group_indexing_mode", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.custom_field.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.device_defender_indexing_mode", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.managed_field.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.named_shadow_indexing_mode", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.thing_connectivity_indexing_mode", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.thing_indexing_mode", "OFF"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIndexingConfiguration_allAttributes(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_iot_indexing_configuration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexingConfigurationConfig_allAttributes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "thing_group_indexing_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "thing_group_indexing_configuration.0.custom_field.#", "0"),
					acctest.CheckResourceAttrGreaterThanValue(resourceName, "thing_group_indexing_configuration.0.managed_field.#", 0),
					resource.TestCheckResourceAttr(resourceName, "thing_group_indexing_configuration.0.thing_group_indexing_mode", "ON"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.custom_field.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "thing_indexing_configuration.0.custom_field.*", map[string]string{
						"name": "attributes.version",
						"type": "Number",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "thing_indexing_configuration.0.custom_field.*", map[string]string{
						"name": "shadow.name.thing1shadow.desired.DefaultDesired",
						"type": "String",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "thing_indexing_configuration.0.custom_field.*", map[string]string{
						"name": "deviceDefender.securityProfile1.NUMBER_VALUE_BEHAVIOR.lastViolationValue.number",
						"type": "Number",
					}),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.device_defender_indexing_mode", "VIOLATIONS"),
					acctest.CheckResourceAttrGreaterThanValue(resourceName, "thing_group_indexing_configuration.0.managed_field.#", 0),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.named_shadow_indexing_mode", "ON"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.thing_connectivity_indexing_mode", "STATUS"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.thing_indexing_mode", "REGISTRY_AND_SHADOW"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.filter.0.named_shadow_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "thing_indexing_configuration.0.filter.0.named_shadow_names.0", "thing1shadow"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIndexingConfigurationConfig_basic = `
resource "aws_iot_indexing_configuration" "test" {
  thing_group_indexing_configuration {
    thing_group_indexing_mode = "OFF"
  }

  thing_indexing_configuration {
    thing_indexing_mode = "OFF"
  }
}
`

const testAccIndexingConfigurationConfig_allAttributes = `
resource "aws_iot_indexing_configuration" "test" {
  thing_group_indexing_configuration {
    thing_group_indexing_mode = "ON"
  }

  thing_indexing_configuration {
    thing_indexing_mode              = "REGISTRY_AND_SHADOW"
    thing_connectivity_indexing_mode = "STATUS"
    device_defender_indexing_mode    = "VIOLATIONS"
    named_shadow_indexing_mode       = "ON"

	filter {
      named_shadow_names = [ "thing1shadow" ]
    }

    custom_field {
      name = "attributes.version"
      type = "Number"
    }
    custom_field {
      name = "shadow.name.thing1shadow.desired.DefaultDesired"
      type = "String"
    }
    custom_field {
      name = "deviceDefender.securityProfile1.NUMBER_VALUE_BEHAVIOR.lastViolationValue.number"
      type = "Number"
    }
  }
}
`
