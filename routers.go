package vortex

import (
	"net/http"
	"runtime"

	vortexMw "github.com/dzjyyds666/VortexCore/middleware"
	vortexUtil "github.com/dzjyyds666/VortexCore/utils"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func prepareDefaultHttpRouter() []*HttpRouter {
	return []*HttpRouter{
		AppendHttpRouter([]string{http.MethodGet}, "/system/info", GetSystemInfo, "获取系统信息", vortexMw.JwtSkipMw()),
		AppendHttpRouter([]string{http.MethodGet}, "/checkAlive", checkAlive, "检查服务是否正常", vortexMw.JwtSkipMw()),
	}
}

type systemInfo struct {
	TotalMemory uint64 `json:"totalMemory"` // 总内存
	UsedMemory  uint64 `json:"usedMemory"`  // 已使用的内存

	TotalDisk uint64 `json:"totalDisk"` // 磁盘总空间
	UsedDisk  uint64 `json:"usedDisk"`  // 已使用的磁盘空间
}

// 获取系统信息
func GetSystemInfo(ctx VortexContext) error {
	// 查询系统的最大内存和当前使用的内存
	memory, err := mem.VirtualMemory()
	if nil != err {
		vortexUtil.Errorf("Vortex|GetSystemInfo|GetMemoryInfo|Error|%v", err)
		return HttpJsonResponse(ctx, http.StatusInternalServerError, map[string]interface{}{
			"message": "GetSystemInfo Error",
		})
	}

	var diskTotal uint64
	var diskUsed uint64
	if runtime.GOOS == "windows" {
		partitions, err := disk.Partitions(false)
		if err != nil {
			vortexUtil.Errorf("Vortex|GetSystemInfo|GetDiskInfo|Error|%v", err)
			return HttpJsonResponse(ctx, http.StatusInternalServerError, map[string]interface{}{
				"message": "GetSystemInfo Error",
			})
		}
		for _, partition := range partitions {
			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				vortexUtil.Errorf("Vortex|GetSystemInfo|GetDiskInfo|Error|%v", err)
				return HttpJsonResponse(ctx, http.StatusInternalServerError, map[string]interface{}{
					"message": "GetSystemInfo Error",
				})
			}
			diskTotal += usage.Total
			diskUsed += usage.Used
		}
	} else {
		diskUsage, err := disk.Usage("/")
		if nil != err {
			vortexUtil.Errorf("Vortex|GetSystemInfo|GetDiskInfo|Error|%v", err)
			return HttpJsonResponse(ctx, http.StatusInternalServerError, map[string]interface{}{
				"message": "GetSystemInfo Error",
			})
		}
		diskTotal = diskUsage.Total
		diskUsed = diskUsage.Used
	}
	return HttpJsonResponse(ctx, http.StatusOK, systemInfo{
		TotalMemory: memory.Total / vortexUtil.GB,
		UsedMemory:  memory.Used / vortexUtil.GB,
		TotalDisk:   diskTotal / vortexUtil.GB,
		UsedDisk:    diskUsed / vortexUtil.GB,
	})
}

// 检查服务是否正常
func checkAlive(ctx VortexContext) error {
	return HttpJsonResponse(ctx, http.StatusOK, vortexUtil.Map{
		"msg": "service success",
	})
}
