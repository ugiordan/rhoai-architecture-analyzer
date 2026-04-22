package renderer

// NotableDependencyPrefixes lists Go module path prefixes considered "notable"
// for display in dependency diagrams and reports.
var NotableDependencyPrefixes = []string{
	"sigs.k8s.io/controller-runtime",
	"k8s.io/api",
	"k8s.io/apimachinery",
	"k8s.io/client-go",
	"github.com/operator-framework",
	"github.com/prometheus",
	"google.golang.org/grpc",
	"github.com/go-logr",
}

// KnownInternalPrefixes lists Go module path prefixes for known platform organizations.
var KnownInternalPrefixes = []string{
	"github.com/opendatahub-io/",
	"github.com/red-hat-data-services/",
}
