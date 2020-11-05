/*
Copyright 2020 The Crossplane Authors.

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

package v1alpha1

import (
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Borrowed from https://github.com/rudoi/cruster-api/blob/062c878029262197b24eff53a3e7913abf9298e6/api/v1/pizzaorder_types.go

// Product is a configurable field of Order
type Product struct {
	Name string `json:"name"`
}

// Customer is a configurable field of Order
type Customer struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// Address is a configurable field of Order
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`

	// +kubebuilder:validation:MaxLength=2
	Region string `json:"region"`

	PostalCode string `json:"postalCode"`

	// +kubebuilder:validation:Pattern="[2-9]\\d{9}$"
	Phone string `json:"phone"`
}

// StoreStatus is a configurable field of Order
type StoreStatus struct {
	ID      string `json:"id,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

// OrderObservation are the observable fields of a Order.
type OrderObservation struct {
	ObservableField string      `json:"observableField,omitempty"`
	Store           StoreStatus `json:"store,omitempty"`
	ManagerName     string      `json:"managerName,omitempty"`
	Price           string      `json:"price,omitempty"`
	OrderStage      string      `json:"orderStage,omitempty"`
}

// OrderParameters are the configurable fields of a Order.
type OrderParameters struct {
	PlaceOrder    bool      `json:"placeOrder"`
	Address       Address   `json:"address"`
	Customer      Customer  `json:"customer"`
	OrderID       string    `json:"orderID,omitempty"`
	Placed        bool      `json:"placed,omitempty"`
	ServiceMethod string    `json:"serviceMethod,omitempty"`
	PaymentMethod string    `json:"paymentMethod,omitempty"`
	Products      []Product `json:"products"`
	// PaymentSecret runtimev1alpha1.ProviderCredentials `json:"paymentSecret,omitempty"`
}

// OrderSpec defines the desired state of an Order
type OrderSpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`
	ForProvider                  OrderParameters `json:"forProvider"`
}

// OrderStatus defines the observed state of an Order
type OrderStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
	AtProvider                     OrderObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// Order manipulates the Domino's API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ORDER STATUS",type="string",JSONPath=".status.atProvider.orderStage"
// +kubebuilder:printcolumn:name="PRICE",type="string",JSONPath=".status.atProvider.price"
// +kubebuilder:printcolumn:name="STORE ADDRESS",type="string",JSONPath=".status.atProvider.store.address"
// +kubebuilder:printcolumn:name="STORE PHONE",type="string",JSONPath=".status.atProvider.store.phone"
// +kubebuilder:printcolumn:name="MANAGER",type="string",JSONPath=".status.atProvider.managerName"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster
type Order struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrderSpec   `json:"spec"`
	Status OrderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OrderList contains a list of Order
type OrderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Order `json:"items"`
}
