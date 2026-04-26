# Vault CLI
> CLI hacking mini-game from the Fallout franchise

# Install

```bash
# Replace x.x.x with the version number
sudo apt install ./vault-cli_x.x.x_amd64.deb
```

## Usage
```bash
vault-cli -d <difficulty>
```
`<difficulty>` is either `easy`, `medium`, or `hard`.

### Rolling text speed
With the `-s` flag, you can change the speed of the rolling text. The lower the number, the faster the text scrolls. The speed is in milliseconds
```bash
vault-cli -d <difficulty> -s <speed>
```


# Build yourself

1. Ensure you have [Go](https://go.dev/doc/install) and [Git](https://git-scm.com/install/) installed
2. ```git clone``` this repository
3. ```go install```
4. Be sure your `go/bin` folder is added to your `PATH`:
    
    a. Get your go folder using
    ```bash
    go env GOPATH
    ```
    b. It should return something like `/home/yourname/go`. If you exedcuted 3. then in this path in `/bin` you should find an executable called `vault-cli`
    c. Edit your bash config:
    ```bash
    nano ~/.bashrc
    ```
    d. Go to the last line and paste this line:
    ```bash
    export PATH=$PATH:$GOPATH/bin
    ```
    where `$GOPATH` is the path you got from a. and b. so for example:
    ```bash
    export PATH=$PATH:/home/USERNAME/go/bin
    ```
    e. Save and exit: Ctrl+O -> Enter -> Ctrl+X
    f. Reload your terminal!
