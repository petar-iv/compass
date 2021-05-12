/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"text/template"

	"github.com/kyma-incubator/compass/components/director/pkg/inputvalidation"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/pkg/errors"
)

type Mode string

// Resource is used to identify entities which can be part of a webhook's request data
type Resource interface {
	Sentinel()
	Template() map[string]interface{}
}

// RequestObject struct contains parts of request that might be needed for later processing of a Webhook request
type RequestObject struct {
	Application        Resource
	BundleInstanceAuth Resource
	Type               resource.Type
	TenantID           string
	ExternalTenantID   string
	Headers            map[string]string
	ApplicationLabels  map[string]interface{}
}

// ResponseObject struct contains parts of response that might be needed for later processing of Webhook response
type ResponseObject struct {
	Body    map[string]interface{}
	Headers map[string]string
}

type URL struct {
	Method *string `json:"method"`
	Path   *string `json:"path"`
}

// Response defines the schema for Webhook output templates
type Response struct {
	Location          *string `json:"location"`
	SuccessStatusCode *int    `json:"success_status_code"`
	GoneStatusCode    *int    `json:"gone_status_code"`
	Error             *string `json:"error"`
}

// ResponseStatus defines the schema for Webhook status templates when dealing with async webhooks
type ResponseStatus struct {
	Location                   *string `json:"location"`
	Status                     *string `json:"status"`
	SuccessStatusCode          *int    `json:"success_status_code"`
	SuccessStatusIdentifier    *string `json:"success_status_identifier"`
	InProgressStatusIdentifier *string `json:"in_progress_status_identifier"`
	FailedStatusIdentifier     *string `json:"failed_status_identifier"`
	Error                      *string `json:"error"`
}

func (u *URL) Validate() error {
	if u.Method == nil {
		return errors.New("missing URL Template method field")
	}

	if u.Path == nil {
		return errors.New("missing URL Template path field")
	}

	_, err := url.ParseRequestURI(*u.Path)
	if err != nil {
		return errors.Wrap(err, "failed to parse URL Template path field")
	}

	return nil
}

func (r *Response) Validate() error {
	if r.Location == nil {
		return errors.New("missing Output Template location field")
	}

	if r.SuccessStatusCode == nil {
		return errors.New("missing Output Template success status code field")
	}

	if r.Error == nil {
		return errors.New("missing Output Template error field")
	}

	return nil
}

func (rs *ResponseStatus) Validate() error {
	if rs.Status == nil {
		return errors.New("missing Status Template status field")
	}

	if rs.SuccessStatusCode == nil {
		return errors.New("missing Status Template success status code field")
	}

	if rs.SuccessStatusIdentifier == nil {
		return errors.New("missing Status Template success status identifier field")
	}

	if rs.InProgressStatusIdentifier == nil {
		return errors.New("missing Status Template in progress status identifier field")
	}

	if rs.FailedStatusIdentifier == nil {
		return errors.New("missing Status Template failed status identifier field")
	}

	if rs.Error == nil {
		return errors.New("missing Status Template error field")
	}

	return nil
}

func (rd *RequestObject) ParseURLTemplate(tmpl *string) (*URL, error) {
	var url URL
	return &url, parseTemplate(tmpl, *rd, &url)
}

func (rd *RequestObject) ParseInputTemplate(tmpl *string) ([]byte, error) {
	res := json.RawMessage{}
	temp := struct {
		Application        map[string]interface{}
		BundleInstanceAuth map[string]interface{}
		Type               resource.Type
		TenantID           string
		ExternalTenantID   string
		Headers            map[string]string
		ApplicationLabels map[string]interface{}
	}{
		Application:        rd.Application.Template(),
		BundleInstanceAuth: rd.BundleInstanceAuth.Template(),
		Type:               rd.Type,
		TenantID:           rd.TenantID,
		ExternalTenantID:   rd.ExternalTenantID,
		Headers:            rd.Headers,
		ApplicationLabels: rd.ApplicationLabels,
	}
	return res, parseTemplate(tmpl, temp, &res)
}

func (rd *RequestObject) ParseHeadersTemplate(tmpl *string) (http.Header, error) {
	var headers http.Header
	return headers, parseTemplate(tmpl, *rd, &headers)
}

func (rd *ResponseObject) ParseOutputTemplate(tmpl *string) (*Response, error) {
	var resp Response
	return &resp, parseTemplate(tmpl, *rd, &resp)
}

func (rd *ResponseObject) ParseStatusTemplate(tmpl *string) (*ResponseStatus, error) {
	var respStatus ResponseStatus
	return &respStatus, parseTemplate(tmpl, *rd, &respStatus)
}

type CredentialsResponse struct {
	SuccessStatusCode *int    `json:"success_status_code"`
	Error             *string `json:"error"`
	URL               *string `json:"url"`
	APIKey            *string `json:"api_key"`
	ClientID          *string `json:"client_id"`
	ClientSecret      *string `json:"client_secret"`
	Username          *string `json:"username"`
	Password          *string `json:"password"`
}

// TODO implement
func (cr *CredentialsResponse) Validate() error {
	return nil
}

func (rd *ResponseObject) ParseCredentialsResultTemplate(tmpl *string) (*CredentialsResponse, error) {
	var resp CredentialsResponse
	return &resp, parseTemplate(tmpl, *rd, &resp)
}

func parseTemplate(tmpl *string, data interface{}, dest interface{}) error {
	t, err := template.New("").Option("missingkey=zero").Parse(*tmpl)
	if err != nil {
		return err
	}

	res := new(bytes.Buffer)
	if err = t.Execute(res, data); err != nil {
		return err
	}

	if err = json.Unmarshal(res.Bytes(), dest); err != nil {
		return err
	}

	if validatable, ok := dest.(inputvalidation.Validatable); ok {
		return validatable.Validate()
	}

	return nil
}
