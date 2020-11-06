package order

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/grantgumina/provider-pizza/apis/order/v1alpha1"

	apisv1alpha1 "github.com/grantgumina/provider-pizza/apis/v1alpha1"
	"github.com/pkg/errors"
	pz "github.com/rudoi/pizza-go/pkg/pizza"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Setup adds a controller which will reconcile Orders?
func Setup(mgr ctrl.Manager, l logging.Logger) error {
	name := managed.ControllerName(v1alpha1.OrderKind)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.OrderGroupVersionKind),
		managed.WithExternalConnecter(
			&connector{
				kube:  mgr.GetClient(),
				usage: resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
				// newServiceFn: pizzaClient,
			},
		),
		managed.WithLogger(l.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.Order{}).
		Complete(r)

}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube  client.Client
	usage resource.Tracker
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the ProviderConfig's credentials secret.
// 4. Using the credentials secret to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Order)
	if !ok {
		return nil, errors.New("errNotMyType")
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, "errTrackPCUsage")
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, "errGetPC")
	}

	// A secret is the most common way to authenticate to a provider, but some
	// providers additionally support alternative authentication methods such as
	// IAM, so a reference is not required.
	ref := pc.Spec.Credentials.SecretRef
	if ref == nil {
		return nil, errors.New("errNoSecretRef")
	}

	s := &v1.Secret{}
	if err := c.kube.Get(ctx, types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}, s); err != nil {
		return nil, errors.Wrap(err, "errGetSecret")
	}

	client := &pz.Client{Client: http.Client{Timeout: 10 * time.Second}}

	return &external{service: nil, pizzaClient: client}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	service     interface{}
	pizzaClient *pz.Client
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {

	cr, ok := mg.(*v1alpha1.Order)
	if !ok {
		return managed.ExternalObservation{}, errors.New("Not the right type")
	}

	// These fmt statements should be removed in the real implementation.
	fmt.Printf("Observing: %+v", cr)

	// Look at order yaml to see if we should check to see if this is a test or not

	realOrder := cr.Spec.ForProvider.PaymentMethod

	if realOrder != "Test" {
		trackingURL, err := c.pizzaClient.GetTrackingUrl(cr.Spec.ForProvider.Address.Phone)

		if err != nil {
			return managed.ExternalObservation{}, errors.New("Could not get tracking URL on this order")
		}

		trackerStatus, err := c.pizzaClient.Track(trackingURL)

		if err != nil {
			return managed.ExternalObservation{}, errors.New("Could not get tracking information on this order")
		}

		cr.Status.AtProvider.ManagerName = trackerStatus.ManagerName
		cr.Status.AtProvider.OrderStage = trackerStatus.OrderStatus
		cr.Status.AtProvider.Store.Phone = trackerStatus.Phone
	} else {
		cr.Status.AtProvider.ManagerName = "John Doe"
		cr.Status.AtProvider.OrderStage = "Routing Station"
		cr.Status.AtProvider.Store.Phone = "2068675309"
	}

	return managed.ExternalObservation{
		// Return false when the external resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: false,

		// Return false when the external resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: false,

		// Return any details that may be required to connect to the external
		// resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil

}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Order)
	if !ok {
		return managed.ExternalCreation{}, errors.New("errNotMyType")
	}

	fmt.Printf("Creating: %+v", cr)

	// Create the object from the user's YAML
	customer := cr.Spec.ForProvider.Customer
	address := &pz.Address{
		Street:     cr.Spec.ForProvider.Address.Street,
		City:       cr.Spec.ForProvider.Address.City,
		Region:     cr.Spec.ForProvider.Address.Region,
		PostalCode: cr.Spec.ForProvider.Address.PostalCode,
	}

	// The the store closest to this address
	store, err := c.pizzaClient.GetNearestStore(address)

	cr.Status.AtProvider.Store.Address = strings.Replace(store.Address, "\n", " ", -1)

	if err != nil {
		errors.New(err.Error())
		return managed.ExternalCreation{}, errors.New("Could not find a store near your address")
	}

	// Create an order object with the relevant information
	order := pz.NewOrder().
		WithAddress(address).
		WithCustomerInfo(customer.FirstName, customer.LastName, customer.Email).
		WithPhoneNumber(cr.Spec.ForProvider.Address.Phone).
		WithStoreID(store.StoreID)

	// Generate the product codes...
	// Search through the menu and find the codes for what's being requested

	for _, product := range cr.Spec.ForProvider.Products {
		p := &pz.OrderProduct{}

		switch product.Name {
		case "Bottle of Coke":
			p.Code = "20BDCOKE"
		case "2L Coke":
			p.Code = "2LCOKE"
		case "Cheese Pizza":
			p.Code = "14SCREEN"
		}

		p.Qty = 1
		order.AddProduct(p)
	}

	fmt.Println(order.Products)

	order.ServiceMethod = cr.Spec.ForProvider.ServiceMethod

	// For some reason the order isn't valid!!!!!

	// Get order price
	price, err := c.pizzaClient.ValidateOrder(order)

	if err != nil {
		errors.New(err.Error())
		// set some property on the order telling the observe loop not to get the tracking URL
		return managed.ExternalCreation{}, errors.New("Order not valid")
	}

	cr.Status.SetConditions(runtimev1alpha1.Creating())
	cr.Status.AtProvider.Price = strconv.FormatFloat(price, 'f', 2, 64)

	// Place the order
	paymentMethod := cr.Spec.ForProvider.PaymentMethod

	fmt.Println("<><> METHOD <><>")
	fmt.Println(paymentMethod)

	switch paymentMethod {
	case "Test":
	case "Cash":
		payment := &pz.Payment{
			Type: "Cash",
		}

		order.Payments = append(order.Payments, payment)

		fmt.Println("CASH IS KING")
	case "CreditCard":
		// pc := &apisv1alpha1.ProviderConfig{}

		// if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		// 	return nil, errors.Wrap(err, errGetPC)
		// }

		// ref := pc.Spec.Credentials.SecretRef

		// if ref == nil {
		// 	return managed.ExternalCreation{}, errors.New("No secret ref found")
		// }

		// s := &v1.Secret{}

		// fmt.Println(s.Data[ref.Key])

		// Read the secrets

		// paymentSecret := cr.Spec.ForProvider.PaymentSecret
		// if paymentSecret == nil {
		// 	return managed.ExternalCreation{}, errors.New(errNoSecretRef)
		// }

	}

	// Look at order yaml to see if we should place an order or not...
	returnedOrder, err := c.pizzaClient.PlaceOrder(order)

	if err != nil {
		return managed.ExternalCreation{}, errors.New("Order couldn't be placed")
	}

	fmt.Println(returnedOrder)

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, err
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Order)
	if !ok {
		return managed.ExternalUpdate{}, errors.New("Could not update order")
	}

	fmt.Printf("Updating: %+v", cr)

	cr.Status.SetConditions(runtimev1alpha1.Available())

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Order)
	if !ok {
		return errors.New("Could not delete order")
	}

	fmt.Printf("Deleting: %+v", cr)
	cr.Status.SetConditions(runtimev1alpha1.Deleting())

	return nil
}
