package Linux

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type GPUDriverVersion struct {
	AMD    string
	NVIDIA string
}

type GPU struct {
	GPU_Vendor  string `json:"gpu_vendor"`
	GPU_Name    string `json:"gpu_name"`
	Memory_Size string `json:"memory_size"`
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
	gpusCommand := exec.Command("bash", "-c", "gpus-info --json")

	gpusCommandOutput, err := gpusCommand.Output()
	if err != nil {
		return
	}

	json.Unmarshal(gpusCommandOutput, &gpus)
	return
}
