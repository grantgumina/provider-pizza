package v1alpha1

import (
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Borrowed from https://github.com/rudoi/cruster-api/blob/062c878029262197b24eff53a3e7913abf9298e6/api/v1/pizzaorder_types.go

type Pizza struct {
	// +kubebuilder:validation:Enum=small;medium;large
	Size string `json:"size"`

	Toppings []string `json:"toppings"`
}

type Customer struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`

	// +kubebuilder:validation:MaxLength=2
	Region string `json:"region"`

	PostalCode string `json:"postalCode"`

	// +kubebuilder:validation:Pattern="[2-9]\\d{9}$"
	Phone string `json:"phone"`
}

type StoreStatus struct {
	ID      string `json:"id,omitempty"`
	Address string `json:"address,omitempty"`
}

type Tracker struct {
	Prep           string `json:"prep,omitempty"`
	Bake           string `json:"bake,omitempty"`
	QualityCheck   string `json:"qualityCheck,omitempty"`
	OutForDelivery string `json:"outForDelivery,omitempty"`
	Delivered      string `json:"delivered,omitempty"`
}

// PizzaOrderSpec defines the desired state of PizzaOrder
type OrderSpec struct {
	PlaceOrder    bool                                `json:"placeOrder"`
	Address       *Address                            `json:"address"`
	Customer      *Customer                           `json:"customer"`
	PaymentSecret runtimev1alpha1.ProviderCredentials `json:"paymentSecret,omitempty"`

	// +kubebuilder:validation:MinItems=1
	Pizzas []*Pizza `json:"pizzas"`
}

// PizzaOrderStatus defines the observed state of PizzaOrder
type OrderStatus struct {
	OrderID   string       `json:"orderID,omitempty"`
	Price     string       `json:"price,omitempty"`
	Placed    bool         `json:"placed,omitempty"`
	Delivered bool         `json:"delivered,omitempty"`
	Store     *StoreStatus `json:"store,omitempty"`
	Tracker   *Tracker     `json:"tracker,omitempty"`
}

// +kubebuilder:object:root=true
type Order struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrderSpec   `json:"spec"`
	Status OrderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type OrderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Order `json:"items"`
}
