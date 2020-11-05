module github.com/grantgumina/provider-pizza

go 1.13

require (
	github.com/crossplane/crossplane-runtime v0.10.0
	github.com/crossplane/crossplane-tools v0.0.0-20201007233256-88b291e145bb
	github.com/crossplane/provider-template v0.0.0-20201019145102-837a2ee9aeeb
	github.com/google/go-cmp v0.4.0
	github.com/pkg/errors v0.9.1
	github.com/rudoi/pizza-go v0.0.0-20190722033559-b192c8d29127
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/controller-tools v0.2.4
)
