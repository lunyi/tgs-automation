package main

import (
	"fmt"
	"os/exec"
)

func helmInstall(chart, releaseName, namespace string) error {
	// 建立命令 "helm install"，並設置相關參數
	cmd := exec.Command("helm", "install", releaseName, chart, "--namespace", namespace)

	// 執行命令並獲取其輸出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("helm install failed with %s: %s", err, output)
	}

	fmt.Printf("Helm install output: %s\n", output)
	return nil
}
