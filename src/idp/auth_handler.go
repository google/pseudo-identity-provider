// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package idp

import (
	"customidp/config"
	sessionmgmt "customidp/session"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// authHandler takes in a oauth2 auth request and redirects a target with either passed-through or replaced parameters.
// It can also return errors, or time out the request.
func authHandler(w http.ResponseWriter, r *http.Request) {
	action := config.GetGlobalConfig().AuthAction
	input := getInputData(r)
	addRequestLogEntry(input, action.Action)

	switch action.Action {
	case "redirect":
		authRedirect(w, r, input)
	case "error":
		errorResponse(w, r, &action.Error)
	case "block":
		blockResponse(w)
	}
}

// authRedirect creates a http redirect based on configuration.
func authRedirect(w http.ResponseWriter, r *http.Request, input *sessionmgmt.RequestInput) {
	redirect := config.GetGlobalConfig().AuthAction.Redirect
	paramsMap := make(map[string]config.Parameter)
	for _, param := range redirect.Parameters {
		paramsMap[param.ID] = param
	}

	// For each input request parameter check configured param action.
	// For Authz endpoint we only need to look at URL parameters.
	redirectParams, err := getQueryParams(input, paramsMap, redirect.DefaultParamAction, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURI, err := getRedirectURI(input, &redirect)
	if err != nil {
		http.Error(w, fmt.Sprintf("No redirect_uri present %v", err), http.StatusBadRequest)
		return
	}

	reqURI := redirectURI
	if redirect.UseHashFragment {
		reqURI += "#" + redirectParams.Encode()
	} else {
		reqURI += "?" + redirectParams.Encode()
	}

	sessionmgmt.CreateSession(input, redirectParams)

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	http.Redirect(w, r, reqURI, http.StatusFound)
}

// getRedirectURI gets a set or custom redirectURI if specified, otherwise
// uses the input parameter.
func getRedirectURI(input *sessionmgmt.RequestInput, c *config.AuthRedirect) (string, error) {
	if c.RedirectTarget.UseCustomRedirectURI {
		if c.RedirectTarget.Target != "" {
			return c.RedirectTarget.Target, nil
		} else if c.RedirectTarget.CustomKey != "" {
			// Construct a Paramter representing the Custom  value.
			p := config.Parameter{
				Action:    "custom",
				CustomKey: c.RedirectTarget.CustomKey,
			}
			res, err := p.Get(input)
			if err != nil {
				return "", err
			}

			if len(res) == 0 {
				return "", errors.New("custom function did not return any values")
			}
			return res[0], nil
		}
		return "", errors.New("missing custom redirect config")
	}

	// Otherwise use input parameter.
	uriVals := input.URLParams["redirect_uri"]
	if len(uriVals) == 0 {
		return "", errors.New("missing redirect_uri")
	}

	targetURI, err := url.QueryUnescape(uriVals[0])
	if err != nil {
		// Normally we would error here, but maybe they are trying something clever.
		// so let it through.
		targetURI = uriVals[0]
	}
	return targetURI, nil
}

// getQueryParams builds the redirection query parameters from config.
// We first handle input params against config, then configured params
// that are missing from input are handled.
func getQueryParams(input *sessionmgmt.RequestInput, paramConfig map[string]config.Parameter, defaultAction string, r *http.Request) (url.Values, error) {
	redirectParams := url.Values{}
	for id := range input.URLParams {
		configParam, ok := paramConfig[id]
		if !ok {
			// This parameter is not configured.
			// Get the default parameter value based on configured default action.
			p := config.Parameter{
				ID:     id,
				Action: defaultAction,
			}
			redirectParams[id] = p.GetDefaultValue(input)
			continue
		}

		vals, err := configParam.Get(input)
		if err != nil {
			return nil, err
		}

		if len(vals) != 0 {
			redirectParams[id] = vals
		}
	}

	// Add configured params that were not present in request.
	for id, configParam := range paramConfig {
		_, ok := redirectParams[id]
		if ok {
			continue
		}
		vals, err := configParam.Get(input)
		if err != nil {
			return nil, err
		}

		if len(vals) != 0 {
			redirectParams[id] = vals
		}
	}
	return redirectParams, nil
}
