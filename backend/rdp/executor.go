package rdp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Connect gera o arquivo .rdp temporário na pasta Temp e executa o mstsc
func Connect(host, username, dpapiHex string) error {
	tempDir := os.TempDir()
	
	// Nome único para evitar conflitos se o usuário abrir 2 conexões rápidas
	fileName := fmt.Sprintf("rdpman_%d.rdp", time.Now().UnixNano())
	filePath := filepath.Join(tempDir, fileName)

	// Montando o arquivo com injeção segura (Password 51:b: aceita o hex da DPAPI)
	content := fmt.Sprintf("full address:s:%s\nusername:s:%s\npassword 51:b:%s\n", host, username, dpapiHex)

	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		return fmt.Errorf("falha ao criar arquivo temporário: %v", err)
	}

	// Mitigação de Path Hijacking: Forçar o caminho absoluto do sistema
	sysRoot := os.Getenv("SystemRoot")
	if sysRoot == "" {
		sysRoot = `C:\Windows` // Fallback seguro
	}
	mstscPath := filepath.Join(sysRoot, "System32", "mstsc.exe")

	cmd := exec.Command(mstscPath, filePath)
	if err := cmd.Start(); err != nil {
		_ = os.Remove(filePath)
		return fmt.Errorf("falha ao iniciar mstsc: %v", err)
	}

	// Fire and Forget: Goroutine limpa o arquivo após 5 segundos
	go func(path string) {
		time.Sleep(5 * time.Second)
		_ = os.Remove(path) // Apaga o arquivo temporário
	}(filePath)

	return nil
}