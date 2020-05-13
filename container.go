package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type container struct {
	Name   string
	Master [3]string
	Data   string
}

func (c *container) getPath() string {
	return contDir + "/" + c.Name + ".cont.json"
}

func (c *container) write() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.getPath(), b, 0600)
	if err != nil {
		return err
	}

	return nil
}

func openContainer(name string) (*container, error) {
	b, err := ioutil.ReadFile(contDir + "/" + name + ".cont.json")
	if err != nil {
		return &container{}, err
	}

	var cont container
	err = json.Unmarshal(b, &cont)
	if err != nil {
		return &container{}, err
	}

	return &cont, nil
}

func newContainer(name, password string) (*container, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return &container{}, err
	}

	encSalt, err := generateSalt(16)
	if err != nil {
		return &container{}, err
	}

	hash := a2Hash(password, salt)
	encHash := a2Hash(password, encSalt)

	emptyData, err := encrypt("{}", encHash)
	if err != nil {
		return &container{}, err
	}

	return &container{
		Name:   name,
		Master: [3]string{bEncode(hash), bEncode(salt), bEncode(encSalt)},
		Data:   bEncode(emptyData),
	}, nil
}

func (c *container) getData(password string) (map[string]string, error) {
	validHash, err := bDecode(c.Master[0])
	if err != nil {
		return nil, err
	}

	salt, err := bDecode(c.Master[1])
	if err != nil {
		return nil, err
	}

	newHash := a2Hash(password, salt)

	if reflect.DeepEqual(newHash, validHash) {
		encSalt, err := bDecode(c.Master[2])
		if err != nil {
			return nil, err
		}

		key := a2Hash(password, encSalt)

		bEnc, err := bDecode(c.Data)
		if err != nil {
			return nil, err
		}

		bData, err := decrypt(bEnc, key)
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

func (c *container) setData(newData map[string]string, password string) error {
	validHash, err := bDecode(c.Master[0])
	if err != nil {
		return err
	}

	salt, err := bDecode(c.Master[1])
	if err != nil {
		return err
	}

	newHash := a2Hash(password, salt)

	if reflect.DeepEqual(newHash, validHash) {
		encSalt, err := bDecode(c.Master[2])
		if err != nil {
			return err
		}

		key := a2Hash(password, encSalt)

		bData, err := json.Marshal(newData)
		if err != nil {
			return err
		}

		bEnc, err := encrypt(string(bData), key)
		if err != nil {
			return err
		}

		c.Data = bEncode(bEnc)

		return nil
	}

	fmt.Println("sest: error: invalid password for container \"" + c.Name + "\"")
	os.Exit(1)
	return nil
}
