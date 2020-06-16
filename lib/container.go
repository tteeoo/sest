package lib

import (
	"errors"
	"encoding/json"
	"io/ioutil"
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
	salt, err := generateSalt(16)
	if err != nil {
		return &Container{}, err
	}

	encSalt, err := generateSalt(16)
	if err != nil {
		return &Container{}, err
	}

	hash := a2Hash(password, salt)
	encHash := a2Hash(password, encSalt)

	emptyData, err := encrypt("{}", encHash)
	if err != nil {
		return &Container{}, err
	}

	return &Container{
		Name:   name,
		Dir:    dir,
		Master: [3]string{bEncode(hash), bEncode(salt), bEncode(encSalt)},
		Data:   bEncode(emptyData),
	}, nil
}

func (c *Container) GetData(password string) (map[string]string, error) {
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

	return nil, errors.New("invalid password for container \"" + c.Name + "\"")
}

func (c *Container) SetData(newData map[string]string, password string) error {
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

	return errors.New("invalid password for container \"" + c.Name + "\"")
}
