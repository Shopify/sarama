package sarama

import (
	"time"
)

// DescribeClientQuotas Response (Version: 0) => throttle_time_ms error_code error_message [entries]
//   throttle_time_ms => INT32
//   error_code => INT16
//   error_message => NULLABLE_STRING
//   entries => [entity] [values]
//     entity => entity_type entity_name
//       entity_type => STRING
//       entity_name => NULLABLE_STRING
//     values => key value
//       key => STRING
//       value => FLOAT64

type DescribeClientQuotasResponse struct {
	ThrottleTime time.Duration                // The duration in milliseconds for which the request was throttled due to a quota violation, or zero if the request did not violate any quota.
	ErrorCode    KError                       // The error code, or `0` if the quota description succeeded.
	ErrorMsg     string                       // The error message, or `null` if the quota description succeeded.
	Entries      []*DescribeClientQuotasEntry // A result entry.
}

type DescribeClientQuotasEntry struct {
	Entity map[string]*string // The quota entity description.
	Values map[string]float64 // The quota values for the entity.
}

func (d *DescribeClientQuotasResponse) encode(pe packetEncoder) error {
	// ThrottleTime
	pe.putInt32(int32(d.ThrottleTime / time.Millisecond))

	// ErrorCode
	pe.putInt16(int16(d.ErrorCode))

	// ErrorMsg
	var err error
	if d.ErrorMsg == "" {
		err = pe.putNullableString(nil)
	} else {
		err = pe.putString(d.ErrorMsg)
	}
	if err != nil {
		return err
	}

	// Entries
	if err = pe.putArrayLength(len(d.Entries)); err != nil {
		return err
	}
	for _, e := range d.Entries {
		if err = e.encode(pe); err != nil {
			return err
		}
	}

	return nil
}

func (d *DescribeClientQuotasResponse) decode(pd packetDecoder, version int16) error {
	// ThrottleTime
	throttleTime, err := pd.getInt32()
	if err != nil {
		return err
	}
	d.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	// ErrorCode
	errCode, err := pd.getInt16()
	if err != nil {
		return err
	}
	d.ErrorCode = KError(errCode)

	// ErrorMsg
	errMsg, err := pd.getString()
	if err != nil {
		return err
	}
	d.ErrorMsg = errMsg

	// Entries
	entryCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if entryCount > 0 {
		d.Entries = make([]*DescribeClientQuotasEntry, entryCount)
		for i := range d.Entries {
			e := &DescribeClientQuotasEntry{}
			if err = e.decode(pd, version); err != nil {
				return err
			}
			d.Entries[i] = e
		}
	} else {
		d.Entries = []*DescribeClientQuotasEntry{}
	}

	return nil
}

func (d *DescribeClientQuotasEntry) encode(pe packetEncoder) error {
	// Entity
	if err := pe.putArrayLength(len(d.Entity)); err != nil {
		return err
	}
	for entityType, entityName := range d.Entity {
		// entity_type
		if err := pe.putString(entityType); err != nil {
			return err
		}
		// entity_name
		if err := pe.putNullableString(entityName); err != nil {
			return err
		}
	}

	// Values
	if err := pe.putArrayLength(len(d.Values)); err != nil {
		return err
	}
	for key, value := range d.Values {
		// key
		if err := pe.putString(key); err != nil {
			return err
		}
		// value
		pe.putFloat64(value)
	}

	return nil
}

func (d *DescribeClientQuotasEntry) decode(pd packetDecoder, version int16) error {
	// Entity
	entityCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if entityCount > 0 {
		d.Entity = make(map[string]*string, entityCount)
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

			d.Entity[entityType] = entityName
		}
	} else {
		d.Entity = map[string]*string{}
	}

	// Values
	valueCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if valueCount > 0 {
		d.Values = make(map[string]float64, valueCount)
		for i := 0; i < valueCount; i++ {
			// key
			key, err := pd.getString()
			if err != nil {
				return err
			}
			// value
			value, err := pd.getFloat64()
			if err != nil {
				return err
			}
			d.Values[key] = value
		}
	} else {
		d.Values = map[string]float64{}
	}

	return nil
}

func (d *DescribeClientQuotasResponse) key() int16 {
	return 48
}

func (d *DescribeClientQuotasResponse) version() int16 {
	return 0
}

func (d *DescribeClientQuotasResponse) headerVersion() int16 {
	return 0
}

func (d *DescribeClientQuotasResponse) requiredVersion() KafkaVersion {
	return V2_6_0_0
}
