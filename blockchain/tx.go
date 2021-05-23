package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gochain/wallet"
	"math/big"
)



type TxInput struct {
	PrevTxID []byte // 32 bytes little-endian
	Out []byte // 4 bytes little-endian
	ScriptSig []byte
	Sequence []byte // 4 bytes little-endian
}
func (in *TxInput) NewInput(prevTx,prevIndex,scriptSig,sequence []byte) {
	in.PrevTxID = prevTx
	in.Out = prevIndex
	if scriptSig == nil{
		in.ScriptSig = Script()
	}else{
		in.ScriptSig = scriptSig
	}
	in.Sequence = sequence
}
func (in TxInput) Serialize() []byte{
	result := toLittleEndian(in.PrevTxID,len(in.PrevTxID))
	result = append(result, toLittleEndian(in.Out,4)...)
	result = append(result, in.ScriptSig...)
	result = append(result, toLittleEndian(in.Sequence,4)...)
	
	return result
}

func Script()[]byte{
	return nil
}

type TxOutput struct{
	Amount *big.Int
	ScriptPubKey []byte
}
func (out TxOutput) Serialize()[]byte{
	amount := out.Amount.Bytes()
	result := toLittleEndian(amount,8)
	result = append(result, out.ScriptPubKey...)

	return result
}

func DeserializeOutputs(data []byte) TxOutputs{
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	Handle(err)
	return outputs
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool{
	lockingHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Compare(lockingHash,pubKeyHash) == 0
}

func (out *TxOutput) IsLockedWithKey(scriptPubKey string) bool{
	return out.scriptPubKey == scriptPubKey	
}

func toLittleEndian(bytes []byte, length int) []byte{
	le := make([]byte,length)
	for i := len(le)-1;i >= 0;i--{
		if bytes[i] != 0x00{
			le = append(le, bytes[i])
		}
		le = append(le, 0x00)
	}
	return le
}