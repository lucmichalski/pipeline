// Copyright © 2019 Banzai Cloud
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
	"context"

	"emperror.dev/errors"
	corev1 "k8s.io/api/core/v1"
)

const DeleteUserNamespacesActivityName = "delete-user-namespaces"

type DeleteUserNamespacesActivityInput struct {
	OrganizationID uint
	ClusterName    string
	K8sSecretID    string
}

type DeleteUserNamespacesActivityOutput struct {
	NamespacesLeft []string
}

type DeleteUserNamespacesActivity struct {
	deleter         UserNamespaceDeleter
	k8sConfigGetter K8sConfigGetter
}

type UserNamespaceDeleter interface {
	Delete(ctx context.Context, organizationID uint, clusterName string, nsFilter *corev1.NamespaceList, k8sConfig []byte) ([]string, error)
}

func MakeDeleteUserNamespacesActivity(deleter UserNamespaceDeleter, k8sConfigGetter K8sConfigGetter) DeleteUserNamespacesActivity {
	return DeleteUserNamespacesActivity{
		deleter:         deleter,
		k8sConfigGetter: k8sConfigGetter,
	}
}

func (a DeleteUserNamespacesActivity) Execute(ctx context.Context, input DeleteUserNamespacesActivityInput) (DeleteUserNamespacesActivityOutput, error) {
	k8sConfig, err := a.k8sConfigGetter.Get(input.OrganizationID, input.K8sSecretID)
	if err != nil {
		return DeleteUserNamespacesActivityOutput{}, errors.WrapIf(err, "failed to get k8s config")
	}
	left, err := a.deleter.Delete(ctx, input.OrganizationID, input.ClusterName, nil, k8sConfig)
	return DeleteUserNamespacesActivityOutput{left}, errors.WrapIf(err, "failed to delete user namespaces")
}
