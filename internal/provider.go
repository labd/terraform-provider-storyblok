package internal

import (
	"context"
	"github.com/hashicorp/go-retryablehttp"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
	"log"
	"net/http"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &storyblokProvider{}
)

type OptionFunc func(p *storyblokProvider)

func WithRetryableClient(retries int) OptionFunc {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = retries

	return func(p *storyblokProvider) {
		p.httpClient = retryClient.StandardClient()
	}
}

func WithDebugClient() OptionFunc {
	return func(p *storyblokProvider) {
		p.httpClient.Transport = debugTransport
	}
}

func WithRecorderClient(file string, mode recorder.Mode) (OptionFunc, func() error) {
	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       file,
		Mode:               mode,
		SkipRequestLatency: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	//Strip all fields we are not interested in
	hook := func(i *cassette.Interaction) error {
		i.Response.Headers = cleanHeaders(i.Response.Headers, "Content-Type")
		i.Request.Headers = cleanHeaders(i.Request.Headers)
		return nil
	}
	r.AddHook(hook, recorder.AfterCaptureHook)

	stop := func() error {
		return r.Stop()
	}

	return func(p *storyblokProvider) {
		p.httpClient = r.GetDefaultClient()
	}, stop
}

// New is a helper function to simplify provider server and testing implementation.
func New(opts ...OptionFunc) provider.Provider {
	var p = &storyblokProvider{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// storyblokProvider is the provider implementation.
type storyblokProvider struct {
	httpClient *http.Client
}

// storyblokProviderModel maps provider schema data to a Go type.
type storyblokProviderModel struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *storyblokProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "storyblok"
}

// Schema defines the provider-level schema for configuration data.
func (p *storyblokProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Storyblok.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "Management API base URL",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Personal access token",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Storyblok API client for data sources and resources.
func (p *storyblokProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Storyblok client")

	// Retrieve provider data from configuration
	var config storyblokProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	url := os.Getenv("STORYBLOK_URL")
	token := os.Getenv("STORYBLOK_TOKEN")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if resp.Diagnostics.HasError() {
		return
	}

	if url == "" {
		url = "https://mapi.storyblok.com"
	}

	ctx = tflog.SetField(ctx, "storyblok_url", url)
	ctx = tflog.SetField(ctx, "storyblok_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "storyblok_token")

	tflog.Debug(ctx, "Creating Storyblok client")

	apiKeyProvider, err := securityprovider.NewSecurityProviderApiKey("header", "Authorization", token)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Storyblok API Client", err.Error())
	}

	// Create a new Storyblok client using the configuration values
	client, err := sbmgmt.NewClientWithResponses(
		url,
		sbmgmt.WithHTTPClient(p.httpClient),
		sbmgmt.WithRequestEditorFn(apiKeyProvider.Intercept))

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Storyblok API Client",
			"An unexpected error occurred when creating the Storyblok API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Storyblok Client Error: "+err.Error(),
		)
		return
	}

	// Make the Storyblok client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Storyblok client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *storyblokProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *storyblokProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewComponentResource,
		NewComponentGroupResource,
		NewSpaceRoleResource,
		NewAssetFolderResource,
	}
}
