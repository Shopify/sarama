package sarama

import "testing"

var (
	describeClientQuotasResponseError = []byte{
		0, 0, 0, 0, // ThrottleTime
		0, 35, // ErrorCode
		0, 41, 'C', 'u', 's', 't', 'o', 'm', ' ', 'e', 'n', 't', 'i', 't', 'y', ' ', 't', 'y', 'p', 'e', ' ', '\'', 'f', 'a', 'u', 'l', 't', 'y', '\'', ' ', 'n', 'o', 't', ' ', 's', 'u', 'p', 'p', 'o', 'r', 't', 'e', 'd',
		0, 0, 0, 0, // Entries
	}

	describeClientQuotasResponseSingleValue = []byte{
		0, 0, 0, 0, // ThrottleTime
		0, 0, // ErrorCode
		255, 255, // ErrorMsg (nil)
		0, 0, 0, 1, // Entries
		0, 0, 0, 1, // Entity
		0, 4, 'u', 's', 'e', 'r', // Entity type
		255, 255, // Entity name (nil)
		0, 0, 0, 1, // Values
		0, 18, 'p', 'r', 'o', 'd', 'u', 'c', 'e', 'r', '_', 'b', 'y', 't', 'e', '_', 'r', 'a', 't', 'e',
		65, 46, 132, 128, 0, 0, 0, 0, // 1000000
	}

	describeClientQuotasResponseMultiValue = []byte{
		0, 0, 0, 0, // ThrottleTime
		0, 0, // ErrorCode
		255, 255, // ErrorMsg (nil)
		0, 0, 0, 1, // Entries
		0, 0, 0, 1, // Entity
		0, 4, 'u', 's', 'e', 'r', // Entity type
		255, 255, // Entity name (nil)
		0, 0, 0, 2, // Values
		0, 18, 'p', 'r', 'o', 'd', 'u', 'c', 'e', 'r', '_', 'b', 'y', 't', 'e', '_', 'r', 'a', 't', 'e',
		65, 46, 132, 128, 0, 0, 0, 0, // 1000000
		0, 18, 'c', 'o', 'n', 's', 'u', 'm', 'e', 'r', '_', 'b', 'y', 't', 'e', '_', 'r', 'a', 't', 'e',
		65, 46, 132, 128, 0, 0, 0, 0, // 1000000
	}

	describeClientQuotasResponseComplexEntity = []byte{
		0, 0, 0, 0, // ThrottleTime
		0, 0, // ErrorCode
		255, 255, // ErrorMsg (nil)
		0, 0, 0, 2, // Entries
		0, 0, 0, 1, // Entity
		0, 4, 'u', 's', 'e', 'r', // Entity type
		255, 255, // Entity name (nil)
		0, 0, 0, 1, // Values
		0, 18, 'p', 'r', 'o', 'd', 'u', 'c', 'e', 'r', '_', 'b', 'y', 't', 'e', '_', 'r', 'a', 't', 'e',
		65, 46, 132, 128, 0, 0, 0, 0, // 1000000
		0, 0, 0, 1, // Entity
		0, 9, 'c', 'l', 'i', 'e', 'n', 't', '-', 'i', 'd', // Entity type
		0, 6, 's', 'a', 'r', 'a', 'm', 'a', // Entity name
		0, 0, 0, 1, // Values
		0, 18, 'c', 'o', 'n', 's', 'u', 'm', 'e', 'r', '_', 'b', 'y', 't', 'e', '_', 'r', 'a', 't', 'e',
		65, 46, 132, 128, 0, 0, 0, 0, // 1000000
	}
)

func TestDescribeClientQuotasResponse(t *testing.T) {
	// Response With Error
	res := &DescribeClientQuotasResponse{
		ThrottleTime: 0,
		ErrorCode:    ErrUnsupportedVersion,
		ErrorMsg:     "Custom entity type 'faulty' not supported",
		Entries:      []*DescribeClientQuotasEntry{},
	}
	testResponse(t, "Response With Error", res, describeClientQuotasResponseError)

	// Single Quota entry
	entry := &DescribeClientQuotasEntry{
		Entity: map[string]*string{"user": nil},
		Values: map[string]float64{"producer_byte_rate": 1000000},
	}
	res = &DescribeClientQuotasResponse{
		ThrottleTime: 0,
		ErrorCode:    ErrNoError,
		ErrorMsg:     "",
		Entries:      []*DescribeClientQuotasEntry{entry},
	}
	testResponse(t, "Single Value", res, describeClientQuotasResponseSingleValue)

	// Multi Quota entry
	entry = &DescribeClientQuotasEntry{
		Entity: map[string]*string{"user": nil},
		Values: map[string]float64{
			"producer_byte_rate": 1000000,
			"consumer_byte_rate": 1000000,
		},
	}
	res = &DescribeClientQuotasResponse{
		ThrottleTime: 0,
		ErrorCode:    ErrNoError,
		ErrorMsg:     "",
		Entries:      []*DescribeClientQuotasEntry{entry},
	}
	testResponse(t, "Multi Value", res, describeClientQuotasResponseMultiValue)

	// Complex Quota entry
	clientId := "sarama"
	userEntry := &DescribeClientQuotasEntry{
		Entity: map[string]*string{"user": nil},
		Values: map[string]float64{"producer_byte_rate": 1000000},
	}
	clientEntry := &DescribeClientQuotasEntry{
		Entity: map[string]*string{"client-id": &clientId},
		Values: map[string]float64{"consumer_byte_rate": 1000000},
	}
	res = &DescribeClientQuotasResponse{
		ThrottleTime: 0,
		ErrorCode:    ErrNoError,
		ErrorMsg:     "",
		Entries:      []*DescribeClientQuotasEntry{userEntry, clientEntry},
	}
	testResponse(t, "Complex Quota", res, describeClientQuotasResponseComplexEntity)
}
