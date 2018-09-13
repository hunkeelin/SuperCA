### Introduction ###
The is the super CA written in Go. Can handle a lot of csr request at the same time. 

### Feature ###
- low memory footprint
- fine grain access control 
- super fast 
- Security is conformed to current standard

### Usage ###
I've included the binary built for linux. However, if you want to compiled it your self following the below instructions: 

1. download latest version of `go`
2. clone repo
```
run build.sh
./SuperCA -h
./SuperCA -server
```
- The software will automatically generate rootca.crt in your `$work/CA/certs` and `$work/CA/keys` directory upon startup. However, you can replace rootca.crt and rootca.key with your own intermca you got from outside CA like symantec or godaddy. 

### Public certs ###
- All certs files are located in https://$hostname:$port/cacerts

### Config file ###
- The documentation of CA.conf is pretty much self explanationary in CA.conf_template. Only thing that needs to take caution is workdir. The workdir specifys the directory to do fine grain control. For example let's say you set workdir=/tmp

When a request come in, it will search for the hostname and do a reverse hostname lookup and split the directory by "." and read the configuration accordingly. By default it will try to read the config /tmp/config. The config should follow the following format:

```
// The name of the ca to sign any csr request. All the cacert name should be located in $capath in CA.conf. E.G devca,prodca you can have multiple signca seperate by the comma. Request with header "signca: $ca_name"
signca=
// Same as above but for keys. E.G devca.key
isca=
// Whether it is even allow to sign any request at all. True/False (Default: False)
allow=
// How many days the .crt is vald for. 
validdays=
```
The software will recursively search for config file base on hostname. The config file that is more fine grained to the specific host will be read. Example, let's say test1.production.abc.com send a csr request to SuperCA server. SuperCA will first look for the config `/tmp/com/abc/production/test1/config`, if it doesn't exist it will try to look `/tmp/com/abc/production/config`, and slowly recursive out to /tmp/config. This way, you can choose which CA sign which CSR down to the exact host/environment/organization. 

### USAGE ###
Send request to SuperCA inform of. 
```
Header {
    "signca": $ca // specify which CA you want your csr to be signed. Will be doubel checked in the "signca" field in the config on the server side. 
}

body {
    csr in beytes pem encoded. 
}
```
### Tricks ###
With the example above, you can make a `config` file and put it in the hostname that is hosting the SuperCA server and set `isca` to `true`. Then use the SuperCAclient (https://github.com/hunkeelin/supercaclient) and request certificates. That way you created your own intermediate CA for singing, since you don't want everything to be signed by rootca.crt. For how to use the SuperCAClient, visit the page for documentation.
### Notes ###
For feature request ping me in reddit


