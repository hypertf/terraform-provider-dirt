// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-provider-dirt/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ObjectResource{}
var _ resource.ResourceWithImportState = &ObjectResource{}

func NewObjectResource() resource.Resource {
	return &ObjectResource{}
}

// ObjectResource defines the resource implementation.
type ObjectResource struct {
	client *client.Client
}

// ObjectResourceModel describes the resource data model.
type ObjectResourceModel struct {
	ID            types.String `tfsdk:"id"`
	BucketID      types.String `tfsdk:"bucket_id"`
	Path          types.String `tfsdk:"path"`
	ContentBase64 types.String `tfsdk:"content_base64"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func (r *ObjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

func (r *ObjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DirtCloud object resource (stored under a bucket)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Object identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bucket_id": schema.StringAttribute{
				MarkdownDescription: "ID of the bucket this object belongs to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Object path (unique within bucket)",
				Required:            true,
			},
			"content_base64": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded content to store",
				Required:            true,
				Sensitive:           true,
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Object creation timestamp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Object last updated timestamp",
			},
		},
	}
}

func (r *ObjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	bucketID := data.BucketID.ValueString()
	path := data.Path.ValueString()
	contentB64 := data.ContentBase64.ValueString()
	if bucketID == "" {
		resp.Diagnostics.AddError("Invalid bucket_id", "bucket_id must be provided")
		return
	}
	if path == "" || len(path) > 1024 {
		resp.Diagnostics.AddError("Invalid path", "path must be non-empty and at most 1024 characters")
		return
	}
	if contentB64 == "" || !isLikelyBase64(contentB64) {
		resp.Diagnostics.AddError("Invalid content_base64", "content_base64 must be non-empty and base64-encoded")
		return
	}

	createReq := client.CreateObjectRequest{ BucketID: bucketID, Path: path, Content: contentB64 }
	obj, err := r.client.CreateObject(ctx, bucketID, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create object, got error: %s", err))
		return
	}

	data.ID = types.StringValue(obj.ID)
	data.BucketID = types.StringValue(obj.BucketID)
	data.Path = types.StringValue(obj.Path)
	data.ContentBase64 = types.StringValue(obj.Content)
	data.CreatedAt = types.StringValue(obj.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(obj.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := r.client.GetObject(ctx, data.BucketID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read object, got error: %s", err))
		return
	}

	data.Path = types.StringValue(obj.Path)
	data.ContentBase64 = types.StringValue(obj.Content)
	data.CreatedAt = types.StringValue(obj.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(obj.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	path := data.Path.ValueString()
	contentB64 := data.ContentBase64.ValueString()
	if path == "" || len(path) > 1024 {
		resp.Diagnostics.AddError("Invalid path", "path must be non-empty and at most 1024 characters")
		return
	}
	if contentB64 != "" && !isLikelyBase64(contentB64) {
		resp.Diagnostics.AddError("Invalid content_base64", "content_base64 must be base64-encoded if provided")
		return
	}

	updateReq := client.UpdateObjectRequest{}
	updateReq.Path = &path
	updateReq.Content = &contentB64

	obj, err := r.client.UpdateObject(ctx, data.BucketID.ValueString(), data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update object, got error: %s", err))
		return
	}

	data.Path = types.StringValue(obj.Path)
	data.ContentBase64 = types.StringValue(obj.Content)
	data.UpdatedAt = types.StringValue(obj.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteObject(ctx, data.BucketID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete object, got error: %s", err))
		return
	}
}

func (r *ObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect ID format: {bucket_id}/{object_id}
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: {bucket_id}/{object_id}")
		return
	}

	bucketID := parts[0]
	objectID := parts[1]

	obj, err := r.client.GetObject(ctx, bucketID, objectID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import object, got error: %s", err))
		return
	}

	data := ObjectResourceModel{
		ID:            types.StringValue(obj.ID),
		BucketID:      types.StringValue(obj.BucketID),
		Path:          types.StringValue(obj.Path),
		ContentBase64: types.StringValue(obj.Content),
		CreatedAt:     types.StringValue(obj.CreatedAt.Format("2006-01-02T15:04:05Z07:00")),
		UpdatedAt:     types.StringValue(obj.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// isLikelyBase64 is a lightweight check to reduce obvious mistakes without adding dependencies.
func isLikelyBase64(s string) bool {
	// very permissive: length multiple of 4 and only valid base64 charset
	if len(s) == 0 || len(s)%4 != 0 {
		return false
	}
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '+' || r == '/' || r == '=' || r == '\n' || r == '\r' {
			continue
		}
		return false
	}
	return true
}
