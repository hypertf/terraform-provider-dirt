// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-provider-dirt/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InstanceResource{}
var _ resource.ResourceWithImportState = &InstanceResource{}

func NewInstanceResource() resource.Resource {
	return &InstanceResource{}
}

// InstanceResource defines the resource implementation.
type InstanceResource struct {
	client *client.Client
}

// InstanceResourceModel describes the resource data model.
type InstanceResourceModel struct {
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

func (r *InstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (r *InstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DirtCloud instance resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Instance identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the project this instance belongs to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Instance name",
				Required:            true,
			},
			"cpu": schema.Int64Attribute{
				MarkdownDescription: "Number of CPU cores",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
			},
			"memory_mb": schema.Int64Attribute{
				MarkdownDescription: "Memory in MB",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2048),
			},
			"image": schema.StringAttribute{
				MarkdownDescription: "Instance image",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ubuntu:20.04"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Instance status (running, stopped)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("running"),
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Instance creation timestamp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Instance last updated timestamp",
			},
		},
	}
}

func (r *InstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InstanceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the instance
	createReq := client.CreateInstanceRequest{
		ProjectID: data.ProjectID.ValueString(),
		Name:      data.Name.ValueString(),
		CPU:       int(data.CPU.ValueInt64()),
		MemoryMB:  int(data.MemoryMB.ValueInt64()),
		Image:     data.Image.ValueString(),
		Status:    data.Status.ValueString(),
	}

	instance, err := r.client.CreateInstance(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create instance, got error: %s", err))
		return
	}

	// Update the model with the response data
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

func (r *InstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InstanceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the instance from the API
	instance, err := r.client.GetInstance(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance, got error: %s", err))
		return
	}

	// Update the model with the latest data
	data.ProjectID = types.StringValue(instance.ProjectID)
	data.Name = types.StringValue(instance.Name)
	data.CPU = types.Int64Value(int64(instance.CPU))
	data.MemoryMB = types.Int64Value(int64(instance.MemoryMB))
	data.Image = types.StringValue(instance.Image)
	data.Status = types.StringValue(instance.Status)
	data.CreatedAt = types.StringValue(instance.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(instance.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InstanceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the instance - send all fields (server handles immutability)
	updateReq := client.UpdateInstanceRequest{}

	name := data.Name.ValueString()
	updateReq.Name = &name

	cpu := int(data.CPU.ValueInt64())
	updateReq.CPU = &cpu

	memoryMB := int(data.MemoryMB.ValueInt64())
	updateReq.MemoryMB = &memoryMB

	image := data.Image.ValueString()
	updateReq.Image = &image

	status := data.Status.ValueString()
	updateReq.Status = &status

	instance, err := r.client.UpdateInstance(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update instance, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Name = types.StringValue(instance.Name)
	data.CPU = types.Int64Value(int64(instance.CPU))
	data.MemoryMB = types.Int64Value(int64(instance.MemoryMB))
	data.Image = types.StringValue(instance.Image)
	data.Status = types.StringValue(instance.Status)
	data.UpdatedAt = types.StringValue(instance.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InstanceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the instance
	err := r.client.DeleteInstance(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete instance, got error: %s", err))
		return
	}
}

func (r *InstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID from the import request
	data := InstanceResourceModel{
		ID: types.StringValue(req.ID),
	}

	// Read the instance to populate other fields
	instance, err := r.client.GetInstance(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import instance, got error: %s", err))
		return
	}

	data.ProjectID = types.StringValue(instance.ProjectID)
	data.Name = types.StringValue(instance.Name)
	data.CPU = types.Int64Value(int64(instance.CPU))
	data.MemoryMB = types.Int64Value(int64(instance.MemoryMB))
	data.Image = types.StringValue(instance.Image)
	data.Status = types.StringValue(instance.Status)
	data.CreatedAt = types.StringValue(instance.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(instance.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Save imported data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
