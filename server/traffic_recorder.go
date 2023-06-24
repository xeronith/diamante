package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	. "os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	. "github.com/xeronith/diamante/actor"
	. "github.com/xeronith/diamante/contracts/server"
)

type trafficRecorder struct {
	sync.RWMutex
	recording bool
	entries   []entry
}

type entry struct {
	timestamp uint64
	dataType  uint32
	token     string
	data      []byte
}

func NewTrafficRecorder() ITrafficRecorder {
	return &trafficRecorder{
		recording: false,
	}
}

func (recorder *trafficRecorder) Record(dataType uint32, input ...interface{}) {
	if !recorder.recording {
		return
	}

	timestamp := uint64(time.Now().UnixNano())

	var token string
	var data []byte
	switch dataType {
	case TEXT_REQUEST, TEXT_RESULT:
		data = []byte(input[0].(string))
	case BINARY_REQUEST, BINARY_RESULT:
		data = input[0].([]byte)
	}

	recorder.Lock()
	defer recorder.Unlock()

	recorder.entries = append(recorder.entries, entry{
		timestamp: timestamp,
		dataType:  dataType,
		token:     token,
		data:      data,
	})
}

func (recorder *trafficRecorder) Start() {
	recorder.Lock()
	defer recorder.Unlock()

	if recorder.recording {
		return
	}

	recorder.recording = true
	recorder.entries = make([]entry, 0)
}

func (recorder *trafficRecorder) Stop() error {
	recorder.Lock()
	defer recorder.Unlock()

	if !recorder.recording {
		return nil
	}

	recorder.recording = false

	err := MkdirAll("./replays", ModePerm)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("./replays/SERVER_TRAFFIC_%d.dump", time.Now().UnixNano())
	file, err := OpenFile(fileName, O_WRONLY|O_CREATE|O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	size := len(recorder.entries)
	for i := 0; i < size; i++ {
		entry := recorder.entries[i]
		recorder.writeUInt64(file, entry.timestamp)
		recorder.writeUInt32(file, entry.dataType)
		recorder.writeString(file, entry.token)
		recorder.writeByteArray(file, entry.data)
	}

	recorder.entries = make([]entry, 0)

	runtime.GC()
	debug.FreeOSMemory()

	return nil
}

func (recorder *trafficRecorder) Load(dumpFile string) error {
	if recorder.recording {
		return errors.New("replay failed while recording")
	}

	file, err := Open(dumpFile)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	recorder.entries = make([]entry, 0)
	for {
		buffer := make([]byte, 8)
		count, err := file.Read(buffer)
		if count == 0 {
			break
		}

		if count != 8 || err != nil {
			return errors.New("replay file is invalid or corrupt")
		}

		timestamp := binary.BigEndian.Uint64(buffer)

		buffer = make([]byte, 4)
		count, err = file.Read(buffer)
		if count != 4 || err != nil {
			return errors.New("replay file is invalid or corrupt")
		}

		dataType := binary.BigEndian.Uint32(buffer)

		count, err = file.Read(buffer)
		if count != 4 || err != nil {
			return errors.New("replay file is invalid or corrupt")
		}

		dataLength := binary.BigEndian.Uint32(buffer)
		buffer = make([]byte, dataLength)
		count, err = file.Read(buffer)
		if uint32(count) != dataLength || err != nil {
			return errors.New("replay file is invalid or corrupt")
		}

		token := string(buffer)

		buffer = make([]byte, 4)
		count, err = file.Read(buffer)
		if count != 4 || err != nil {
			return errors.New("replay file is invalid or corrupt")
		}

		dataLength = binary.BigEndian.Uint32(buffer)
		buffer = make([]byte, dataLength)
		count, err = file.Read(buffer)
		if uint32(count) != dataLength || err != nil {
			return errors.New("replay file is invalid or corrupt")
		}

		recorder.entries = append(recorder.entries, entry{
			timestamp: timestamp,
			dataType:  dataType,
			token:     token,
			data:      buffer,
		})
	}

	return nil
}

func (recorder *trafficRecorder) Replay(server IServer, speed float32) error {
	if recorder == nil || len(recorder.entries) == 0 {
		return nil
	}

	offset := recorder.entries[0].timestamp / 1000000
	interval := time.Duration(1000000 * (1 / speed))
	ticker := time.NewTicker(interval)
	size := len(recorder.entries)

	fmt.Println("\n\nSERVER TRAFFIC REPLAY")
	fmt.Println("──────────────────────────────────────")

	for range ticker.C {
		for index := 0; index < size; index++ {
			entry := recorder.entries[index]
			timestamp := entry.timestamp / 1000000

			if timestamp == offset {
				switch entry.dataType {
				//---------------------------------------------------------------------------------------------------
				case BINARY_REQUEST:
					fmt.Printf("▒▒ BO/I ▒▒ %08X %d % X\n", index+1, entry.timestamp, entry.data)
					actor := CreateActor(nil, false, "", "")
					actor.SetToken(entry.token)
					server.OnActorBinaryData(actor, entry.data)
				case BINARY_RESULT:
					fmt.Printf("▒▒ BO/O ▒▒ %08X %d % X\n", index+1, entry.timestamp, entry.data)
				//---------------------------------------------------------------------------------------------------
				case TEXT_REQUEST:
					fmt.Printf("▓▓ TO/I ▓▓ %08X %d % X\n", index+1, entry.timestamp, entry.data)
					actor := CreateActor(nil, false, "", "")
					actor.SetToken(entry.token)
					server.OnActorTextData(actor, string(entry.data))
				case TEXT_RESULT:
					fmt.Printf("▓▓ TO/O ▓▓ %08X %d % X\n", index+1, entry.timestamp, entry.data)
				}
				//---------------------------------------------------------------------------------------------------
				if index == size-1 {
					fmt.Printf("──────────────────────────────────────\nTOTAL OPERATIONS: %d\n\n\n\n", size)
					ticker.Stop()
				}
			}
		}

		offset++
	}

	return nil
}

func (recorder *trafficRecorder) writeUInt64(file *File, value uint64) {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, value)
	if _, err := file.Write(buffer); err != nil {
		fmt.Println(err)
	}
}

func (recorder *trafficRecorder) writeUInt32(file *File, value uint32) {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, value)
	if _, err := file.Write(buffer); err != nil {
		fmt.Println(err)
	}
}

func (recorder *trafficRecorder) writeString(file *File, value string) {
	length := uint32(len(value))
	recorder.writeUInt32(file, length)
	if _, err := file.Write([]byte(value)); err != nil {
		fmt.Println(err)
	}
}

func (recorder *trafficRecorder) writeByteArray(file *File, buffer []byte) {
	length := uint32(len(buffer))
	recorder.writeUInt32(file, length)
	if _, err := file.Write(buffer); err != nil {
		fmt.Println(err)
	}
}
