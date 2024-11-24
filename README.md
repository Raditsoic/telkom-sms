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

## **How to export database**

### **Linux**

```sh
cd export
chmod +x export_item.sh export_transaction.sh
./export_item.sh && ./export_transaction.sh
```

### **Windows**
```ps1
cd export
\.export_item.ps1
\.export_transaction.ps1
```


