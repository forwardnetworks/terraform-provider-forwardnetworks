package main

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceAWSPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAWSPolicyRead,

		Schema: map[string]*schema.Schema{
			"aws_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AWS policy document.",
			},
		},
	}
}

func dataSourceAWSPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	awsPolicy, err := client.GetAwsPolicy()
	if err != nil {
		return diag.FromErr(err)
	}

	awsPolicyJson, err := json.Marshal(awsPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("forwardnetworks_aws_policy")
	d.Set("aws_policy", string(awsPolicyJson))

	return nil
}
