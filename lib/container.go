package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type Container struct {
	Name   string
	Master [3]string
	Data   string
	Dir    string
}

func (c *Container) GetPath() string {
	return c.Dir + "/" + c.Name + ".cont.json"
}

func (c *Container) Write() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.GetPath(), b, 0600)
	if err != nil {
		return err
	}

	return nil
}

func OpenContainer(name, dir string) (*Container, error) {
	b, err := ioutil.ReadFile(dir + "/" + name + ".cont.json")
	if err != nil {
		return &Container{}, err
	}

	var cont Container
	err = json.Unmarshal(b, &cont)
	if err != nil {
		return &Container{}, err
	}

	return &cont, nil
}

func NewContainer(name, dir, password string) (*Container, error) {
	salt, err := GenerateSalt(16)
	if err != nil {
		return &Container{}, err
	}

	encSalt, err := GenerateSalt(16)
	if err != nil {
		return &Container{}, err
	}

	hash := A2Hash(password, salt)
	encHash := A2Hash(password, encSalt)

	emptyData, err := Encrypt("{}", encHash)
	if err != nil {
		return &Container{}, err
	}

	return &Container{
		Name:   name,
		Dir:    dir,
		Master: [3]string{BEncode(hash), BEncode(salt), BEncode(encSalt)},
		Data:   BEncode(emptyData),
	}, nil
}

func (c *Container) Read(password string) (map[string]string, error) {
	validHash, err := BDecode(c.Master[0])
	if err != nil {
		return nil, err
	}

	salt, err := BDecode(c.Master[1])
	if err != nil {
		return nil, err
	}

	newHash := A2Hash(password, salt)

	if reflect.DeepEqual(newHash, validHash) {
		encSalt, err := BDecode(c.Master[2])
		if err != nil {
			return nil, err
		}

		key := A2Hash(password, encSalt)

		bEnc, err := BDecode(c.Data)
		if err != nil {
			return nil, err
		}

		bData, err := Decrypt(bEnc, key)
		if err != nil {
			return nil, err
		}

		var data map[string]string
		err = json.Unmarshal(bData, &data)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	fmt.Println("sest: error: invalid password for container \"" + c.Name + "\"")
	os.Exit(1)
	return nil, nil
}

func (c *Container) SetData(newData map[string]string, password string) error {
	validHash, err := BDecode(c.Master[0])
	if err != nil {
		return err
	}

	salt, err := BDecode(c.Master[1])
	if err != nil {
		return err
	}

	newHash := A2Hash(password, salt)

	if reflect.DeepEqual(newHash, validHash) {
		encSalt, err := BDecode(c.Master[2])
		if err != nil {
			return err
		}

		key := A2Hash(password, encSalt)

		bData, err := json.Marshal(newData)
		if err != nil {
			return err
		}

		bEnc, err := Encrypt(string(bData), key)
		if err != nil {
			return err
		}

		c.Data = BEncode(bEnc)

		return nil
	}

	fmt.Println("sest: error: invalid password for container \"" + c.Name + "\"")
	os.Exit(1)
	return nil
}
