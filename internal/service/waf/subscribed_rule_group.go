package waf

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	DSNameSubscribedRuleGroup = "Subscribed Rule Group Data Source"
)

func DataSourceSubscribedRuleGroup() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSubscribedRuleGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceSubscribedRuleGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).WAFConn
	name, nameOk := d.Get("name").(string)
	metricName, metricNameOk := d.Get("metric_name").(string)

	// Error out if string-assertion fails for either name or metricName
	if !nameOk || !metricNameOk {
		if !nameOk {
			name = DSNameSubscribedRuleGroup
		}

		err := errors.New("unable to read attributes")
		return names.DiagError(names.WAF, names.ErrActionReading, DSNameSubscribedRuleGroup, name, err)
	}

	output, err := FindSubscribedRuleGroupByNameOrMetricName(ctx, conn, name, metricName)

	if err != nil {
		return names.DiagError(names.WAF, names.ErrActionReading, DSNameSubscribedRuleGroup, name, err)
	}

	d.SetId(aws.StringValue(output.RuleGroupId))
	d.Set("metric_name", aws.StringValue(output.MetricName))
	d.Set("name", aws.StringValue(output.Name))

	return nil
}
