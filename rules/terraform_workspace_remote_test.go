package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformWorkspaceRemoteRule(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "terraform.workspace in resource with remote backend",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "remote" {}
}
resource "null_resource" "a" {
	triggers = {
		w = terraform.workspace
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformWorkspaceRemoteRule(),
					Message: "terraform.workspace should not be used with a 'remote' backend",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 8, Column: 7},
						End:      hcl.Pos{Line: 8, Column: 26},
					},
				},
			},
		},
		{
			Name: "terraform.workspace with no backend",
			Content: `
terraform {
	required_version = ">= 1.0"
}
resource "null_resource" "a" {
	triggers = {
		w = terraform.workspace
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "terraform.workspace with non-remote backend",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "local" {}
}
resource "null_resource" "a" {
	triggers = {
		w = terraform.workspace
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "terraform.workspace with remote backend, but required_version does not support 1.0.x",
			Content: `
terraform {
	required_version = ">= 1.1"
	backend "remote" {}
}
resource "null_resource" "a" {
	triggers = {
		w = terraform.workspace
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "terraform.workspace with remote backend, but required_version is not declared",
			Content: `
terraform {
	backend "remote" {}
}
resource "null_resource" "a" {
	triggers = {
		w = terraform.workspace
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "terraform.workspace without terraform setting",
			Content: `
resource "null_resource" "a" {
	triggers = {
		w = terraform.workspace
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "terraform.workspace in data source with remote backend",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "remote" {}
}
data "null_data_source" "a" {
	inputs = {
		w = terraform.workspace
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformWorkspaceRemoteRule(),
					Message: "terraform.workspace should not be used with a 'remote' backend",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 8, Column: 7},
						End:      hcl.Pos{Line: 8, Column: 26},
					},
				},
			},
		},
		{
			Name: "terraform.workspace in module call with remote backend",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "remote" {}
}
module "a" {
	source = "./module"
	w = terraform.workspace
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformWorkspaceRemoteRule(),
					Message: "terraform.workspace should not be used with a 'remote' backend",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 8, Column: 6},
						End:      hcl.Pos{Line: 8, Column: 25},
					},
				},
			},
		},
		{
			Name: "terraform.workspace in provider config with remote backend",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "remote" {}
}
provider "aws" {
	assume_role	{
		role_arn = terraform.workspace
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformWorkspaceRemoteRule(),
					Message: "terraform.workspace should not be used with a 'remote' backend",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 8, Column: 14},
						End:      hcl.Pos{Line: 8, Column: 33},
					},
				},
			},
		},
		{
			Name: "terraform.workspace in locals with remote backend",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "remote" {}
}
locals {
	w = terraform.workspace
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformWorkspaceRemoteRule(),
					Message: "terraform.workspace should not be used with a 'remote' backend",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 7, Column: 6},
						End:      hcl.Pos{Line: 7, Column: 25},
					},
				},
			},
		},
		{
			Name: "terraform.workspace in output with remote backend",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "remote" {}
}
output "o" {
	value = terraform.workspace
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformWorkspaceRemoteRule(),
					Message: "terraform.workspace should not be used with a 'remote' backend",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 7, Column: 10},
						End:      hcl.Pos{Line: 7, Column: 29},
					},
				},
			},
		},
		{
			Name: "nonmatching expressions with remote backend",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "remote" {}
}
locals {
	a = "terraform.workspace"
	b = path.module
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "meta-arguments",
			Content: `
terraform {
	required_version = ">= 1.0"
	backend "remote" {}
}
resource "aws_instance" "foo" {
  instance_type = "t3.nano"
  lifecycle {
    ignore_changes = [instance_type]
  }
  providers = {
    aws = aws
  }
  depends_on = [aws_instance.bar]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "terraform.workspace in JSON syntax",
			JSON: true,
			Content: `
{
  "terraform": {
	"required_version": ">= 1.0",
    "backend": {
      "remote": {}
	}
  },
  "resource": {
    "null_resource": {
      "a": {
        "triggers": {
          "w": "${terraform.workspace}"
        }
      }
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformWorkspaceRemoteRule(),
					Message: "terraform.workspace should not be used with a 'remote' backend",
					Range: hcl.Range{
						Filename: "config.tf.json",
						Start:    hcl.Pos{Line: 9, Column: 15},
						End:      hcl.Pos{Line: 17, Column: 4},
					},
				},
			},
		},
		{
			Name: "with ignore_changes",
			Content: `
terraform {
  required_version = ">= 1.0"
  backend "remote" {}
}
resource "kubernetes_secret" "my_secret" {
  data = {}
  lifecycle {
    ignore_changes = [
      data
    ]
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "required_version does not support 1.0.x by multiple setting",
			Content: `
terraform {
	required_version = "< 1.9"
	backend "remote" {}
}
terraform {
	required_version = ">= 1.1"
}
resource "null_resource" "a" {
	triggers = {
		w = terraform.workspace
	}
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewTerraformWorkspaceRemoteRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "config.tf"
			if tc.JSON {
				filename = "config.tf.json"
			}
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
