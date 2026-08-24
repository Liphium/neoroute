package connector

import "github.com/tinylib/msgp/msgp"

type PackedAny struct {
	Value any
}

func (u *PackedAny) UnmarshalMsg(b []byte) ([]byte, error) {
	v, rest, err := msgp.ReadIntfBytes(b)
	if err != nil {
		return rest, err
	}

	u.Value = v
	return rest, nil
}

func (u PackedAny) MarshalMsg(b []byte) ([]byte, error) {
	return msgp.AppendIntf(b, u.Value)
}

func (u PackedAny) Msgsize() int {
	s, _ := msgp.AppendIntf(nil, u.Value)
	return len(s)
}

// streaming impl -> full msgp.Decodable/Encodable
func (u *PackedAny) DecodeMsg(dc *msgp.Reader) error {
	v, err := dc.ReadIntf()
	if err != nil {
		return err
	}
	u.Value = v
	return nil
}

func (u PackedAny) EncodeMsg(en *msgp.Writer) error {
	return en.WriteIntf(u.Value)
}
