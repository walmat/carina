package http

type H2Headers string

const (
	H2HeadersAmps = "amps"
	H2HeadersMasp = "masp"
	H2HeadersMpas = "mpas"
)

var (
	DefaultFingerprint = &FingerprintConfig{
		H2Headers: H2HeadersAmps,
	}

	FingerprintChrome = &FingerprintConfig{
		H2Headers: H2HeadersMasp,
		H2Settings: []H2Setting{
			{ID: H2SettingHeaderTableSize, Val: 65536},
			{ID: H2SettingMaxConcurrentStreams, Val: 1000},
			{ID: H2SettingInitialWindowSize, Val: 6291456},
			{ID: H2SettingMaxHeaderListSize, Val: 262144},
		},
		WindowUpdate: 15728640,
		HeaderPriority: H2PriorityParam{
			StreamDep: 0,
			Exclusive: true,
			Weight:    255,
		},
	}

	FingerprintFirefox = &FingerprintConfig{
		H2Headers: H2HeadersMpas,
		H2Settings: []H2Setting{
			{ID: H2SettingHeaderTableSize, Val: 65536},
			{ID: H2SettingInitialWindowSize, Val: 131072},
			{ID: H2SettingMaxFrameSize, Val: 16384},
		},

		WindowUpdate: 12517377,

		InitialStreamID: 15,
		PriorityFrames: []H2PriorityFrame{
			{StreamID: 3, Param: H2PriorityParam{StreamDep: 0, Weight: 201, Exclusive: false}},
			{StreamID: 5, Param: H2PriorityParam{StreamDep: 0, Weight: 101, Exclusive: false}},
			{StreamID: 7, Param: H2PriorityParam{StreamDep: 0, Weight: 1, Exclusive: false}},
			{StreamID: 9, Param: H2PriorityParam{StreamDep: 7, Weight: 1, Exclusive: false}},
			{StreamID: 11, Param: H2PriorityParam{StreamDep: 3, Weight: 1, Exclusive: false}},
			{StreamID: 13, Param: H2PriorityParam{StreamDep: 0, Weight: 241, Exclusive: false}},
		},
	}
)

type H2PriorityFrame struct {
	StreamID uint32
	Param    H2PriorityParam
}

type FingerprintConfig struct {
	H2Headers
	H2Settings []H2Setting

	WindowUpdate uint32

	InitialStreamID uint32
	PriorityFrames  []H2PriorityFrame
	HeaderPriority  H2PriorityParam
}
