## **Docker**

### Prerequisite
- Docker Windows/Linux

### **Start Container**
To start the container use: `docker compose up --build`

#### **Seed Database**

To run the seeder use: 
```sh
docker-compose --profile seeder up seeder # One time only
docker-compose up --build seeder
```


