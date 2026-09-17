package gridphysics

// BuildIEEE14BusBenchmark returns the standard IEEE 14-Bus test feeder topology.
// It consists of 14 buses (1 Slack, 4 PV generators, 9 PQ loads) and 20 transmission lines.
func BuildIEEE14BusBenchmark() ([]Bus, []TransmissionLine) {
	buses := []Bus{
		{ID: 1, Name: "Bus 1 (Slack)", Type: BusSlack, PGenMW: 232.4, QGenMVAR: -16.9, PLoadMW: 0.0, QLoadMVAR: 0.0, VMagPU: 1.06, VAngRad: 0.0},
		{ID: 2, Name: "Bus 2 (Gen)", Type: BusPV, PGenMW: 40.0, QGenMVAR: 42.4, PLoadMW: 21.7, QLoadMVAR: 12.7, VMagPU: 1.045, VAngRad: 0.0},
		{ID: 3, Name: "Bus 3 (Sync)", Type: BusPV, PGenMW: 0.0, QGenMVAR: 23.4, PLoadMW: 94.2, QLoadMVAR: 19.0, VMagPU: 1.01, VAngRad: 0.0},
		{ID: 4, Name: "Bus 4 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 47.8, QLoadMVAR: -3.9, VMagPU: 1.0, VAngRad: 0.0},
		{ID: 5, Name: "Bus 5 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 7.6, QLoadMVAR: 1.6, VMagPU: 1.0, VAngRad: 0.0},
		{ID: 6, Name: "Bus 6 (Gen)", Type: BusPV, PGenMW: 0.0, QGenMVAR: 12.2, PLoadMW: 11.2, QLoadMVAR: 7.5, VMagPU: 1.07, VAngRad: 0.0},
		{ID: 7, Name: "Bus 7 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 0.0, QLoadMVAR: 0.0, VMagPU: 1.0, VAngRad: 0.0},
		{ID: 8, Name: "Bus 8 (Sync)", Type: BusPV, PGenMW: 0.0, QGenMVAR: 17.4, PLoadMW: 0.0, QLoadMVAR: 0.0, VMagPU: 1.09, VAngRad: 0.0},
		{ID: 9, Name: "Bus 9 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 29.5, QLoadMVAR: 16.6, VMagPU: 1.0, VAngRad: 0.0},
		{ID: 10, Name: "Bus 10 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 9.0, QLoadMVAR: 5.8, VMagPU: 1.0, VAngRad: 0.0},
		{ID: 11, Name: "Bus 11 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 3.5, QLoadMVAR: 1.8, VMagPU: 1.0, VAngRad: 0.0},
		{ID: 12, Name: "Bus 12 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 6.1, QLoadMVAR: 1.6, VMagPU: 1.0, VAngRad: 0.0},
		{ID: 13, Name: "Bus 13 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 13.5, QLoadMVAR: 5.8, VMagPU: 1.0, VAngRad: 0.0},
		{ID: 14, Name: "Bus 14 (Load)", Type: BusPQ, PGenMW: 0.0, QGenMVAR: 0.0, PLoadMW: 14.9, QLoadMVAR: 5.0, VMagPU: 1.0, VAngRad: 0.0},
	}

	lines := []TransmissionLine{
		{ID: "L1-2", FromBus: 1, ToBus: 2, Resistance: 0.01938, Reactance: 0.05917, MVARating: 200.0},
		{ID: "L1-5", FromBus: 1, ToBus: 5, Resistance: 0.05403, Reactance: 0.22304, MVARating: 100.0},
		{ID: "L2-3", FromBus: 2, ToBus: 3, Resistance: 0.04699, Reactance: 0.19797, MVARating: 100.0},
		{ID: "L2-4", FromBus: 2, ToBus: 4, Resistance: 0.05811, Reactance: 0.17632, MVARating: 100.0},
		{ID: "L2-5", FromBus: 2, ToBus: 5, Resistance: 0.05695, Reactance: 0.17388, MVARating: 100.0},
		{ID: "L3-4", FromBus: 3, ToBus: 4, Resistance: 0.06701, Reactance: 0.17103, MVARating: 100.0},
		{ID: "L4-5", FromBus: 4, ToBus: 5, Resistance: 0.01335, Reactance: 0.04211, MVARating: 100.0},
		{ID: "L4-7", FromBus: 4, ToBus: 7, Resistance: 0.0, Reactance: 0.20912, MVARating: 100.0},
		{ID: "L4-9", FromBus: 4, ToBus: 9, Resistance: 0.0, Reactance: 0.55618, MVARating: 100.0},
		{ID: "L5-6", FromBus: 5, ToBus: 6, Resistance: 0.0, Reactance: 0.25202, MVARating: 100.0},
		{ID: "L6-11", FromBus: 6, ToBus: 11, Resistance: 0.09498, Reactance: 0.19890, MVARating: 50.0},
		{ID: "L6-12", FromBus: 6, ToBus: 12, Resistance: 0.12291, Reactance: 0.25581, MVARating: 50.0},
		{ID: "L6-13", FromBus: 6, ToBus: 13, Resistance: 0.06615, Reactance: 0.13027, MVARating: 50.0},
		{ID: "L7-8", FromBus: 7, ToBus: 8, Resistance: 0.0, Reactance: 0.17615, MVARating: 100.0},
		{ID: "L7-9", FromBus: 7, ToBus: 9, Resistance: 0.0, Reactance: 0.11001, MVARating: 100.0},
		{ID: "L9-10", FromBus: 9, ToBus: 10, Resistance: 0.03181, Reactance: 0.08450, MVARating: 50.0},
		{ID: "L9-14", FromBus: 9, ToBus: 14, Resistance: 0.12711, Reactance: 0.27038, MVARating: 50.0},
		{ID: "L10-11", FromBus: 10, ToBus: 11, Resistance: 0.08205, Reactance: 0.19207, MVARating: 50.0},
		{ID: "L12-13", FromBus: 12, ToBus: 13, Resistance: 0.22092, Reactance: 0.19988, MVARating: 50.0},
		{ID: "L13-14", FromBus: 13, ToBus: 14, Resistance: 0.17093, Reactance: 0.34802, MVARating: 50.0},
	}

	return buses, lines
}
