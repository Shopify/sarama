package sarama

import "testing"

var (
	alterClientQuotasResponseError = []byte{
		0, 0, 0, 0, // ThrottleTime
		0, 0, 0, 1, // Entries len
		0, 42, // ErrorCode (ErrInvalidRequest)
		0, 42, 'U', 'n', 'h', 'a', 'n', 'd', 'l', 'e', 'd', ' ', 'c', 'l', 'i', 'e', 'n', 't', ' ', 'q', 'u', 'o', 't', 'a', ' ', 'e', 'n', 't', 'i', 't', 'y', ' ', 't', 'y', 'p', 'e', ':', ' ', 'f', 'a', 'u', 'l', 't', 'y', // ErrorMsg
		0, 0, 0, 1, // Entity len
		0, 6, 'f', 'a', 'u', 'l', 't', 'y', // entityType
		255, 255, // entityName
	}

	alterClientQuotasResponseSingleEntry = []byte{
		0, 0, 0, 0, // ThrottleTime
		0, 0, 0, 1, // Entries len
		0, 0, // ErrorCode
		255, 255, // ErrorMsg
		0, 0, 0, 1, // Entity len
		0, 4, 'u', 's', 'e', 'r', // entityType
		255, 255, // entityName
	}

	alterClientQuotasResponseMultipleEntries = []byte{
		0, 0, 0, 0, // ThrottleTime
		0, 0, 0, 2, // Entries len
		0, 0, // ErrorCode
		255, 255, // ErrorMsg
		0, 0, 0, 1, // Entity len
		0, 4, 'u', 's', 'e', 'r', // entityType
		255, 255, // entityName
		0, 0, // ErrorCode
		255, 255, // ErrorMsg
		0, 0, 0, 1, // Entity len
		0, 9, 'c', 'l', 'i', 'e', 'n', 't', '-', 'i', 'd', // entityType
		255, 255, // entityName
	}
)

func TestAlterClientQuotasResponse(t *testing.T) {
	// Response with error
	entry := &AlterClientQuotasEntryResponse{
		ErrorCode: KError(42),
		ErrorMsg:  "Unhandled client quota entity type: faulty",
		Entity:    map[string]*string{"faulty": nil},
	}
	res := &AlterClientQuotasResponse{
		ThrottleTime: 0,
		Entries:      []*AlterClientQuotasEntryResponse{entry},
	}
	testResponse(t, "Response With Error", res, alterClientQuotasResponseError)

	// Response Altered single entry
	entry = &AlterClientQuotasEntryResponse{
		Entity: map[string]*string{"user": nil},
	}
	res = &AlterClientQuotasResponse{
		ThrottleTime: 0,
		Entries:      []*AlterClientQuotasEntryResponse{entry},
	}
	testResponse(t, "Altered single entry", res, alterClientQuotasResponseSingleEntry)

	// Response Altered multiple entries
	entry1 := &AlterClientQuotasEntryResponse{
		Entity: map[string]*string{"user": nil},
	}
	entry2 := &AlterClientQuotasEntryResponse{
		Entity: map[string]*string{"client-id": nil},
	}
	res = &AlterClientQuotasResponse{
		ThrottleTime: 0,
		Entries:      []*AlterClientQuotasEntryResponse{entry1, entry2},
	}
	testResponse(t, "Altered multiple entries", res, alterClientQuotasResponseMultipleEntries)
}
