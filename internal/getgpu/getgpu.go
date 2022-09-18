package getgpu

import (
	"encoding/json"
	"os/exec"
	"strings"

	pb "github.com/laxamore/mineralos/config/mineralos_proto"
)

type GPUDriverVersion struct {
	AMD    string
	NVIDIA string
}

type GPU struct {
	GpuVendor  string `json:"gpu_vendor"`
	GpuName    string `json:"gpu_name"`
	MemorySize string `json:"memory_size"`
}

func GetGPUDriverVersion() (drivers GPUDriverVersion, err error) {
	amdCommand := exec.Command("bash", "-c", "pacman -Qi opencl-amd | grep Version | cut -d':' -f 2 | awk '{sub(/ /,\"\"); print}'")
	nvidiaCommand := exec.Command("bash", "-c", "pacman -Qi nvidia | grep Version | cut -d':' -f 2 | awk '{sub(/ /,\"\"); print}'")

	amdCommandOutput, err := amdCommand.Output()
	if err != nil {
		return
	}
	nvidiaCommandOutput, err := nvidiaCommand.Output()
	if err != nil {
		return
	}

	drivers.AMD = string(amdCommandOutput)
	drivers.NVIDIA = string(nvidiaCommandOutput)

	drivers.AMD = strings.Replace(drivers.AMD, "\n", "", -1)
	drivers.NVIDIA = strings.Replace(drivers.NVIDIA, "\n", "", -1)

	return
}

func GetGPU() (gpus []GPU, err error) {
	gpusCommand := exec.Command("/mineralos/bin/gpus-info", "--json")

	gpusCommandOutput, err := gpusCommand.Output()
	if err != nil {
		return
	}

	json.Unmarshal(gpusCommandOutput, &gpus)
	return
}

func ArrGPUSToPBGPUS(gpus []GPU) (pbGPUS []*pb.GPUS) {
	for i := 0; i < len(gpus); i++ {
		pbGPUS = append(pbGPUS, &pb.GPUS{
			GpuVendor:  gpus[i].GpuVendor,
			GpuName:    gpus[i].GpuName,
			MemorySize: gpus[i].MemorySize,
		})
	}

	return
}
