// Copyright © 2020 Banzai Cloud
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

package workflow

import (
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/testsuite"

	"github.com/banzaicloud/pipeline/internal/cluster/distribution/eks"
)

func TestDeleteStoredNodePoolActivityExecute(t *testing.T) {
	testCases := []struct {
		caseName      string
		expectedError error
		input         DeleteStoredNodePoolActivityInput
		inputActivity *DeleteStoredNodePoolActivity
	}{
		{
			caseName:      "example error",
			expectedError: errors.New("test error: example"),
			input:         DeleteStoredNodePoolActivityInput{},
			inputActivity: &DeleteStoredNodePoolActivity{
				nodePoolStore: &eks.MockNodePoolStore{},
			},
		},
		{
			caseName:      "example success",
			expectedError: nil,
			input:         DeleteStoredNodePoolActivityInput{},
			inputActivity: &DeleteStoredNodePoolActivity{
				nodePoolStore: &eks.MockNodePoolStore{},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.caseName, func(t *testing.T) {
			environment := (&testsuite.WorkflowTestSuite{}).NewTestActivityEnvironment()

			environment.RegisterActivityWithOptions(testCase.inputActivity.Execute, activity.RegisterOptions{Name: t.Name()})

			testCase.inputActivity.nodePoolStore.(*eks.MockNodePoolStore).On(
				"DeleteNodePool",
				mock.Anything,
				testCase.input.OrganizationID,
				testCase.input.ClusterID,
				testCase.input.ClusterName,
				testCase.input.NodePoolName,
			).Return(testCase.expectedError).Once()

			_, actualError := environment.ExecuteActivity(t.Name(), testCase.input)

			if testCase.expectedError == nil {
				require.Nil(t, actualError)
			} else {
				require.EqualError(t, actualError, testCase.expectedError.Error())
			}
		})
	}
}

func TestNewDeleteStoredNodePoolActivity(t *testing.T) {
	testCases := []struct {
		caseName           string
		expectedActivity   *DeleteStoredNodePoolActivity
		inputNodePoolStore eks.NodePoolStore
	}{
		{
			caseName:           "nil node pool store",
			expectedActivity:   &DeleteStoredNodePoolActivity{},
			inputNodePoolStore: nil,
		},
		{
			caseName: "not nil node pool store",
			expectedActivity: &DeleteStoredNodePoolActivity{
				nodePoolStore: &eks.MockNodePoolStore{},
			},
			inputNodePoolStore: &eks.MockNodePoolStore{},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.caseName, func(t *testing.T) {
			actualActivity := NewDeleteStoredNodePoolActivity(testCase.inputNodePoolStore)

			require.Equal(t, testCase.expectedActivity, actualActivity)
		})
	}
}
