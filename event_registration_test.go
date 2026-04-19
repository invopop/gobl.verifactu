package verifactu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventFingerprint(t *testing.T) {
	tests := []struct {
		name     string
		event    *Event
		prev     *EventChainData
		expected string
	}{
		{
			name: "First event with NIF",
			event: &Event{
				Software: &EventSoftware{
					NIF:                "89890001K",
					SoftwareID:         "77",
					Version:            "1.0.03",
					InstallationNumber: "383",
				},
				Issuer: &EventIssuer{
					NIF: "A28083806",
				},
				EventType:           "01",
				GenerationTimestamp: "2024-11-20T19:00:55+01:00",
			},
			prev:     nil,
			expected: "96E4F799E89D44E5DF9339384C275CEB44F271F8FC9D88F60AEFFD809EFCA84A",
		},
		{
			name: "With previous event",
			event: &Event{
				Software: &EventSoftware{
					NIF:                "89890001K",
					SoftwareID:         "77",
					Version:            "1.0.03",
					InstallationNumber: "383",
				},
				Issuer: &EventIssuer{
					NIF: "A28083806",
				},
				EventType:           "01",
				GenerationTimestamp: "2024-11-21T10:00:55+01:00",
			},
			prev: &EventChainData{
				EventType:           "01",
				GenerationTimestamp: "2024-11-20T19:00:55+01:00",
				Fingerprint:         "96E4F799E89D44E5DF9339384C275CEB44F271F8FC9D88F60AEFFD809EFCA84A",
			},
			expected: "9E628216C2F194123D14D18FD8D014C4D5BB699FC9055E60B3D89D99A031336B",
		},
		{
			name: "First event with IDOtro instead of NIF",
			event: &Event{
				Software: &EventSoftware{
					IDOther: &EventOtherID{
						ID: "ESJ1234567",
					},
					SoftwareID:         "77",
					Version:            "1.0.03",
					InstallationNumber: "383",
				},
				Issuer: &EventIssuer{
					NIF: "B08194359",
				},
				EventType:           "02",
				GenerationTimestamp: "2024-11-20T20:00:55+01:00",
			},
			prev:     nil,
			expected: "C9D12D1E39829D7DE69384604FFED76DBCB280DE32A49BF04F0E670C1995936C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.event.fingerprint(tt.prev)
			assert.Equal(t, tt.expected, tt.event.Fingerprint)

			if tt.prev == nil {
				assert.Equal(t, "S", tt.event.Chaining.FirstEvent)
				assert.Nil(t, tt.event.Chaining.PreviousEvent)
			} else {
				assert.Empty(t, tt.event.Chaining.FirstEvent)
				assert.NotNil(t, tt.event.Chaining.PreviousEvent)
				assert.Equal(t, tt.prev.EventType, tt.event.Chaining.PreviousEvent.EventType)
				assert.Equal(t, tt.prev.GenerationTimestamp, tt.event.Chaining.PreviousEvent.GenerationTimestamp)
				assert.Equal(t, tt.prev.Fingerprint, tt.event.Chaining.PreviousEvent.Fingerprint)
			}

			cd := tt.event.ChainData()
			assert.Equal(t, tt.event.EventType, cd.EventType)
			assert.Equal(t, tt.event.GenerationTimestamp, cd.GenerationTimestamp)
			assert.Equal(t, tt.expected, cd.Fingerprint)
		})
	}
}
