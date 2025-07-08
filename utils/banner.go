package vUtil

import "fmt"

func ShowBanner(port string) {
	Infof(fmt.Sprintf(`
 __     __                 _                 
 \ \   / /   ___    _ __  | |_    ___  __  __
  \ \ / /   / _ \  | '__| | __|  / _ \ \ \/ /
   \ V /   | (_) | | |    | |_  |  __/  >  < 
    \_/     \___/  |_|     \__|  \___| /_/\_\
----------------------------------------------
 Vortex Server Start Success on %s
==============================================
`, port))
}
