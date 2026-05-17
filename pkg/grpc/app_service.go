package grpc

// TODO: Generate proto stubs by running `make generate` or `buf generate` after
// adding apps.proto to the buf configuration. Once the generated code is
// available under pkg/protos/apps, implement the Apps gRPC service interface here.

// AppService handles gRPC requests for the Apps service.
// It will implement the pb_apps.AppsServer interface once proto generation is complete.
type AppService struct{}

// NewAppService creates a new AppService instance.
func NewAppService() *AppService {
	return &AppService{}
}
