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
var _ resource.Resource = &BucketResource{}
var _ resource.ResourceWithImportState = &BucketResource{}

func NewBucketResource() resource.Resource {
	return &BucketResource{}
}

// BucketResource defines the resource implementation.
type BucketResource struct {
	client *client.Client
}

// BucketResourceModel describes the resource data model.
type BucketResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *BucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket"
}

func (r *BucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DirtCloud bucket resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Bucket identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Bucket name (unique, max 255, /^[a-zA-Z0-9_-]+$/)",
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Bucket creation timestamp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Bucket last updated timestamp",
			},
		},
	}
}

func (r *BucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	if name == "" || len(name) > 255 {
		resp.Diagnostics.AddError("Invalid name", "Bucket name must be non-empty and at most 255 characters")
		return
	}

	createReq := client.CreateBucketRequest{ Name: name }
	bucket, err := r.client.CreateBucket(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create bucket, got error: %s", err))
		return
	}

	data.ID = types.StringValue(bucket.ID)
	data.Name = types.StringValue(bucket.Name)
	data.CreatedAt = types.StringValue(bucket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(bucket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	bucket, err := r.client.GetBucket(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read bucket, got error: %s", err))
		return
	}

	data.Name = types.StringValue(bucket.Name)
	data.CreatedAt = types.StringValue(bucket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(bucket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	if name == "" || len(name) > 255 {
		resp.Diagnostics.AddError("Invalid name", "Bucket name must be non-empty and at most 255 characters")
		return
	}

	updateReq := client.UpdateBucketRequest{ Name: name }
	bucket, err := r.client.UpdateBucket(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update bucket, got error: %s", err))
		return
	}

	data.Name = types.StringValue(bucket.Name)
	data.UpdatedAt = types.StringValue(bucket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteBucket(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete bucket, got error: %s", err))
		return
	}
}

func (r *BucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data := BucketResourceModel{ ID: types.StringValue(req.ID) }

	bucket, err := r.client.GetBucket(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import bucket, got error: %s", err))
		return
	}

	data.Name = types.StringValue(bucket.Name)
	data.CreatedAt = types.StringValue(bucket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(bucket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
