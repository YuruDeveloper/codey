package browser

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	appError "github.com/YuruDeveloper/codey/internal/error"
)


func Browser(url string) error {
	// WSL 환경 감지
	if isWSL() {
		// WSL에서는 Windows의 rundll32를 사용하여 브라우저 열기
		// 이 방법이 URL 특수문자를 가장 안전하게 처리함
		cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url)
		err := cmd.Start()
		if err != nil {
			return appError.NewError(appError.FailOpenBrowser,err)
		}
		return nil
	}

	// 일반 Linux/Mac/Windows
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return appError.NewValidError(appError.UnexpectedOS,"Fail open Browser")
	}
	err := cmd.Start()
	if err != nil {
		return appError.NewError(appError.FailOpenBrowser,err)
	}
	return nil
}

// isWSL은 WSL 환경인지 확인합니다
func isWSL() bool {
	// /proc/version에 "microsoft" 또는 "WSL"이 있으면 WSL
	if data, err := os.ReadFile("/proc/version"); err == nil {
		version := strings.ToLower(string(data))
		if strings.Contains(version, "microsoft") || strings.Contains(version, "wsl") {
			return true
		}
	}
	return false
}