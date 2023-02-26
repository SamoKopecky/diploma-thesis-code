package app

import (
	"io"
	"os"
	"strings"

	myio "github.com/SamoKopecky/pqcom/main/io"
	"github.com/SamoKopecky/pqcom/main/network"
	log "github.com/sirupsen/logrus"
)

func Chat(destAddr string, srcPort, destPort int, connect bool) {
	if connect {
		stream := network.Connect(destAddr, destPort)
		go printer(stream, false)
		for {
			data := []byte(myio.ReadUserInput(""))
			stream.Send(data, network.ContentT)
		}
	} else {
		streamFactory := make(chan network.Stream)
		go network.Listen(srcPort, streamFactory, false)
		stream := <-streamFactory
		go printer(stream, false)
		for {
			data := []byte(myio.ReadUserInput(""))
			stream.Send(data, network.ContentT)
		}
	}
}

func Send(destAddr string, srcPort, destPort int, filePath string) {
	stream := network.Connect(destAddr, destPort)
	chunks := make(chan []byte)
	var source io.Reader
	var err error
	var fileName string

	if filePath != "" {
		splitFilePath := strings.Split(filePath, string(os.PathSeparator))
		fileName = splitFilePath[len(splitFilePath)-1]
		source, err = os.Open(filePath)
		if err != nil {
			log.WithField("error", err).Error("Error opening file")
		}
	} else {
		source = os.Stdin
	}
	go func() {
		myio.ReadByChunks(source, chunks, network.CHUNK_SIZE)
		close(chunks)
	}()
	if fileName != "" {
		stream.Send([]byte(fileName), network.FileNameT)
	}
	for msg := range chunks {
		stream.Send(msg, network.ContentT)
	}
	log.WithField("addr", stream.Conn.RemoteAddr()).Info("Done sending")
}

func Receive(destAddr string, srcPort, destPort int, dir string) {
	streamFactory := make(chan network.Stream)
	go network.Listen(srcPort, streamFactory, true)
	for {
		stream := <-streamFactory
		if dir != "" {
			go dirFileWriter(stream.Msg, dir)
		} else {
			go printer(stream, true)
		}
	}
}
