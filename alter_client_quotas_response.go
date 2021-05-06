package sarama

import (
	"time"
)

// AlterClientQuotas Response (Version: 0) => throttle_time_ms [entries]
//   throttle_time_ms => INT32
//   entries => error_code error_message [entity]
//     error_code => INT16
//     error_message => NULLABLE_STRING
//     entity => entity_type entity_name
//       entity_type => STRING
//       entity_name => NULLABLE_STRING

type AlterClientQuotasResponse struct {
	ThrottleTime time.Duration                     // The duration in milliseconds for which the request was throttled due to a quota violation, or zero if the request did not violate any quota.
	Entries      []*AlterClientQuotasEntryResponse // The quota configuration entries altered.
}

type AlterClientQuotasEntryResponse struct {
	ErrorCode KError             // The error code, or `0` if the quota alteration succeeded.
	ErrorMsg  string             // The error message, or `null` if the quota alteration succeeded.
	Entity    map[string]*string // The quota entity altered.
}

func (a *AlterClientQuotasResponse) encode(pe packetEncoder) error {
	// ThrottleTime
	pe.putInt32(int32(a.ThrottleTime / time.Millisecond))

	// Entries
	if err := pe.putArrayLength(len(a.Entries)); err != nil {
		return err
	}
	for _, e := range a.Entries {
		if err := e.encode(pe); err != nil {
			return err
		}
	}

	return nil
}

func (a *AlterClientQuotasResponse) decode(pd packetDecoder, version int16) error {
	// ThrottleTime
	throttleTime, err := pd.getInt32()
	if err != nil {
		return err
	}
	a.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	// Entries
	entryCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if entryCount > 0 {
		a.Entries = make([]*AlterClientQuotasEntryResponse, entryCount)
		for i := range a.Entries {
			e := &AlterClientQuotasEntryResponse{}
			if err = e.decode(pd, version); err != nil {
				return err
			}
			a.Entries[i] = e
		}
	} else {
		a.Entries = []*AlterClientQuotasEntryResponse{}
	}

	return nil
}

func (a *AlterClientQuotasEntryResponse) encode(pe packetEncoder) error {
	// ErrorCode
	pe.putInt16(int16(a.ErrorCode))

	// ErrorMsg
	var err error
	if a.ErrorMsg == "" {
		err = pe.putNullableString(nil)
	} else {
		err = pe.putString(a.ErrorMsg)
	}
	if err != nil {
		return err
	}

	// Entity
	if err = pe.putArrayLength(len(a.Entity)); err != nil {
		return err
	}

	for entityType, entityName := range a.Entity {
		// entity_type
		if err = pe.putString(entityType); err != nil {
			return err
		}
		// entity_name
		if err = pe.putNullableString(entityName); err != nil {
			return err
		}
	}

	return nil
}

func (a *AlterClientQuotasEntryResponse) decode(pd packetDecoder, version int16) error {
	// ErrorCode
	errCode, err := pd.getInt16()
	if err != nil {
		return err
	}
	a.ErrorCode = KError(errCode)

	// ErrorMsg
	errMsg, err := pd.getString()
	if err != nil {
		return err
	}
	a.ErrorMsg = errMsg

	// Entity
	entityCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if entityCount > 0 {
		a.Entity = make(map[string]*string, entityCount)
		for i := 0; i < entityCount; i++ {
			// entity_type
			entityType, err := pd.getString()
			if err != nil {
				return err
			}
			// entity_name
			entityName, err := pd.getNullableString()
			if err != nil {
				return err
			}
			a.Entity[entityType] = entityName
		}
	} else {
		a.Entity = map[string]*string{}
	}

	return nil
}

func (a *AlterClientQuotasResponse) key() int16 {
	return 49
}

func (a *AlterClientQuotasResponse) version() int16 {
	return 0
}

func (a *AlterClientQuotasResponse) headerVersion() int16 {
	return 0
}

func (a *AlterClientQuotasResponse) requiredVersion() KafkaVersion {
	return V2_6_0_0
}
