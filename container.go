package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type container struct {
	name   string
	master [3]string
	data   string
}

func (c container) getPath() string {
	return contDir + "/" + c.name + ".cont.json"
}

func (c container) write() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.getPath(), b, 0700)
	if err != nil {
		return err
	}

	return nil
}

func newContainer(name, password string) (container, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return container{}, err
	}

	encSalt, err := generateSalt(16)
	if err != nil {
		return container{}, err
	}

	hash := a2Hash(password, salt)
	encHash := a2Hash(password, encSalt)

	emptyData, err := encrypt("{}", encHash)
	if err != nil {
		return container{}, err
	}

	return container{
		name:   name,
		master: [3]string{bEncode(hash), bEncode(salt), bEncode(encSalt)},
		data:   bEncode(emptyData),
	}, nil
}

func (c container) getData(password string) (map[string]string, error) {
	validHash, err := bDecode(c.master[0])
	if err != nil {
		return nil, err
	}

	salt, err := bDecode(c.master[1])
	if err != nil {
		return nil, err
	}

	newHash := a2Hash(password, salt)

	if reflect.DeepEqual(newHash, validHash) {
		encSalt, err := bDecode(c.master[2])
		if err != nil {
			return nil, err
		}

		key := a2Hash(password, encSalt)

		bEnc, err := bDecode(c.data)
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

	fmt.Println("sest: error: invalid password for container", c.name)
	os.Exit(1)
	return nil, nil
}

func (c container) setData(newData map[string]string, password string) error {
	validHash, err := bDecode(c.master[0])
	if err != nil {
		return err
	}

	salt, err := bDecode(c.master[1])
	if err != nil {
		return err
	}

	newHash := a2Hash(password, salt)

	if reflect.DeepEqual(newHash, validHash) {
		encSalt, err := bDecode(c.master[2])
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

		c.data = bEncode(bEnc)

		return nil
	}

	fmt.Println("sest: error: invalid password for container", c.name)
	os.Exit(1)
	return nil
}
