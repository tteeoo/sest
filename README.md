# sest: secure strings
`sest` (pronounced "es ee es tee") is effectively a local command-line password manager, but can really be used to store any sensitive string of characters.

`sest` stores information in "containers", which are really just json files containing a password hash, salts, and some encrypted json.

Each container has its own master password which is used to access its contents. Depending on your setup, you may only end up using one container (which is fine).

A container stores data in key-value pairs, with a main value, and another optional value, used to store usernames.

The cli `sest` app works on Linux based systems and probably most other unix based systems (not tested).

## Installation
If you have Go installed then simply clone the repo and run `go install`.

Otherwise, a Linux binary is provided with the latest release on GitHub.

The default directory where containers are stored is `$HOME/.sest`, set the environment variable `SEST_DIR` to change this (no slash at the end).

In order for the `cp` command to copy the secret to your clipboard you will need `xclip` installed, and of course you'll need to be running Xorg for xclip to work.

## Usage
```
sest [-h | --help ] 
     [-V | --verison]
     [<command> [arguments]]
```

### Commands
```
ls                     lists all containers
mk  <container>        makes a new container, asks for a master password
ln  <container>        lists all keys in a container, asks for a master password
del <container>        deletes a container, asks for confirmation
in  <container> <key>  stores a new key-value pair in a container, asks for a master password and a value
cp  <container> <key>  copies the value of a key from a container to the clipboard (requires xclip), asks for a master password
rm  <container> <key>  removes a key-value pair from a container, asks for a master password
out <container> <key>  prints out the value of a key from a container, asks for a master password
```

## Security
To be frank, I am no cryptography expert, and one may find a flaw in this system. I'm (interperet this how you wish) 99% sure that it's perfectly fine for storing sensitive information.

So, here's how it works:
* A random salt is generated and used with your password in an Argon2id hash
* Another random salt is generated and also stored alongside the above two values (base64 encoded) in every container
* The data of each container is encrypted with AES-256 GCM using your (verified) password Argon2 hashed with the other salt as the key

## License
`sest` is licensed under the [BSD 2-clause license](https://github.com/tteeoo/sest/blob/master/LICENSE), use this program at your own risk; it offers no warranty for stolen information.
