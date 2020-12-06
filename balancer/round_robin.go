package balancer

import "google.golang.org/grpc/balancer/roundrobin"

// Just use grpc Round Robin balancer.
// No need making such wheel ourselves.
const RoundRobin = roundrobin.Name
