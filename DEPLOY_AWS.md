# Deploying to AWS Lightsail (Fastest & Cheapest)

For a simple Go + SQLite application, **AWS Lightsail** is the best balance of speed, cost, and low overhead. It avoids the complexity of EC2/EKS and the storage persistence issues of Lambda.

**Cost**: ~$5/month (First 3 months free on some plans).

## Step 1: Create Instance
1. Log into AWS Console -> **Lightsail**.
2. Click **Create instance**.
3. Platform: **Linux/Unix**.
4. Blueprint: **OS Only** -> **Ubuntu 22.04 LTS**.
5. Choose instance plan (The cheapest $3.50 or $5 plan is plenty).
6. Name it `fieldnotes-api`.
7. Click **Create instance**.

## Step 2: Prepare the Server
Once the instance is "Running", click the orange terminal icon (Connecting using SSH) or use your own terminal.

Run these commands on the server to install Docker:
```bash
# Update and install Docker
sudo apt-get update
sudo apt-get install -y docker.io docker-compose
sudo usermod -aG docker ubuntu
```
*Close the terminal window and reopen it to apply the group change.*

## Step 3: Deploy the App
You can deploy by copying files directly. (In a real CI/CD you'd push the image to a registry, but this is faster for assessment).

### Option A: Manual SCP (Fastest for today)
From your local machine (where your code is):

1. **Build the image locally and save it** (optional, or just copy source):
   ```bash
   # Create a directory on the server
   ssh ubuntu@<YOUR_INSTANCE_IP> "mkdir -p fieldnotes"
   
   # Copy files (including new Caddyfile)
   scp Dockerfile docker-compose.yml Caddyfile go.mod go.sum *.go ubuntu@<YOUR_INSTANCE_IP>:~/fieldnotes/
   ```

2. **Run on Server**:
   ```bash
   ssh ubuntu@<YOUR_INSTANCE_IP>
   cd fieldnotes
   
   # Build and Run
   docker-compose up -d --build
   ```

## Step 4: Open Firewall
1. In Lightsail Console, go to **Networking** tab.
2. Under "IPv4 Firewall":
   - **Keep/Add** Port **80** (HTTP)
   - **Add** Port **443** (HTTPS)
   - **Remove/Ignore** Port 8080 (It is now internal only)
3. Access your API at: `http://<YOUR_INSTANCE_IP>/v1/notes` (No port needed!)

## HTTPS (SSL)
If you point a domain (e.g., `api.myapp.com`) to your instance IP:
1. Edit `Caddyfile` on the server: `nano Caddyfile`
2. Change `:80` to `api.myapp.com`.
3. Restart: `docker-compose restart proxy`.
4. Caddy automatically gets a standard encryption certificate for you!

## Maintenance
- **Logs**: `docker-compose logs -f`
- **Data**: Your data is safe in `~/fieldnotes/data/notes.db` on the server disk.
