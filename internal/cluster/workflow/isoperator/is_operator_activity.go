// Copyright © 2021 Banzai Cloud
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

package isoperator

import (
	"context"

	"emperror.dev/errors"

	"github.com/banzaicloud/pipeline/internal/helm"
)

const (
	IntegratedServiceOperatorInstallerActivityName = "integrated-service-operator-installer-activity"
	GetNextClusterRefActivityName                  = "get-next-cluster-id-activity"
)

type IntegratedServicesOperatorInstallerActivityInput struct {
	OrgID     uint
	ClusterID uint
}

func NewInstallerActivityInput(orgID uint, clusterID uint) IntegratedServicesOperatorInstallerActivityInput {
	return IntegratedServicesOperatorInstallerActivityInput{
		OrgID:     orgID,
		ClusterID: clusterID,
	}
}

type IntegratedServicesOperatorInstallerActivity struct {
	config              Config
	clusterDataProvider helm.ClusterDataProvider
	repoUpdater         helm.Service
	chartReleaser       helm.UnifiedReleaser
}

func NewInstallerActivity(repoUpdater helm.Service, chartReleaser helm.UnifiedReleaser, config Config) IntegratedServicesOperatorInstallerActivity {
	return IntegratedServicesOperatorInstallerActivity{
		config:        config,
		repoUpdater:   repoUpdater,
		chartReleaser: chartReleaser,
	}
}

func (r IntegratedServicesOperatorInstallerActivity) Execute(ctx context.Context, input IntegratedServicesOperatorInstallerActivityInput) error {
	if err := r.repoUpdater.UpdateRepository(ctx,
		input.OrgID,
		helm.Repository{
			Name: r.config.RepoName,
			URL:  r.config.RepoURL,
		}); err != nil {
		return errors.WrapIf(err, "failed to update helm repository")
	}

	if err := r.chartReleaser.InstallOrUpgrade(
		input.OrgID,
		r.clusterDataProvider,
		helm.Release{
			ReleaseName: r.config.ReleaseName,
			ChartName:   r.config.ChartName,
			Namespace:   r.config.Namespace,
			Version:     r.config.ChartVersion,
		},
		helm.Options{
			Namespace: r.config.Namespace,
		},
	); err != nil {
		return errors.WrapIf(err, "failed to install or upgrade the chart")
	}
	return nil
}

type NextClusterIDActivity struct {
	NextIDProvider NextIDProvider
}

func NewNextClusterIDActivity(NextidProvider NextIDProvider) NextClusterIDActivity {
	return NextClusterIDActivity{
		NextIDProvider: NextidProvider,
	}
}

type ClusterRef struct {
	ID    uint
	OrgID uint
}

func (n NextClusterIDActivity) Execute(ctx context.Context, lastClusterID uint) (ClusterRef, error) {
	orgID, clusterID, err := n.NextIDProvider(lastClusterID)
	if err != nil {
		return ClusterRef{}, errors.WrapIfWithDetails(err, "failed to retrieve the next cluster reference")
	}

	return ClusterRef{ID: clusterID, OrgID: orgID}, nil
}