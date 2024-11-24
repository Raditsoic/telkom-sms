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

## **How to see database**
```sh
docker exec -i telkom-sms-psql-1 psql -U your_postgres_username -d your_database_name
```

```sql
\dt # Kalo mau liat tabel apa aja yg ada
SELECT * FROM nama_tabel; # Kalo mau liat isi tabel, nama tabel diganti ke nama tabel yg kaya di \dt
```


