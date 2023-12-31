/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Server struct {
	CreateCount      int    `json:"serverCreateCount,omitempty"`
	CreateStartNo    int    `json:"serverCreateStartNo,omitempty"`
	Description      string `json:"serverDescription,omitempty"`
	ImageNo          string `json:"serverImageNo,omitempty"`
	ImageProductCode string `json:"serverImageProductCode,omitempty"`
	Name             string `json:"serverName,omitempty"`
	ProductCode      string `json:"serverProductCode,omitempty"`
	SpecCode         string `json:"serverSpecCode,omitempty"`
}

type BlockStorageMapping struct {
	BlockStorageName           string `json:"blockStorageMappingBlockStorageName,omitempty"`
	BlockStorageSize           string `json:"blockStorageMappingBlockStorageSize,omitempty"`
	BlockStorageVolumeTypeCode string `json:"blockStorageMappingBlockStorageVolumeTypeCode,omitempty"`
	Encrypted                  string `json:"blockStorageMappingEncrypted,omitempty"`
	Order                      int    `json:"blockStorageMappingList,omitempty"`
	SnapshotInstanceNo         string `json:"blockStorageMappingSnapshotInstanceNo,omitempty"`
}

type NetworkInterface struct {
	IP       string `json:"networkInterfaceIp,omitempty"`
	No       string `json:"networkInterfaceNo,omitempty"`
	Order    int    `json:"networkInterfaceList,omitempty"`
	SubnetNo string `json:"networkInterfaceSubnetNo,omitempty"`
}

// ProvisionSpec defines the desired state of Provision
type ProvisionSpec struct {
	AccessControlGroupNoListN         string              `json:"accessControlGroupNoList,omitempty"`
	AssociateWithPublicIp             bool                `json:"associateWithPublicIp,omitempty"`
	BlockDevicePartitionMountPoint    string              `json:"blockDevicePartitionMountPoint,omitempty"`
	BlockDevicePartitionSize          string              `json:"blockDevicePartitionSize,omitempty"`
	FeeSystemTypeCode                 string              `json:"feeSystemTypeCode,omitempty"`
	InitScriptNo                      string              `json:"initScriptNo,omitempty"`
	IsEncryptedBaseBlockStorageVolume bool                `json:"isEncryptedBaseBlockStorageVolume,omitempty"`
	IsProtectServerTermination        bool                `json:"isProtectServerTermination,omitempty"`
	LoginKeyName                      string              `json:"loginKeyName,omitempty"`
	MemberServerImageInstanceNo       string              `json:"memberServerImageInstanceNo,omitempty"`
	PlacementGroupNo                  string              `json:"placementGroupNo,omitempty"`
	RAIDTypeName                      string              `json:"raidTypeName,omitempty"`
	ResponseFormatType                string              `json:"responseFormatType,omitempty"`
	SubnetNo                          string              `json:"subnetNo,omitempty"`
	VpcNo                             string              `json:"vpcNo,omitempty"`
	Server                            Server              `json:"server,omitempty"`
	BlockStorageMapping               BlockStorageMapping `json:"blockStorageMapping,omitempty"`
	NetworkInterface                  NetworkInterface    `json:"networkInterface,omitempty"`
}

type ProvisionPhase string

const (
	ProvisionPhaseCreate ProvisionPhase = "Create"
	ProvisionPhaseUpdate ProvisionPhase = "Update"
	ProvisionPhaseStop   ProvisionPhase = "Stop"
	ProvisionPhaseDelete ProvisionPhase = "Delete"
	ProvisionPhaseGet    ProvisionPhase = "Get"
	JobIsSuccess         ProvisionPhase = "Success"
)

// ProvisionStatus defines the observed state of Provision
type ProvisionStatus struct {
	Phase ProvisionPhase `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Provision is the Schema for the provisions API
type Provision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProvisionSpec   `json:"spec,omitempty"`
	Status ProvisionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProvisionList contains a list of Provision
type ProvisionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provision `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Provision{}, &ProvisionList{})
}
