// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-provider-dirt/internal/client"
)

// Ensure DirtProvider satisfies various provider interfaces.
var _ provider.Provider = &DirtProvider{}

// DirtProvider defines the provider implementation.
type DirtProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DirtProviderModel describes the provider data model.
type DirtProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *DirtProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dirt"
	resp.Version = p.version
}

func (p *DirtProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DirtCloud is a fake local cloud provider for learning and testing Terraform. It does not provision any real infrastructure. Instead, it simulates resources (projects, instances, metadata) and is paired with a local console that looks and behaves like a real cloud so you can practice Terraform workflows safely.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The DirtCloud API endpoint. Defaults to http://localhost:8080/v1. Can also be set via the DIRT_ENDPOINT environment variable.",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The DirtCloud API token for authentication. Can also be set via the DIRT_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *DirtProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DirtProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set defaults and get values from environment if not configured
	endpoint := "http://localhost:8080/v1"
	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	} else if envEndpoint := os.Getenv("DIRT_ENDPOINT"); envEndpoint != "" {
		endpoint = envEndpoint
	}

	token := ""
	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	} else if envToken := os.Getenv("DIRT_TOKEN"); envToken != "" {
		token = envToken
	}

	// Create DirtCloud client
	dirtClient := client.NewClient(endpoint)
	if token != "" {
		dirtClient.Token = token
	}

	// Make the client available to resources and data sources
	resp.DataSourceData = dirtClient
	resp.ResourceData = dirtClient
}

func (p *DirtProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewInstanceResource,
		NewMetadataResource,
		NewBucketResource,
		NewObjectResource,
	}
}

func (p *DirtProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		// No ephemeral resources for now
	}
}

func (p *DirtProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewInstanceDataSource,
		NewMetadataDataSource,
	}
}

func (p *DirtProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// No functions for now
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DirtProvider{
			version: version,
		}
	}
}
