package vortex

import (
	"github.com/dzjyyds666/VortexCore/utils"
	"net/http"
	"runtime"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func prepareDefaultHttpRouter() []*httpRouter {
	return []*httpRouter{
		AppendHttpRouter([]string{http.MethodGet}, "/system/info", HandleGetSystemInfo, "获取系统信息", JwtSkipMw()),
		AppendHttpRouter([]string{http.MethodGet}, "/checkAlive", HandleCheckAlive, "检查服务是否正常", JwtSkipMw()),
	}
}

type systemInfo struct {
	TotalMemory uint64 `json:"totalMemory"` // 总内存
	UsedMemory  uint64 `json:"usedMemory"`  // 已使用的内存

	TotalDisk uint64 `json:"totalDisk"` // 磁盘总空间
	UsedDisk  uint64 `json:"usedDisk"`  // 已使用的磁盘空间
}

// 获取系统信息
func HandleGetSystemInfo(ctx VortexContext) error {
	echoCtx := ctx.GetEcho()
	// 查询系统的最大内存和当前使用的内存
	memory, err := mem.VirtualMemory()
	if nil != err {
		vUtil.Errorf("Vortex|HandleGetSystemInfo|GetMemoryInfo|Error|%v", err)
		return HttpJsonResponse(echoCtx, http.StatusInternalServerError, map[string]interface{}{
			"message": "HandleGetSystemInfo Error",
		})
	}

	var diskTotal uint64
	var diskUsed uint64
	if runtime.GOOS == "windows" {
		partitions, err := disk.Partitions(false)
		if err != nil {
			vUtil.Errorf("Vortex|HandleGetSystemInfo|GetDiskInfo|Error|%v", err)
			return HttpJsonResponse(echoCtx, http.StatusInternalServerError, map[string]interface{}{
				"message": "HandleGetSystemInfo Error",
			})
		}
		for _, partition := range partitions {
			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				vUtil.Errorf("Vortex|HandleGetSystemInfo|GetDiskInfo|Error|%v", err)
				return HttpJsonResponse(echoCtx, http.StatusInternalServerError, map[string]interface{}{
					"message": "HandleGetSystemInfo Error",
				})
			}
			diskTotal += usage.Total
			diskUsed += usage.Used
		}
	} else {
		diskUsage, err := disk.Usage("/")
		if nil != err {
			vUtil.Errorf("Vortex|HandleGetSystemInfo|GetDiskInfo|Error|%v", err)
			return HttpJsonResponse(echoCtx, http.StatusInternalServerError, map[string]interface{}{
				"message": "HandleGetSystemInfo Error",
			})
		}
		diskTotal = diskUsage.Total
		diskUsed = diskUsage.Used
	}
	return HttpJsonResponse(echoCtx, http.StatusOK, systemInfo{
		TotalMemory: memory.Total / vUtil.GB,
		UsedMemory:  memory.Used / vUtil.GB,
		TotalDisk:   diskTotal / vUtil.GB,
		UsedDisk:    diskUsed / vUtil.GB,
	})
}

// 检查服务是否正常
func HandleCheckAlive(ctx VortexContext) error {
	return HttpJsonResponse(ctx.GetEcho(), http.StatusOK, vUtil.Map{
		"msg": "service success",
	})
}
