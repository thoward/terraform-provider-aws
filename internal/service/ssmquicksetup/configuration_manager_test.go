// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssmquicksetup_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/service/ssmquicksetup"
	"github.com/aws/aws-sdk-go-v2/service/ssmquicksetup/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tfssmquicksetup "github.com/hashicorp/terraform-provider-aws/internal/service/ssmquicksetup"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSSMQuickSetupConfigurationManager_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var cm ssmquicksetup.GetConfigurationManagerOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_ssmquicksetup_configuration_manager.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.SSMQuickSetupEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMQuickSetupServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationManagerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationManagerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationManagerExists(ctx, resourceName, &cm),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "maintenance_window_start_time.0.day_of_week"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "user.*", map[string]string{
						"console_access": "false",
						"groups.#":       "0",
						"username":       "Test",
						"password":       "TestTest1234",
					}),
					acctest.MatchResourceAttrRegionalARN(resourceName, "manager_arn", "ssmquicksetup", regexache.MustCompile(`configurationmanager:+.`)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apply_immediately", "user"},
			},
		},
	})
}

func TestAccSSMQuickSetupConfigurationManager_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var cm ssmquicksetup.GetConfigurationManagerOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_ssmquicksetup_configuration_manager.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.SSMQuickSetupEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMQuickSetupServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationManagerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationManagerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationManagerExists(ctx, resourceName, &cm),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfssmquicksetup.ResourceConfigurationManager, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckConfigurationManagerDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).SSMQuickSetupClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_ssmquicksetup_configuration_manager" {
				continue
			}
			managerARN := rs.Primary.Attributes["manager_arn"]

			_, err := tfssmquicksetup.FindConfigurationManagerByID(ctx, conn, managerARN)
			if errs.IsA[*types.ResourceNotFoundException](err) {
				return nil
			}
			if err != nil {
				return create.Error(names.SSMQuickSetup, create.ErrActionCheckingDestroyed, tfssmquicksetup.ResNameConfigurationManager, managerARN, err)
			}

			return create.Error(names.SSMQuickSetup, create.ErrActionCheckingDestroyed, tfssmquicksetup.ResNameConfigurationManager, managerARN, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckConfigurationManagerExists(ctx context.Context, name string, configurationmanager *ssmquicksetup.GetConfigurationManagerOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.SSMQuickSetup, create.ErrActionCheckingExistence, tfssmquicksetup.ResNameConfigurationManager, name, errors.New("not found"))
		}

		managerARN := rs.Primary.Attributes["manager_arn"]
		if managerARN == "" {
			return create.Error(names.SSMQuickSetup, create.ErrActionCheckingExistence, tfssmquicksetup.ResNameConfigurationManager, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SSMQuickSetupClient(ctx)

		out, err := tfssmquicksetup.FindConfigurationManagerByID(ctx, conn, managerARN)
		if err != nil {
			return create.Error(names.SSMQuickSetup, create.ErrActionCheckingExistence, tfssmquicksetup.ResNameConfigurationManager, managerARN, err)
		}

		*configurationmanager = *out

		return nil
	}
}

func testAccConfigurationManagerConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssmquicksetup_configuration_manager" "test" {
  configuration_manager_name = %[1]q
}
`, rName)
}
