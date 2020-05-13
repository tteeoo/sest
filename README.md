# sest
## sest: secure strings

`sest` is effectively a local command-line password manager, but can really be used to store any sensitive string of characters.

`sest` stores information in "containers", which are really just json files containing a password hash, salts, and some encrypted json.

Each container has its own master password which is used to access its contents. Depending on your setup, you may only end up using one container (which is fine).

A container stores data in key-value pairs.

`sest` is compatible with pretty much all UNIXes (technically only tested on Linux, make an issue if it's not working for you)

## Installation:
If you have go installed, simply run `go install` [(install go here)](https://golang.org/doc/install#install)

Otherwise, a linux binary is provided in the `bin/` directory (compiled on arch btw)

## Usage:
sest [--version | -V] | [--help | -h] | [<command> [arguments]]

### Commands:
**mk (container name):** makes a new container, will ask for a master password

**del (container name):** deletes a container, will ask for confirmation

**ls:** lists all containers

**in (container name) (key name):** stores a new key-value pair in a container, will ask for a master password and a value

**out (container name) (key name):** prints out the value of a key from a container, will ask for a master password

**ln:** lists all keys in a container, will ask for a master password

**rm (container name) (key name):** removes a key-value pair from a container, will ask for a master password

## Security
To be frank, I am no cryptography expert, and one may find a flaw in this system (as such I, nor any other contributers are responsible for stolen data), although I trust this program and I'm 99% sure that it's perfectly fine for storing any sensitive information.

So, here's how it works:

* A random salt is generated and used with your password in an Argon2id hash
* Another random salt is generated and also stored alongside the above two values (base64 encoded) in every container
* The data of each container is encrypted with AES-256 GCM using your (verified) password Argon2 hashed with the other salt as the key

I know, the salts aren't necassarily needed, but they make hashing everything a bit easier and don't have any negative side effects.

## License

`sest` is licensed under the [BSD 2-clause license](https://github.com/tteeoo/sest/blob/master/LICENSE), use this product at your own risk.
