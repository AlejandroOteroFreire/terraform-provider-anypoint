package anypoint

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	auth "github.com/mulesoft-anypoint/anypoint-client-go/authorization"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_AUTH_TYPE", "connected_app"),
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "connected_app" && v != "user" {
						errs = append(errs, fmt.Errorf("%q must be 'connected_app' or 'user', got: %s", key, v))
					}
					return
				},
				Description: "Authentication type. Valid values: 'connected_app' (default, client credentials grant) or 'user' (password grant — uses connected app credentials to authenticate on behalf of a user, granting that user's permissions for operations like Access Management).",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_CLIENT_ID", nil),
				Description: "the connected app's id",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_CLIENT_SECRET", nil),
				Description: "the connected app's secret",
			},
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_ACCESS_TOKEN", nil),
				Description: "the connected app's access token",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_USERNAME", nil),
				Description: "Username for user authentication (only required when auth_type is 'user').",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_PASSWORD", nil),
				Description: "Password for user authentication (only required when auth_type is 'user').",
			},
			"cplane": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_CPLANE", "us"),
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "us" && v != "eu" && v != "gov" {
						errs = append(errs, fmt.Errorf("%q must be 'eu', 'us' or 'gov', got: %s", key, v))
					}
					return
				},
				Description: "the anypoint control plane",
			},
		},
		ResourcesMap:         RESOURCES_MAP,
		DataSourcesMap:       DATASOURCES_MAP,
		ConfigureContextFunc: providerConfigure,
		TerraformVersion:     "v1.0.1",
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	var diags diag.Diagnostics

	auth_type := d.Get("auth_type").(string)
	client_id := d.Get("client_id").(string)
	client_secret := d.Get("client_secret").(string)
	access_token := d.Get("access_token").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	cplane := d.Get("cplane").(string)

	server_index := cplane2serverindex(cplane)
	auth_ctx := context.WithValue(ctx, auth.ContextServerIndex, server_index)

	// Pre-signed access token takes precedence regardless of auth_type
	if access_token != "" {
		return newProviderConfOutput(access_token, server_index), diags
	}

	switch auth_type {
	case "user":
		// OAuth2 password grant — uses connected app + user credentials to authenticate on behalf of the user.
		if client_id == "" || client_secret == "" || username == "" || password == "" {
			return newProviderConfOutput("", server_index), append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Incomplete user authentication configuration",
				Detail:   "When auth_type is 'user', all of client_id, client_secret, username and password must be provided.",
			})
		}
		authres, d := userOAuth2Auth(auth_ctx, client_id, client_secret, username, password)
		if d != nil {
			return newProviderConfOutput("", server_index), d
		}
		return newProviderConfOutput(authres.GetAccessToken(), server_index), diags

	case "", "connected_app":
		// Default — OAuth2 client_credentials grant
		if client_id == "" || client_secret == "" {
			return newProviderConfOutput("", server_index), append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing connected app credentials",
				Detail:   "client_id and client_secret must be provided for connected app authentication.",
			})
		}
		authres, d := connectedAppAuth(auth_ctx, client_id, client_secret)
		if d != nil {
			return newProviderConfOutput("", server_index), d
		}
		return newProviderConfOutput(authres.GetAccessToken(), server_index), diags

	default:
		return newProviderConfOutput("", server_index), append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid auth_type",
			Detail:   fmt.Sprintf("auth_type must be 'connected_app' or 'user', got: %s", auth_type),
		})
	}
}

/*
Authenticates a connected app on behalf of a user using OAuth2 password grant.
This is the "admin" mode — combines connected app credentials with user credentials.
The resulting token inherits the user's permissions, which is required for operations
like Access Management (teams, team_roles, etc.) that the connected app alone cannot perform.
*/
func userOAuth2Auth(ctx context.Context, client_id, client_secret, username, password string) (*auth.ApiV2Oauth2TokenPost200Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	creds := auth.NewCredentials()
	grantType := "password"
	creds.GrantType = &grantType
	creds.SetClientId(client_id)
	creds.SetClientSecret(client_secret)
	creds.SetUsername(username)
	creds.SetPassword(password)

	cfgauth := auth.NewConfiguration()
	authclient := auth.NewAPIClient(cfgauth)
	authres, httpr, err := authclient.DefaultApi.ApiV2Oauth2TokenPost(ctx).Credentials(*creds).Execute()
	if err != nil {
		var details string
		if httpr != nil {
			b, _ := io.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to authenticate on behalf of user (password grant)",
			Detail:   details,
		})
		return auth.NewApiV2Oauth2TokenPost200Response(), diags
	}
	defer httpr.Body.Close()
	return authres, diags
}

/*
Authenticates a connected app using OAuth2 client_credentials grant.
*/
func connectedAppAuth(ctx context.Context, client_id string, client_secret string) (*auth.ApiV2Oauth2TokenPost200Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	creds := auth.NewCredentialsWithDefaults()
	creds.SetClientId(client_id)
	creds.SetClientSecret(client_secret)
	cfgauth := auth.NewConfiguration()
	authclient := auth.NewAPIClient(cfgauth)
	authres, httpr, err := authclient.DefaultApi.ApiV2Oauth2TokenPost(ctx).Credentials(*creds).Execute()
	if err != nil {
		var details string
		if httpr != nil {
			b, _ := io.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to authenticate using connected app",
			Detail:   details,
		})
		return auth.NewApiV2Oauth2TokenPost200Response(), diags
	}
	defer httpr.Body.Close()
	return authres, diags
}

/*
returns the server index depending on the control plane name
if the control plane is not recognized, returns -1
*/
func cplane2serverindex(cplane string) int {
	switch cplane {
	case "eu":
		return 1
	case "us":
		return 0
	case "gov":
		return 2
	}
	return -1
}
