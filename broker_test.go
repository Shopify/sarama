package sarama

import (
	"fmt"
	"testing"
	"time"

	"github.com/rcrowley/go-metrics"
)

func ExampleBroker() {
	broker := NewBroker("localhost:9092")
	err := broker.Open(nil)
	if err != nil {
		panic(err)
	}

	request := MetadataRequest{Topics: []string{"myTopic"}}
	response, err := broker.GetMetadata(&request)
	if err != nil {
		_ = broker.Close()
		panic(err)
	}

	fmt.Println("There are", len(response.Topics), "topics active in the cluster.")

	if err = broker.Close(); err != nil {
		panic(err)
	}
}

type mockEncoder struct {
	bytes []byte
}

func (m mockEncoder) encode(pe packetEncoder) error {
	return pe.putRawBytes(m.bytes)
}

func TestBrokerAccessors(t *testing.T) {
	broker := NewBroker("abc:123")

	if broker.ID() != -1 {
		t.Error("New broker didn't have an ID of -1.")
	}

	if broker.Addr() != "abc:123" {
		t.Error("New broker didn't have the correct address")
	}

	broker.id = 34
	if broker.ID() != 34 {
		t.Error("Manually setting broker ID did not take effect.")
	}
}

func TestSimpleBrokerCommunication(t *testing.T) {
	for _, tt := range brokerTestTable {
		Logger.Printf("Testing broker communication for %s", tt.name)
		mb := NewMockBroker(t, 0)
		mb.Returns(&mockEncoder{tt.response})
		broker := NewBroker(mb.Addr())
		// Set the broker id in order to validate local broker metrics
		broker.id = 0
		conf := NewConfig()
		conf.Version = V0_10_0_0
		// Use a new registry every time to prevent side effect caused by the global one
		conf.MetricRegistry = metrics.NewRegistry()
		err := broker.Open(conf)
		if err != nil {
			t.Fatal(err)
		}
		tt.runner(t, broker)
		err = broker.Close()
		if err != nil {
			t.Error(err)
		}
		// Wait up to 500 ms for the remote broker to process requests
		// in order to have consistent metrics
		if err := mb.WaitForExpectations(500 * time.Millisecond); err != nil {
			t.Error(err)
		}
		mb.Close()
		validateBrokerMetrics(t, broker, mb)
	}

}

// We're not testing encoding/decoding here, so most of the requests/responses will be empty for simplicity's sake
var brokerTestTable = []struct {
	name     string
	response []byte
	runner   func(*testing.T, *Broker)
}{
	{"MetadataRequest",
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := MetadataRequest{}
			response, err := broker.GetMetadata(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("Metadata request got no response!")
			}
		}},

	{"ConsumerMetadataRequest",
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 't', 0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := ConsumerMetadataRequest{}
			response, err := broker.GetConsumerMetadata(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("Consumer Metadata request got no response!")
			}
		}},

	{"ProduceRequest (NoResponse)",
		[]byte{},
		func(t *testing.T, broker *Broker) {
			request := ProduceRequest{}
			request.RequiredAcks = NoResponse
			response, err := broker.Produce(&request)
			if err != nil {
				t.Error(err)
			}
			if response != nil {
				t.Error("Produce request with NoResponse got a response!")
			}
		}},

	{"ProduceRequest (WaitForLocal)",
		[]byte{0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := ProduceRequest{}
			request.RequiredAcks = WaitForLocal
			response, err := broker.Produce(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("Produce request without NoResponse got no response!")
			}
		}},

	{"FetchRequest",
		[]byte{0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := FetchRequest{}
			response, err := broker.Fetch(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("Fetch request got no response!")
			}
		}},

	{"OffsetFetchRequest",
		[]byte{0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := OffsetFetchRequest{}
			response, err := broker.FetchOffset(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("OffsetFetch request got no response!")
			}
		}},

	{"OffsetCommitRequest",
		[]byte{0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := OffsetCommitRequest{}
			response, err := broker.CommitOffset(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("OffsetCommit request got no response!")
			}
		}},

	{"OffsetRequest",
		[]byte{0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := OffsetRequest{}
			response, err := broker.GetAvailableOffsets(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("Offset request got no response!")
			}
		}},

	{"JoinGroupRequest",
		[]byte{0x00, 0x17, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := JoinGroupRequest{}
			response, err := broker.JoinGroup(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("JoinGroup request got no response!")
			}
		}},

	{"SyncGroupRequest",
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := SyncGroupRequest{}
			response, err := broker.SyncGroup(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("SyncGroup request got no response!")
			}
		}},

	{"LeaveGroupRequest",
		[]byte{0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := LeaveGroupRequest{}
			response, err := broker.LeaveGroup(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("LeaveGroup request got no response!")
			}
		}},

	{"HeartbeatRequest",
		[]byte{0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := HeartbeatRequest{}
			response, err := broker.Heartbeat(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("Heartbeat request got no response!")
			}
		}},

	{"ListGroupsRequest",
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := ListGroupsRequest{}
			response, err := broker.ListGroups(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("ListGroups request got no response!")
			}
		}},

	{"DescribeGroupsRequest",
		[]byte{0x00, 0x00, 0x00, 0x00},
		func(t *testing.T, broker *Broker) {
			request := DescribeGroupsRequest{}
			response, err := broker.DescribeGroups(&request)
			if err != nil {
				t.Error(err)
			}
			if response == nil {
				t.Error("DescribeGroups request got no response!")
			}
		}},
}

func validateBrokerMetrics(t *testing.T, broker *Broker, mockBroker *MockBroker) {
	metricValidators := newMetricValidators()
	mockBrokerBytesRead := 0
	mockBrokerBytesWritten := 0

	// Compute socket bytes
	for _, requestResponse := range mockBroker.History() {
		mockBrokerBytesRead += requestResponse.RequestSize
		mockBrokerBytesWritten += requestResponse.ResponseSize
	}

	// Check that the number of bytes sent corresponds to what the mock broker received
	metricValidators.registerForAllBrokers(broker, countMeterValidator("incoming-byte-rate", mockBrokerBytesWritten))
	if mockBrokerBytesWritten == 0 {
		// This a ProduceRequest with NoResponse
		metricValidators.registerForAllBrokers(broker, countMeterValidator("response-rate", 0))
		metricValidators.registerForAllBrokers(broker, countHistogramValidator("response-size", 0))
		metricValidators.registerForAllBrokers(broker, minMaxHistogramValidator("response-size", 0, 0))
	} else {
		metricValidators.registerForAllBrokers(broker, countMeterValidator("response-rate", 1))
		metricValidators.registerForAllBrokers(broker, countHistogramValidator("response-size", 1))
		metricValidators.registerForAllBrokers(broker, minMaxHistogramValidator("response-size", mockBrokerBytesWritten, mockBrokerBytesWritten))
	}

	// Check that the number of bytes received corresponds to what the mock broker sent
	metricValidators.registerForAllBrokers(broker, countMeterValidator("outgoing-byte-rate", mockBrokerBytesRead))
	metricValidators.registerForAllBrokers(broker, countMeterValidator("request-rate", 1))
	metricValidators.registerForAllBrokers(broker, countHistogramValidator("request-size", 1))
	metricValidators.registerForAllBrokers(broker, minMaxHistogramValidator("request-size", mockBrokerBytesRead, mockBrokerBytesRead))

	// Run the validators
	metricValidators.run(t, broker.conf.MetricRegistry)
}
