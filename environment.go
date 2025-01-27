package runscope

import (
	"encoding/json"
	"fmt"
	"time"
)

// Environment stores details for shared and test-specific environments. See https://www.runscope.com/docs/api/environments
type Environment struct {
	ID                  string                    `json:"id,omitempty"`
	Name                string                    `json:"name,omitempty"`
	Script              string                    `json:"script,omitempty"`
	PreserveCookies     bool                      `json:"preserve_cookies"`
	TestID              string                    `json:"test_id,omitempty"`
	InitialVariables    map[string]string         `json:"initial_variables,omitempty"`
	Integrations        []*EnvironmentIntegration `json:"integrations,omitempty"`
	Regions             []string                  `json:"regions,omitempty"`
	VerifySsl           bool                      `json:"verify_ssl"`
	ExportedAt          *time.Time                `json:"exported_at,omitempty"`
	RetryOnFailure      bool                      `json:"retry_on_failure"`
	RemoteAgents        []*LocalMachine           `json:"remote_agents,omitempty"`
	WebHooks            []string                  `json:"webhooks,omitempty"`
	ParentEnvironmentID string                    `json:"parent_environment_id,omitempty"`
	EmailSettings       *EmailSettings            `json:"emails,omitempty"`
	ClientCertificate   string                    `json:"client_certificate,omitempty"`
	Headers             map[string][]string       `json:"headers,omitempty"`
}

// EmailSettings determining how test failures trigger notifications
type EmailSettings struct {
	NotifyAll       bool       `json:"notify_all"`
	NotifyOn        string     `json:"notify_on"`
	NotifyThreshold int        `json:"notify_threshold"`
	Recipients      []*Contact `json:"recipients"`
}

// EnvironmentIntegration represents an integration with a third-party. See https://www.runscope.com/docs/api/integrations
type EnvironmentIntegration struct {
	ID              string `json:"id"`
	IntegrationType string `json:"integration_type"`
	Description     string `json:"description,omitempty"`
}

// LocalMachine used in an environment to represent a remote agent
type LocalMachine struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return new(Environment)
}

// CreateSharedEnvironment creates a new shared environment. See https://www.runscope.com/docs/api/environments#create-shared
func (client *Client) CreateSharedEnvironment(environment *Environment, bucket *Bucket) (*Environment, error) {
	return client.createEnvironment(environment, fmt.Sprintf("/buckets/%s/environments", bucket.Key))
}

// CreateTestEnvironment creates a new test environment. See https://www.runscope.com/docs/api/environments#create
func (client *Client) CreateTestEnvironment(environment *Environment, test *Test) (*Environment, error) {
	return client.createEnvironment(environment, fmt.Sprintf("/buckets/%s/tests/%s/environments",
		test.Bucket.Key, test.ID))
}

// ListSharedEnvironment lists all shared environments for a given bucket. See https://www.runscope.com/docs/api/environments#list-shared
func (client *Client) ListSharedEnvironment(bucket *Bucket) ([]*Environment, error) {
	return client.listEnvironments(bucket, fmt.Sprintf("/buckets/%s/environments", bucket.Key))
}

// ListTestEnvironment lists all tests environments in a given test. See https://api.blazemeter.com/api-monitoring/#test-envrionment-list
func (client *Client) ListTestEnvironment(bucket *Bucket, test *Test) ([]*Environment, error) {
	return client.listEnvironments(bucket, fmt.Sprintf("/buckets/%s/tests/%s/environments", bucket.Key, test.ID))
}

// ReadSharedEnvironment lists details about an existing shared environment. See https://www.runscope.com/docs/api/environments#detail
func (client *Client) ReadSharedEnvironment(environment *Environment, bucket *Bucket) (*Environment, error) {
	return client.readEnvironment(environment, fmt.Sprintf("/buckets/%s/environments/%s",
		bucket.Key, environment.ID))
}

// ReadTestEnvironment lists details about an existing test environment. See https://www.runscope.com/docs/api/environments#detail
func (client *Client) ReadTestEnvironment(environment *Environment, test *Test) (*Environment, error) {
	return client.readEnvironment(environment, fmt.Sprintf("/buckets/%s/tests/%s/environments/%s",
		test.Bucket.Key, test.ID, environment.ID))
}

// UpdateSharedEnvironment updates details about an existing shared environment. See https://www.runscope.com/docs/api/environments#modify
func (client *Client) UpdateSharedEnvironment(environment *Environment, bucket *Bucket) (*Environment, error) {
	return client.updateEnvironment(environment,
		fmt.Sprintf("/buckets/%s/environments/%s", bucket.Key, environment.ID))
}

// UpdateTestEnvironment updates details about an existing test environment. See https://www.runscope.com/docs/api/environments#modify
func (client *Client) UpdateTestEnvironment(environment *Environment, test *Test) (*Environment, error) {
	return client.updateEnvironment(environment,
		fmt.Sprintf("/buckets/%s/tests/%s/environments/%s", test.Bucket.Key, test.ID, environment.ID))
}

// DeleteEnvironment deletes an existing shared environment. https://www.runscope.com/docs/api/environments#delete
func (client *Client) DeleteEnvironment(environment *Environment, bucket *Bucket) error {
	return client.deleteResource("environment", environment.ID,
		fmt.Sprintf("/buckets/%s/environments/%s", bucket.Key, environment.ID))
}

func (environment *Environment) String() string {
	value, err := json.Marshal(environment)
	if err != nil {
		return ""
	}

	return string(value)
}

func (client *Client) createEnvironment(environment *Environment, endpoint string) (*Environment, error) {
	newResource, error := client.createResource(environment, "environment", environment.Name, endpoint)
	if error != nil {
		return nil, error
	}

	newEnvironment, error := getEnvironmentFromResponse(newResource.Data)
	if error != nil {
		return nil, error
	}

	return newEnvironment, nil
}

func (client *Client) listEnvironments(bucket *Bucket, endpoint string) ([]*Environment, error) {
	resource, error := client.readResource("environments", bucket.Key, endpoint)
	if error != nil {
		return nil, error
	}

	list, error := getEnvironmentsFromResponse(resource.Data)
	if error != nil {
		return nil, error
	}

	return list, nil
}

func (client *Client) readEnvironment(environment *Environment, endpoint string) (*Environment, error) {
	resource, error := client.readResource("environment", environment.ID, endpoint)
	if error != nil {
		return nil, error
	}

	readEnvironment, error := getEnvironmentFromResponse(resource.Data)
	if error != nil {
		return nil, error
	}

	return readEnvironment, nil
}

func (client *Client) updateEnvironment(environment *Environment, endpoint string) (*Environment, error) {
	resource, error := client.updateResource(environment, "environment", environment.ID, endpoint)
	if error != nil {
		return nil, error
	}

	updatedEnvironment, error := getEnvironmentFromResponse(resource.Data)
	if error != nil {
		return nil, error
	}

	return updatedEnvironment, nil
}

func getEnvironmentFromResponse(response interface{}) (*Environment, error) {
	environment := new(Environment)
	err := decode(environment, response)
	return environment, err
}

func getEnvironmentsFromResponse(response interface{}) ([]*Environment, error) {
	var environments []*Environment
	err := decode(&environments, response)
	return environments, err
}
