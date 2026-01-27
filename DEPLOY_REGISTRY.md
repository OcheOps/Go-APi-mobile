# Deploying via Docker Registry (Hub, ECR, etc.)

This method is cleaner. You build locally, push to a registry, and just run the containers on AWS.

## 1. Local: Build & Push
Replace `your-registry` with your actual Docker Hub username or Registry URL.

```bash
# Set your image name
export API_IMAGE=your-registry/fieldnotes-api:v1

# Build
docker build -t $API_IMAGE .

# Push
docker push $API_IMAGE
```

## 2. Server (Lightsail): Deploy
1. **Copy only configuration files** (No source code needed!):
   ```bash
   scp docker-compose.yml Caddyfile ubuntu@<YOUR_INSTANCE_IP>:~/fieldnotes/
   ```

2. **Run on Server**:
   ```bash
   ssh ubuntu@<YOUR_INSTANCE_IP>
   cd fieldnotes
   
   # Set the same image variable
   export API_IMAGE=your-registry/fieldnotes-api:v1
   
   # Pull and Start (Caddy will handle HTTPS automatically for your domain)
   docker-compose up -d
   ```

## Note on Domain
Ensure `interviewapi.czettapay.com` points to your Lightsail IP address. Caddy will fail to get an SSL certificate if the DNS is not pointing to the server.
