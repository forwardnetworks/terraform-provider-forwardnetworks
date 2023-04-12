package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-forwardnetworks/forwardnetworks"
)

func resourceCollectionSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectionScheduleCreate,
		ReadContext:   resourceCollectionScheduleRead,
		UpdateContext: resourceCollectionScheduleUpdate,
		DeleteContext: resourceCollectionScheduleDelete,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Description: "The ID of the network to which the collection schedule belongs.",
				Required:    true,
				ForceNew:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the collection schedule.",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "The status of the collection schedule.",
				Required:    true,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Description: "The time zone of the collection schedule.",
				Required:    true,
			},
			"days_of_the_week": {
				Type:        schema.TypeList,
				Description: "The days of the week for the collection schedule.",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"times": {
				Type:        schema.TypeList,
				Description: "The times for the collection schedule.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"period_in_seconds": {
				Type:        schema.TypeInt,
				Description: "The period in seconds for the collection schedule.",
				Optional:    true,
			},
			"start_at": {
				Type:        schema.TypeString,
				Description: "The start time of the collection schedule.",
				Optional:    true,
			},
			"end_at": {
				Type:        schema.TypeString,
				Description: "The end time of the collection schedule.",
				Optional:    true,
			},
		},
	}
}

func resourceCollectionScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)

    schedule := &forwardnetworks.CollectionSchedule{
        Enabled:         d.Get("enabled").(bool),
        TimeZone:        d.Get("time_zone").(string),
        DaysOfTheWeek:   convertInterfaceSliceToIntSlice(d.Get("days_of_the_week").([]interface{})),
        Times:           convertInterfaceSliceToStringSlice(d.Get("times").([]interface{})),
        PeriodInSeconds: d.Get("period_in_seconds").(int),
        StartAt:         d.Get("start_at").(string),
        EndAt:           d.Get("end_at").(string),
    }

    createdSchedule, err := client.CreateCollectionSchedule(networkID, schedule)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(createdSchedule.Id)

    return resourceCollectionScheduleRead(ctx, d, m)
}

func resourceCollectionScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)

    scheduleID := d.Id()

    schedule, err := client.GetCollectionSchedule(networkID, scheduleID)
    if err != nil {
        return diag.FromErr(err)
    }

    d.Set("enabled", schedule.Enabled)
    d.Set("time_zone", schedule.TimeZone)
    d.Set("days_of_the_week", schedule.DaysOfTheWeek)
    d.Set("times", schedule.Times)
    d.Set("period_in_seconds", schedule.PeriodInSeconds)
    d.Set("start_at", schedule.StartAt)
    d.Set("end_at", schedule.EndAt)

    return nil
}

func resourceCollectionScheduleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)
    scheduleID := d.Id()

    schedule := &forwardnetworks.CollectionSchedule{
        Id:              scheduleID,
        Enabled:         d.Get("enabled").(bool),
        TimeZone:        d.Get("time_zone").(string),
        DaysOfTheWeek:   convertInterfaceSliceToIntSlice(d.Get("days_of_the_week").([]interface{})),
        Times:           convertInterfaceSliceToStringSlice(d.Get("times").([]interface{})),
        PeriodInSeconds: d.Get("period_in_seconds").(int),
        StartAt:         d.Get("start_at").(string),
        EndAt:           d.Get("end_at").(string),
    }

    _, err := client.UpdateCollectionSchedule(networkID, scheduleID, schedule)
    if err != nil {
        return diag.FromErr(err)
    }

    return resourceCollectionScheduleRead(ctx, d, m)
}

func resourceCollectionScheduleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)
    scheduleID := d.Id()

    err := client.DeleteCollectionSchedule(networkID, scheduleID)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId("")

    return nil
}

// Helper functions to convert slices of interface{} to slices of specific types
func convertInterfaceSliceToIntSlice(input []interface{}) []int {
	output := make([]int, len(input))
	for i, v := range input {
		output[i] = v.(int)
	}
	return output
}

func convertInterfaceSliceToStringSlice(input []interface{}) []string {
    output := make([]string, len(input))
    for i, v := range input {
        output[i] = v.(string)
    }
    return output
}