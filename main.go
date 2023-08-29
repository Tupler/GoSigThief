package main

import (
	"GoSigThief/pkg"
	"flag"
	"fmt"
)

var (
	TargetExe string
	CertExe   string
	OutPutExe string
)

func init() {
	flag.StringVar(&CertExe, "c", "", "Set the CertExe path")
	flag.StringVar(&TargetExe, "t", "", "Set the TarGetExe path")
	flag.StringVar(&OutPutExe, "o", "", "Set the OutPutExe path(default:Target_signed.exe)")

}

func banner() {
	bann := `

  ________           _________.__        __  .__    .__        _____ 
 /  _____/  ____    /   _____/|__| _____/  |_|  |__ |__| _____/ ____\
/   \  ___ /  _ \   \_____  \ |  |/ ___\   __\  |  \|  |/ __ \   __\ 
\    \_\  (  <_> )  /        \|  / /_/  >  | |   Y  \  \  ___/|  |   
 \______  /\____/  /_______  /|__\___  /|__| |___|  /__|\___  >__|   
        \/                 \/   /_____/           \/        \/
`
	fmt.Println(bann)
}
func main() {
	banner()
	flag.Parse()
	if CertExe == "" || TargetExe == "" {
		fmt.Println("input [-h] to get help")
		return
	}
	if OutPutExe == "" {
		OutPutExe = TargetExe[:4] + "_signed.exe"
	}
	pkg.WriteCert(pkg.CopyCert(CertExe), TargetExe, OutPutExe)
}
