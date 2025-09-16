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
var _ datasource.DataSource = &MetadataDataSource{}

func NewMetadataDataSource() datasource.DataSource {
	return &MetadataDataSource{}
}

// MetadataDataSource defines the data source implementation.
type MetadataDataSource struct {
	client *client.Client
}

// MetadataDataSourceModel describes the data source data model.
type MetadataDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Path      types.String `tfsdk:"path"`
	Value     types.String `tfsdk:"value"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (d *MetadataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metadata"
}

func (d *MetadataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DirtCloud metadata data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Metadata identifier",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Metadata path identifier",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Metadata value",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Metadata creation timestamp",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Metadata last updated timestamp",
				Computed:            true,
			},
		},
	}
}

func (d *MetadataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MetadataDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the metadata from the API
	metadata, err := d.client.GetMetadataByPath(ctx, data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metadata, got error: %s", err))
		return
	}

	// Update the model with the API response data
	data.ID = types.StringValue(metadata.ID)
	data.Path = types.StringValue(metadata.Path)
	data.Value = types.StringValue(metadata.Value)
	data.CreatedAt = types.StringValue(metadata.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(metadata.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}