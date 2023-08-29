package pkg

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func getWord(file *os.File) uint32 {
	fil := make([]byte, 2)
	at, err := file.Read(fil)
	if err != nil || at == 0 {
		log.Fatal(err.Error())
	}
	return uint32(binary.LittleEndian.Uint16(fil))

}
func getDword(file *os.File) uint32 {
	fil := make([]byte, 4)
	at, err := file.Read(fil)
	if err != nil || at == 0 {
		log.Fatal(err.Error())
	}
	return binary.LittleEndian.Uint32(fil)
}

func GetPeInfo(path string) (int64, uint32, uint32) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	fmt.Println("[*] Got the File")
	_, err = file.Seek(0x3c, 0)
	if err != nil {
		log.Fatal(err)
	}

	peHeaderLocation := getDword(file)
	CoffStart := int64(peHeaderLocation) + 4
	OptionalheaderStart := CoffStart + 20
	_, err = file.Seek(OptionalheaderStart, 0)
	if err != nil {
		log.Fatal(err.Error())
	}
	Magic := getWord(file)
	_, err = file.Seek(OptionalheaderStart+24, 0)
	if err != nil {
		log.Fatal(err.Error())
	}
	//	var imgBase uint64
	if Magic != 0x20b {
		file.Seek(4, io.SeekCurrent)

	}

	if Magic != 0x20b {
		file.Seek(4, io.SeekCurrent)
		//imgBase = uint64(getDword(file))
	} else {
		file.Seek(8, io.SeekCurrent)
		//imgBase = getQword(file)
	}
	position, _ := file.Seek(0, io.SeekCurrent)

	file.Seek(position+40, 0)

	if Magic == 0x20b {
		file.Seek(32, io.SeekCurrent)
	} else {
		file.Seek(16, io.SeekCurrent)
	}

	CertTableLOC, _ := file.Seek(40, io.SeekCurrent)
	fmt.Println("[*] Got the CertTable")
	CertLOC := getDword(file)
	CertSize := getDword(file)

	return CertTableLOC, CertLOC, CertSize
}

func CopyCert(path string) []byte {
	_, CertLOC, CertSize := GetPeInfo(path)
	if CertSize == 0 || CertLOC == 0 {
		log.Fatal("[*] Input file Not signed! ")
	}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	file.Seek(int64(CertLOC), 0)
	cert := make([]byte, CertSize)
	file.Read(cert)
	fmt.Println("[*] Read the Cert successfully")
	return cert
}

func WriteCert(cert []byte, path string, outputPath string) {
	CertTableLOC, _, _ := GetPeInfo(path)
	//	copyFile(path, outputPath)
	file1, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	file2, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file2.Close()
	defer file1.Close()

	file1Info, err := os.Stat(path)
	file1Len := file1Info.Size()
	file1data := make([]byte, file1Len)
	file1.Read(file1data)
	file2.Write(file1data)
	file2.Seek(CertTableLOC, 0)
	x := make([]byte, 4)
	binary.LittleEndian.PutUint32(x, uint32(file1Len))
	file2.Write(x)
	bCertLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(bCertLen, uint32(len(cert)))
	file2.Write(bCertLen)
	file2.Seek(0, io.SeekEnd)
	file2.Write(cert)
	fmt.Println("[*] Signature appended!")
}
