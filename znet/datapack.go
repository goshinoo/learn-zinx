package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/goshinoo/learn-zinx/utils"
	"github.com/goshinoo/learn-zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	//DataLen uint32+ID uint32
	return 8
}

func (d *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	//创建存放byte字节的缓冲
	buffer := bytes.NewBuffer([]byte{})
	//将dataLen写进buffer
	if err := binary.Write(buffer, binary.LittleEndian, message.GetDataLen()); err != nil {
		return nil, err
	}
	//将msgId写进buffer
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMsgId()); err != nil {
		return nil, err
	}
	//将data写进buffer
	if err := binary.Write(buffer, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (d *DataPack) UnPack(data []byte) (ziface.IMessage, error) {
	//创建存放byte字节的reader
	buffer := bytes.NewReader(data)

	msg := &Message{}

	if err := binary.Read(buffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//判断dataLen是否超出允许最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		fmt.Println(utils.GlobalObject.MaxPackageSize, msg.DataLen)
		return nil, errors.New("too large msg data recv!")
	}

	if err := binary.Read(buffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	return msg, nil
}
