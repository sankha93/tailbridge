package tailbridge

import (
    "io/ioutil"
    "regexp"
    "strings"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Listen int
    Groups map[string]Group
}

type Group struct {
    User string
    Port int
    Machines []string
    Directories []string
}

// Machine-group map
var machines map[string]string
var config Config

func ReadConfig(config_path string) Config {
    yamlFile, err := ioutil.ReadFile(config_path)

    if err != nil {
       panic(err)
    }

    err = yaml.Unmarshal(yamlFile, &config)

    if err != nil {
        panic(err)
    }

    BuildMachinesIndex(config.Groups)
    return config
}

func BuildMachinesIndex(groups map[string]Group) {
    machines = make(map[string]string)

    for group_name := range groups {
        for _, ip := range groups[group_name].Machines {
            machines[ip] = group_name
        }
    }
}

func GetMachineParams(ip string) (user string, port int, succes bool) {
    group_name, ok := machines[ip]

    if ok {
        port := config.Groups[group_name].Port
        user := config.Groups[group_name].User
        return user, port, ok
    } else {
        return "", 0, ok
    }
}

func IsFileAllowed(file_name string, machine_ip string) bool {
    illegalSet := []string{"|", ">", "<", ";", "$", "&"}
    
    var contains bool
    for _, c := range illegalSet {
        contains = strings.Contains(file_name, c)

        if contains {
            return false
        }
    }

    group_name, ok := machines[machine_ip]

    if !ok {
        return false
    }

    dirs := config.Groups[group_name].Directories

    var match bool
    for _, dir := range dirs {
        match, _ = regexp.MatchString(dir, file_name)
        if match {
            return true
        }
    }
    return false
}

