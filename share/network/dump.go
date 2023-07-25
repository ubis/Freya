package network

import (
	"bytes"
	"fmt"

	"github.com/ubis/Freya/share/log"
)

// DumpPacket function dumps packet into console
func DumpPacket(packet interface{}) {
	if writer, ok := packet.(*Writer); ok {
		dump(writer.buffer, writer.index)
	} else if reader, ok := packet.(*Reader); ok {
		dump(reader.buffer, int(reader.Size))
	} else if array, ok := packet.([]byte); ok {
		dump(array, len(array))
	} else {
		log.Error("Unknown packet type!")
	}
}

// Dumps byte array into console
func dump(packet []byte, length int) {
	var buffer bytes.Buffer

	buffer.WriteString("\n\n")
	buffer.WriteString("      00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n")
	buffer.WriteString("-----------------------------------------------------\n")

	var newLine = 16
	var lineCnt = 0

	for i := 0; i < length; i++ {
		if newLine == 16 {
			newLine = 0
			if i > 0 {
				buffer.WriteString("\n")
			}

			buffer.WriteString(fmt.Sprintf("%04X: ", lineCnt))
			lineCnt++
		}

		buffer.WriteString(fmt.Sprintf("%02X ", packet[i]))
		newLine++
	}

	buffer.WriteString("\n\n")
	log.Debug(buffer.String())
}
