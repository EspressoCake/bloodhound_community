# 🔭 Bottom-Line Up Front

This is a generator for Bloodhound Community Edition containerized services.
It generates unique credentials for database instances, alongside dynamic port choices.


## 🤔 Why?

- Manually looking through `Docker`'s `ps -a`/etc isn't optimal
- Finding ports for multiple instances sucks
- Managing credentials sucks


## 👷 Building
Assuming you `Go` version >= 1.24: `go install github.com/EspressoCake/bloodhound_community@latest`


## 🚀 Usage

#### General Help
```
bloodhound_community --help
Generation utility for creating Bloodhound instances and fetching metadata

Usage:
  bloodhound_community [command]

Available Commands:
  generate    Create configurations for projects
  help        Help about any command
  ports       Retrieve ports/connection string for running containers

Flags:
  -h, --help   help for bloodhound_community

Use "bloodhound_community [command] --help" for more information about a command.
```

#### Generation Help
```
bloodhound_community generate --help
Generates respective configruation(s) for Docker instances to deploy Bloodhound Community Edition via docker-compose.yml files

Usage:
  bloodhound_community generate [flags]

Flags:
  -h, --help          help for generate
  -n, --name string   Name for project, in lowercase
  -p, --path string   Filepath on system to desired root directory. Default will be the current working directory
```

#### Idenfitying Ports of a Deployed Instance
```
bloodhound_community ports --help
Retrieve ports/connection string for running containers via SSH or assumed to be local

Usage:
  bloodhound_community ports [flags]

Flags:
  -h, --help            help for ports
  -l, --local           Boolean if the connection is meant to be local or remote (default true)
  -p, --prefix string   Prefix for project containers to query
```

## 🌱 Example Usage
```
bloodhound_community generate -n caspian

Password for Bloodhound CE Web Server:  7Zvw7YVaDmV3
Current password for your Neo4j is:     QUPsoBXyTnBp
Go to the following directory:          /Users/median/neo4j-inst-caspian
Run the following:                      docker compose up -d OR docker-compose up -d

...

bloodhound_community ports -p caspian --local=true
neo4j-inst-caspian
==================
Bloodhound_Web_Port:      55161
NEO4J_Web_Port:           55154
NEO4J_Database_Port:      55155


bloodhound_community ports -p caspian --local=false
neo4j-inst-caspian
==================
ssh -L 8080:localhost:55161 -L 7474:localhost:55154 -L 7687:localhost:55155 username@SERVER_IP
```



