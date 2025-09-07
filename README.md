# Educational platform 

This project demonstrates full cycle of an educational platform development with a focus on microservice architecture, database design, deployment and management. Previously focused on database design, but now it's much more about backend written in GoLang that emulates the platform's work.



## Used technologies and tools

- **GoLang** — backend development
- **GoLang & gofakeit** — data generation (seeding)
- **Docker & Docker Compose** — containerization and orchestration
- **PostgreSQL** — relational database management system
- **Redis** – in-memory data structure storage for caching
- **Patroni & etcd** — automatically replicates DB for high availability and clustering
- **HAProxy** — load balancing
- **Flyway** — database migration management
- **Prometheus & Grafana** — monitoring and visualization


## Project deploy
In order to deploy project you need to manually build patroni inside `/db` folder.
It's needed because existing patroni images don't work as expected or don't work at all.
So it was decided to use manually built one.

**Step 1.**  
Create and set `.env` file as it shown in `.env.example` 
or explicitly indicate your desire to use file with another name.


**Step 2.**  
[Clone and build patroni](https://github.com/patroni/patroni) inside `/db` folder using instruction.  

**Step 3.**   
After setting `.env` file run `docker-compose build`  
Or specify another env file by adding `--env-file env_filename` 

**Step 4.**  
During the first setup it's recommended to use only `initial` profile.  
After successful build run `docker-compose --profile initial up`.  
Add `-d` at the end if you want to run in background.

The initial setup is designed to configure database: apply migrations and seed DB with random data. 
Seeding is heavy computational task so that's why it's not recommended
to use `benchmark` profile with `initail` at the same time 

**Step 5.**  
After configuring and seeding database you can stop and run project again with `benchmark` profile using:  
`docker-compose --profile benchmark up`  
Now you'll be able to see dashboards in grafana with load testing

## Monitoring & Grafana dashboards
**Monitoring and services available at:**
1. Prometheus: http://localhost:9090  
2. Grafana: http://localhost:3000 (login/password: admin/admin)   
3. Go Backend: http://localhost:8080
4. Redis: http://localhost:6379
5. HAProxy: http://localhost:5001

## Project lifecycle
After starting all services the following sequence of actions are execute:  
1. etcd cluster up
2. Patroni replicas starting and choose master
3. HAProxy balancing requests and determining current master and available replicas
4. Flyway apply migrations up to MIGRATION_VERSION set in .env file
5. GoLang + gofakeit generates testing data and filling DB with it
6. Backend service starts
7. Prometheus and Grafana begin monitoring routine
8. Load testing client making different requests to service

## Repo Structure

Backend service was created according to [Go Clean Template](https://github.com/evrone/go-clean-template/) as one of the
best microservice architecture examples.

```
backend/                            # GoLang backend microservice 
├── cmd/                            # Entry point
├── config/                         # Services configs
├── docs/
├── internal/                       # Bussiness logic layers
│   ├── app/                        
│   ├── controller/                 
│   ├── entity/
│   ├── repo/
│   └── usecase/
├── pkg/                            # Tools (DI)
└── Dockerfile
client/                             # Load testing client
├── Dockerfile
└── main.go
db/                                 # DB and provisioning
├── config/                         # Docker files, Grafana, Prometheus and Redis configs
├── init/
├── migrations/                     # DB migrations
├── patroni/                        # Manually built patroni image (check Step 2)
├── scripts/
│   ├── seeding/                    # GoLangg seeding service
│   └── ...                         # Some currently unused scripts
└── 
docs/                               # Project docs
├── DB_ERD.md                       # PlantUML Database ER Diagram
└── DB_FUNCTIONAL_REQUIREMENTS.md   # Russian DB Reqs
.env.example
docker-compose.yml                  # Main docker-compose file with deploy
```

## Docs
1. [Backend Microservice Documentation](backend/docs/SERVICE.md) 
2. [DB ERD](docs/DB_ERD.md)
3. [DB Functional Requirements](docs/DB_FUNCTIONAL_REQUIREMENTS.md) (RU version)
