package crypto

import (
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modcrypt32              = syscall.NewLazyDLL("crypt32.dll")
	procCryptProtectData    = modcrypt32.NewProc("CryptProtectData")
	procCryptUnprotectData  = modcrypt32.NewProc("CryptUnprotectData")
)

const cryptprotectUIForbidden = 0x1

type dataBlob struct {
	cbData uint32
	pbData *byte
}

// Encrypt protege a string usando DPAPI e retorna o hash em formato hexadecimal.
func Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	// O mstsc.exe exige que a senha esteja em formato UTF-16LE 
	// antes de ser criptografada pela DPAPI.
	utf16pw, err := syscall.UTF16FromString(plainText)
	if err != nil {
		return "", fmt.Errorf("falha ao converter string para utf16: %v", err)
	}

	// Convertendo o slice de uint16 (UTF-16) para um array de bytes brutos
	var ptBytes []byte
	for _, w := range utf16pw {
		ptBytes = append(ptBytes, byte(w), byte(w>>8))
	}

	var dataIn dataBlob
	dataIn.cbData = uint32(len(ptBytes))
	dataIn.pbData = &ptBytes[0]

	var dataOut dataBlob

	ret, _, errCall := procCryptProtectData.Call(
		uintptr(unsafe.Pointer(&dataIn)),
		0, 0, 0, 0,
		uintptr(cryptprotectUIForbidden),
		uintptr(unsafe.Pointer(&dataOut)),
	)

	if ret == 0 {
		return "", fmt.Errorf("falha ao criptografar: %v", errCall)
	}
	defer syscall.LocalFree(syscall.Handle(unsafe.Pointer(dataOut.pbData)))

	outBytes := unsafe.Slice(dataOut.pbData, dataOut.cbData)
	return hex.EncodeToString(outBytes), nil
}

// Decrypt recebe o hexadecimal da DPAPI e retorna a string em texto plano.
func Decrypt(hexString string) (string, error) {
	if hexString == "" {
		return "", nil
	}

	encBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return "", fmt.Errorf("hexadecimal inválido: %v", err)
	}

	var dataIn dataBlob
	dataIn.cbData = uint32(len(encBytes))
	dataIn.pbData = &encBytes[0]

	var dataOut dataBlob

	ret, _, callErr := procCryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&dataIn)),
		0, 0, 0, 0,
		uintptr(cryptprotectUIForbidden),
		uintptr(unsafe.Pointer(&dataOut)),
	)

	if ret == 0 {
		return "", fmt.Errorf("falha ao descriptografar (banco movido de máquina?): %v", callErr)
	}
	defer syscall.LocalFree(syscall.Handle(unsafe.Pointer(dataOut.pbData)))

	outBytes := unsafe.Slice(dataOut.pbData, dataOut.cbData)

	// Transformando os bytes (UTF-16LE) de volta para o formato de array numérico do Windows
	u16s := make([]uint16, len(outBytes)/2)
	for i := 0; i < len(outBytes); i += 2 {
		u16s[i/2] = uint16(outBytes[i]) | (uint16(outBytes[i+1]) << 8)
	}

	// Converte de volta para UTF-8 nativo do Go e remove os bytes nulos automaticamente
	return syscall.UTF16ToString(u16s), nil
}