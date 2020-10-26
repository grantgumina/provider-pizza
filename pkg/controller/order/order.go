package order

import (
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	ctrl "sigs.k8s.io/controller-runtime"
)

// "github.com/pkg/errors"
// v1 "k8s.io/api/core/v1"
// "k8s.io/apimachinery/pkg/types"
// ctrl "sigs.k8s.io/controller-runtime"
// "sigs.k8s.io/controller-runtime/pkg/client"

// "github.com/crossplane/crossplane-runtime/pkg/event"
// "github.com/crossplane/crossplane-runtime/pkg/logging"
// "github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
// "github.com/crossplane/crossplane-runtime/pkg/resource"

// "github.com/crossplane/provider-template/apis/sample/v1alpha1"
// apisv1alpha1 "github.com/crossplane/provider-template/apis/v1alpha1"
//

// Adds a controller which will reconcile Orders?
func Setup(mgr ctrl.Manager, l logging.Logger) error {
	// name := managed.ControllerName(v1alpha1.MyTypeGroupKind)
	name := managed.ControllerName()
}
