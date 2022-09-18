package main

const (
	AMD    = "1002"
	NVIDIA = "10de"
)

type gpuInfo struct {
	Vendor   string
	DeviceID string
	Revision string
}

func GetGpusList() (gpuList map[gpuInfo][]string) {
	gpuList = map[gpuInfo][]string{
		// AMD GPUS
		gpuInfo{
			Vendor:   AMD,
			DeviceID: "67df",
			Revision: "0xcf",
		}: {"Radeon RX 470", "Ellesmere"},
		gpuInfo{
			Vendor:   AMD,
			DeviceID: "67df",
			Revision: "0xc7",
		}: {"Radeon RX 480", "Ellesmere"},
		gpuInfo{
			Vendor:   AMD,
			DeviceID: "67df",
			Revision: "0xef",
		}: {"Radeon RX 570", "Ellesmere"},
		gpuInfo{
			Vendor:   AMD,
			DeviceID: "67df",
			Revision: "0xe7",
		}: {"Radeon RX 580", "Ellesmere"},
		gpuInfo{
			Vendor:   AMD,
			DeviceID: "731f",
			Revision: "0xca",
		}: {"Radeon RX 5600 XT", "Navi 10"},
		// NVIDIA GPUS
	}

	return
}
