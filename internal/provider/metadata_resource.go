// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-provider-dirt/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MetadataResource{}
var _ resource.ResourceWithImportState = &MetadataResource{}

func NewMetadataResource() resource.Resource {
	return &MetadataResource{}
}

// MetadataResource defines the resource implementation.
type MetadataResource struct {
	client *client.Client
}

// MetadataResourceModel describes the resource data model.
type MetadataResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Path      types.String `tfsdk:"path"`
	Value     types.String `tfsdk:"value"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *MetadataResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metadata"
}

func (r *MetadataResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DirtCloud metadata resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Metadata identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Metadata path identifier (must be unique)",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Metadata value",
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Metadata creation timestamp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Metadata last updated timestamp",
			},
		},
	}
}

func (r *MetadataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *MetadataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MetadataResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the metadata
	createReq := client.CreateMetadataRequest{
		Path:  data.Path.ValueString(),
		Value: data.Value.ValueString(),
	}

	metadata, err := r.client.CreateMetadata(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create metadata, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.ID = types.StringValue(metadata.ID)
	data.Path = types.StringValue(metadata.Path)
	data.Value = types.StringValue(metadata.Value)
	data.CreatedAt = types.StringValue(metadata.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(metadata.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MetadataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MetadataResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the metadata from the API
	metadata, err := r.client.GetMetadata(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metadata, got error: %s", err))
		return
	}

	// Update the model with the latest data
	data.Path = types.StringValue(metadata.Path)
	data.Value = types.StringValue(metadata.Value)
	data.CreatedAt = types.StringValue(metadata.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(metadata.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MetadataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MetadataResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the metadata
	updateReq := client.UpdateMetadataRequest{}

	path := data.Path.ValueString()
	updateReq.Path = &path

	value := data.Value.ValueString()
	updateReq.Value = &value

	metadata, err := r.client.UpdateMetadata(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update metadata, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Path = types.StringValue(metadata.Path)
	data.Value = types.StringValue(metadata.Value)
	data.UpdatedAt = types.StringValue(metadata.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MetadataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MetadataResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the metadata
	err := r.client.DeleteMetadata(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete metadata, got error: %s", err))
		return
	}
}

func (r *MetadataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID from the import request
	data := MetadataResourceModel{
		ID: types.StringValue(req.ID),
	}

	// Read the metadata to populate other fields
	metadata, err := r.client.GetMetadata(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import metadata, got error: %s", err))
		return
	}

	data.Path = types.StringValue(metadata.Path)
	data.Value = types.StringValue(metadata.Value)
	data.CreatedAt = types.StringValue(metadata.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(metadata.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save imported data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
