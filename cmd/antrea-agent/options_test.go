// Copyright 2022 Antrea Authors
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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	featuregatetesting "k8s.io/component-base/featuregate/testing"

	"antrea.io/antrea/pkg/agent/config"
	agentconfig "antrea.io/antrea/pkg/config/agent"
	"antrea.io/antrea/pkg/features"
)

func TestOptionsValidateTLSOptions(t *testing.T) {
	tests := []struct {
		name        string
		config      *agentconfig.AgentConfig
		expectedErr string
	}{
		{
			name: "empty input",
			config: &agentconfig.AgentConfig{
				TLSCipherSuites: "",
				TLSMinVersion:   "",
			},
			expectedErr: "",
		},
		{
			name: "invalid TLSMinVersion",
			config: &agentconfig.AgentConfig{
				TLSCipherSuites: "",
				TLSMinVersion:   "foo",
			},
			expectedErr: "invalid TLSMinVersion",
		},
		{
			name: "invalid TLSCipherSuites",
			config: &agentconfig.AgentConfig{
				TLSCipherSuites: "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, foo",
				TLSMinVersion:   "VersionTLS10",
			},
			expectedErr: "invalid TLSCipherSuites",
		},
		{
			name: "valid input",
			config: &agentconfig.AgentConfig{
				TLSCipherSuites: "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, TLS_RSA_WITH_AES_128_GCM_SHA256",
				TLSMinVersion:   "VersionTLS12",
			},
			expectedErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{config: tt.config}
			err := o.validateTLSOptions()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}
		})
	}
}

func TestOptionsValidateEgressConfig(t *testing.T) {
	tests := []struct {
		name                 string
		featureGateValue     bool
		trafficEncapMode     config.TrafficEncapModeType
		egressConfig         agentconfig.EgressConfig
		expectedErr          string
		expectedEnableEgress bool
	}{
		{
			name:                 "enabled",
			featureGateValue:     true,
			trafficEncapMode:     config.TrafficEncapModeEncap,
			expectedEnableEgress: true,
		},
		{
			name:                 "unsupported encap mode",
			featureGateValue:     true,
			trafficEncapMode:     config.TrafficEncapModeNoEncap,
			expectedEnableEgress: false,
		},
		{
			name:             "too large maxEgressIPsPerNode",
			featureGateValue: true,
			trafficEncapMode: config.TrafficEncapModeEncap,
			egressConfig: agentconfig.EgressConfig{
				MaxEgressIPsPerNode: 300,
			},
			expectedErr:          "maxEgressIPsPerNode cannot be greater than",
			expectedEnableEgress: false,
		},
		{
			name:             "invalid exceptCIDRs",
			featureGateValue: true,
			trafficEncapMode: config.TrafficEncapModeEncap,
			egressConfig: agentconfig.EgressConfig{
				ExceptCIDRs: []string{"1.1.1.300/32"},
			},
			expectedErr:          "Egress Except CIDR 1.1.1.300/32 is invalid",
			expectedEnableEgress: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer featuregatetesting.SetFeatureGateDuringTest(t, features.DefaultFeatureGate, features.Egress, tt.featureGateValue)()

			o := &Options{config: &agentconfig.AgentConfig{
				Egress: tt.egressConfig,
			}}
			err := o.validateEgressConfig(tt.trafficEncapMode)
			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.expectedErr)
			}
			assert.Equal(t, tt.expectedEnableEgress, o.enableEgress)
		})
	}
}
