package v1alpha1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Package type metadata.
const (
	Group   = "order.provider-pizza.crossplane.io"
	Version = "v1alpha1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// Order type metadata
var (
	OrderKind             = reflect.TypeOf(Order{}).Name()
	OrderGroupKind        = schema.GroupKind{Group: Group, Kind: OrderKind}.String()
	OrderKindAPIVersion   = OrderKind + "." + SchemeGroupVersion.String()
	OrderGroupVersionKind = SchemeGroupVersion.WithKind(OrderKind)
)

func init() {
	SchemeBuilder.Register(&Order{}, &OrderList{})
}
