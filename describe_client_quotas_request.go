package sarama

// DescribeClientQuotas Request (Version: 0) => [components] strict
//   components => entity_type match_type match
//     entity_type => STRING
//     match_type => INT8
//     match => NULLABLE_STRING
//   strict => BOOLEAN

// A filter to be applied to matching client quotas.
// Components: the components to filter on
// Strict: whether the filter only includes specified components
type DescribeClientQuotasRequest struct {
	Components []*DescribeClientQuotasComponent
	Strict     bool
}

// A filter to be applied.
// EntityType: the entity type the filter component applies to ("user", "client-id")
// Match: if present, the name that's matched exactly
//        if empty, matches the default name
//        if null, matches any specified name
type DescribeClientQuotasComponent struct {
	EntityType string
	Match      *string
}

func (d *DescribeClientQuotasRequest) encode(pe packetEncoder) error {
	// Components
	if err := pe.putArrayLength(len(d.Components)); err != nil {
		return err
	}
	for _, c := range d.Components {
		if err := c.encode(pe); err != nil {
			return err
		}
	}

	// Strict
	pe.putBool(d.Strict)

	return nil
}

func (d *DescribeClientQuotasRequest) decode(pd packetDecoder, version int16) error {
	// Components
	componentCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if componentCount > 0 {
		d.Components = make([]*DescribeClientQuotasComponent, componentCount)
		for i := range d.Components {
			c := &DescribeClientQuotasComponent{}
			if err = c.decode(pd, version); err != nil {
				return err
			}
			d.Components[i] = c
		}
	} else {
		d.Components = []*DescribeClientQuotasComponent{}
	}

	// Strict
	strict, err := pd.getBool()
	if err != nil {
		return err
	}
	d.Strict = strict

	return nil
}

// How to match the entity {0 = exact name, 1 = default name, 2 = any specified name}.
func (d *DescribeClientQuotasComponent) getMatchType() int8 {
	if d.Match == nil {
		return 2
	}
	if len(*d.Match) == 0 {
		return 1
	}
	return 0
}

func (d *DescribeClientQuotasComponent) encode(pe packetEncoder) error {
	// EntityType
	if err := pe.putString(d.EntityType); err != nil {
		return err
	}

	// MatchType
	pe.putInt8(d.getMatchType())

	// Match
	if err := pe.putNullableString(d.Match); err != nil {
		return err
	}

	return nil
}

func (d *DescribeClientQuotasComponent) decode(pd packetDecoder, version int16) error {
	// EntityType
	entityType, err := pd.getString()
	if err != nil {
		return err
	}
	d.EntityType = entityType

	// MatchType (ignored)
	_, err = pd.getInt8()
	if err != nil {
		return err
	}

	// Match
	match, err := pd.getNullableString()
	if err != nil {
		return err
	}
	d.Match = match

	return nil
}

func (d *DescribeClientQuotasRequest) key() int16 {
	return 48
}

func (d *DescribeClientQuotasRequest) version() int16 {
	return 0
}

func (d *DescribeClientQuotasRequest) headerVersion() int16 {
	return 1
}

func (d *DescribeClientQuotasRequest) requiredVersion() KafkaVersion {
	return V2_6_0_0
}
