package eds

import (
	"net"
	"testing"

	tassert "github.com/stretchr/testify/assert"

	"github.com/openservicemesh/osm/pkg/endpoint"
	"github.com/openservicemesh/osm/pkg/service"
)

func TestNewClusterLoadAssignment(t *testing.T) {
	assert := tassert.New(t)

	namespacedServices := []service.MeshService{
		{Namespace: "ns1", Name: "bookstore-1", TargetPort: 80},
		{Namespace: "ns2", Name: "bookstore-2", TargetPort: 90},
	}

	allServiceEndpoints := map[service.MeshService][]endpoint.Endpoint{
		namespacedServices[0]: {
			{IP: net.IP("0.0.0.0")},
		},
		namespacedServices[1]: {
			{IP: net.IP("0.0.0.1")},
			{IP: net.IP("0.0.0.2")},
		},
	}

	cla := newClusterLoadAssignment(namespacedServices[0], allServiceEndpoints[namespacedServices[0]])
	assert.NotNil(cla)
	assert.Equal(cla.ClusterName, "ns1/bookstore-1|80")
	assert.Len(cla.Endpoints, 1)
	assert.Len(cla.Endpoints[0].LbEndpoints, 1)
	assert.Equal(cla.Endpoints[0].LbEndpoints[0].GetLoadBalancingWeight().Value, uint32(100))

	cla2 := newClusterLoadAssignment(namespacedServices[1], allServiceEndpoints[namespacedServices[1]])
	assert.NotNil(cla2)
	assert.Equal(cla2.ClusterName, "ns2/bookstore-2|90")
	assert.Len(cla2.Endpoints, 1)
	assert.Len(cla2.Endpoints[0].LbEndpoints, 2)
	assert.Equal(cla2.Endpoints[0].LbEndpoints[0].GetLoadBalancingWeight().Value, uint32(50))
	assert.Equal(cla2.Endpoints[0].LbEndpoints[1].GetLoadBalancingWeight().Value, uint32(50))

	cla3 := newClusterLoadAssignment(namespacedServices[0], []endpoint.Endpoint{})
	assert.NotNil(cla3)
	assert.Equal(cla3.ClusterName, "ns1/bookstore-1|80")
	assert.Len(cla3.Endpoints, 1)
	assert.Len(cla3.Endpoints[0].LbEndpoints, 0)
}
