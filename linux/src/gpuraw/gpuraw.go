package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jaypipes/ghw"
)

type gpu struct {
	Vendor  string
	Name    string
	GhwInfo ghw.GraphicsCard
}

func main() {
	gpuList := GetGpusList()

	ghwGPU, err := ghw.GPU(ghw.WithDisableWarnings())
	if err != nil {
		fmt.Printf("Error getting GPU info: %v", err)
	}

	var gpus []ghw.GraphicsCard

	for _, card := range ghwGPU.GraphicsCards {
		if card.DeviceInfo.Vendor.ID == AMD || card.DeviceInfo.Vendor.ID == NVIDIA {
			gpuDevice := gpuInfo{
				Vendor:   AMD,
				DeviceID: card.DeviceInfo.Product.ID,
				Revision: card.DeviceInfo.Revision,
			}

			if card.DeviceInfo.Vendor.ID == AMD {
				card.DeviceInfo.Vendor.Name = "AMD"
			} else {
				card.DeviceInfo.Vendor.Name = "NVIDIA"
			}

			if gpuList[gpuDevice] != nil {
				card.DeviceInfo.Product.Name = gpuList[gpuDevice][0]
				gpus = append(gpus, *card)
			} else {
				gpus = append(gpus, *card)
			}
		} else {
			gpus = append(gpus, *card)
		}
	}

	if len(gpus) > 0 {
		sort.Slice(gpus, func(i, j int) bool {
			return sort.StringsAreSorted([]string{
				gpus[i].Address,
				gpus[j].Address,
			})
		})

		gpusJson, err := json.MarshalIndent(gpus, "", "\t")
		if err != nil {
			fmt.Printf("error getting gpu json: %v", err)
		} else {
			fmt.Println(string(gpusJson))
		}
	}
}
