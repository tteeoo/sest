# sest: secure strings

`sest` is effectively a local command-line password manager, but can really be used to store any sensitive string of characters.

`sest` stores information in "containers", which are really just json files containing a password hash, salts, and some encrypted json.

Each container has its own master password which is used to access its contents. Depending on your setup, you may only end up using one container (which is fine).

A container stores data in key-value pairs.

`sest` works on Linux based OSes and probably most UNIXes

![usage gif](https://raw.githubusercontent.com/tteeoo/sest/master/usage.gif)

## Installation
If you have Go installed [(install Go here)](https://golang.org/doc/install#install), simply clone the repo and run `go install` 

Otherwise, a Linux binary is provided in the `bin/` directory (compiled on arch btw)

To quickly install, run:
```
# wget https://github.com/tteeoo/sest/releases/download/0.1.5/sest -P /usr/bin
```

The default directory where containers are stored is `$HOME/.sest`, set the environment variable `SEST_DIR` to change this (no slash at the end).

In order for the `cp` command to copy the secret to your clipboard you will need `xclip` installed, and of course you'll need to be running Xorg for xclip to work.

## Usage
`sest [--version | -V] | [--help | -h] | [<command> [arguments]]`

### Commands
```
ls                     lists all containers
mk  <container>        makes a new container, will ask for a master password
ln  <container>        lists all keys in a container, will ask for a master password
del <container>        deletes a container, will ask for confirmation
in  <container> <key>  stores a new key-value pair in a container, will ask for a master password and a value
cp  <container> <key>  copies the value of a key from a container to the clipboard (needs xclip installed), will ask for a master password
rm  <container> <key>  removes a key-value pair from a container, will ask for a master password
out <container> <key>  prints out the value of a key from a container, will ask for a master password
```

## Security
To be frank, I am no cryptography expert, and one may find a flaw in this system (as such I, nor any other contributers are responsible for any stolen data), although (interperet this how you wish) I'm 99% sure that it's perfectly fine for storing sensitive information.

So, here's how it works:

* A random salt is generated and used with your password in an Argon2id hash
* Another random salt is generated and also stored alongside the above two values (base64 encoded) in every container
* The data of each container is encrypted with AES-256 GCM using your (verified) password Argon2 hashed with the other salt as the key

## License

`sest` is licensed under the [BSD 2-clause license](https://github.com/tteeoo/sest/blob/master/LICENSE), use this product at your own risk.
