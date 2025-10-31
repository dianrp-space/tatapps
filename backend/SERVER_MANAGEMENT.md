# Backend Server Management

## 🚀 Quick Commands

### Start Backend
```bash
cd /var/www/html/tatapps/backend
./start.sh
```

### Stop Backend
```bash
cd /var/www/html/tatapps/backend
./stop.sh
```

### Restart Backend (Recommended)
```bash
cd /var/www/html/tatapps/backend
./restart.sh
```

---

## 📋 Manual Commands

### Start Backend Manually
```bash
cd /var/www/html/tatapps/backend
nohup go run cmd/api/main.go > /tmp/backend.log 2>&1 &
```

### Stop Backend Manually
```bash
# Kill by port
lsof -ti:8080 | xargs kill -9

# Or kill by process name
pkill -f "go run cmd/api/main.go"
```

### Check Backend Status
```bash
# Check if running
lsof -ti:8080

# Test API
curl http://localhost:8080/health

# View logs
tail -f /tmp/backend.log

# View last 50 lines
tail -50 /tmp/backend.log
```

---

## 🔍 Troubleshooting

### Port Already in Use
```bash
# Find what's using port 8080
lsof -i:8080

# Kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

### Check Backend Process
```bash
# Find backend process
ps aux | grep "go run cmd/api/main.go"

# Kill specific PID
kill -9 <PID>
```

### Database Connection Error
```bash
# Test database connection
PGPASSWORD='tatapps_password' psql -h localhost -U tatapps_user -d tatapps -c "SELECT 1"

# Check .env file
cat /var/www/html/tatapps/backend/.env
```

### View Real-time Logs
```bash
# Follow logs in real-time
tail -f /tmp/backend.log

# Search logs for errors
grep -i error /tmp/backend.log
```

---

## 🌐 API Endpoints

- **Health Check**: http://localhost:8080/health
- **Login**: http://localhost:8080/api/v1/auth/login
- **API Docs**: See API_DOCUMENTATION.md

---

## ⚙️ Configuration

Database configuration in `.env`:
```properties
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tatapps
DB_USER=tatapps_user
DB_PASSWORD=tatapps_password
```

Server will run on: **http://localhost:8080**
