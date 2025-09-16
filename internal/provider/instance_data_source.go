// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-provider-dirt/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &InstanceDataSource{}

func NewInstanceDataSource() datasource.DataSource {
	return &InstanceDataSource{}
}

// InstanceDataSource defines the data source implementation.
type InstanceDataSource struct {
	client *client.Client
}

// InstanceDataSourceModel describes the data source data model.
type InstanceDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	CPU       types.Int64  `tfsdk:"cpu"`
	MemoryMB  types.Int64  `tfsdk:"memory_mb"`
	Image     types.String `tfsdk:"image"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (d *InstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (d *InstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DirtCloud instance data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Instance identifier",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the project this instance belongs to",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Instance name",
				Computed:            true,
			},
			"cpu": schema.Int64Attribute{
				MarkdownDescription: "Number of CPU cores",
				Computed:            true,
			},
			"memory_mb": schema.Int64Attribute{
				MarkdownDescription: "Memory in MB",
				Computed:            true,
			},
			"image": schema.StringAttribute{
				MarkdownDescription: "Instance image",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Instance status (running, stopped)",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Instance creation timestamp",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Instance last updated timestamp",
				Computed:            true,
			},
		},
	}
}

func (d *InstanceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *InstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InstanceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the instance from the API
	instance, err := d.client.GetInstance(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance, got error: %s", err))
		return
	}

	// Update the model with the API response data
	data.ID = types.StringValue(instance.ID)
	data.ProjectID = types.StringValue(instance.ProjectID)
	data.Name = types.StringValue(instance.Name)
	data.CPU = types.Int64Value(int64(instance.CPU))
	data.MemoryMB = types.Int64Value(int64(instance.MemoryMB))
	data.Image = types.StringValue(instance.Image)
	data.Status = types.StringValue(instance.Status)
	data.CreatedAt = types.StringValue(instance.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(instance.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}