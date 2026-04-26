# Vault CLI
A visually-pleasing terminal recreation of the iconic "Hacking" minigame from the Fallout franchise. Built with Go and the Charm.sh ecosystem.

---

## Installation

Find the latest binaries for your architecture on the [Releases](https://github.com/nicoewok/vault-cli/releases) page.

### Linux (Debian/Ubuntu)
Download the `.deb` file and install via dpkg:
```bash
sudo apt install ./vault-cli_1.0.0_linux_amd64.deb
```

### Windows
1. Download the latest `windows_amd64.zip`.

2. Extract the .zip file.

3. Run `vault-cli.exe` from your terminal (PowerShell or CMD).

Optionally add it to `PATH` using:
```powershell
[System.Environment]::SetEnvironmentVariable("Path", $env:Path + ";" + "[file_location]", "User")
```
where `[file_location]` is the location of the folder where your `vault-cli.exe` resides.

### MacOS
Download the latest `darwin_all.tar.gz` (or `arm64` for M1/M2/M3 chips).
1. Extract the binary:
```bash
tar -xf vault-cli_1.0.0_darwin_all.tar.gz
```
Move to your path:
```bash
mv vault-cli /usr/local/bin/
```

# Usage
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
