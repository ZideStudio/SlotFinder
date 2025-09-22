# PostgreSQL Dockerfile Usage

This directory contains a custom `Dockerfile` for running a PostgreSQL database for the SlotFinder project.

## How to Use

1. **Build the Docker image**

   From the `db/` directory:

   ```bash
   docker build -t slotfinder-postgres .
   ```

2. **Run the PostgreSQL container**

   You must provide all required PostgreSQL environment variables in the `docker run` command. These variables **must match** those defined in your `back/.env` file to ensure your backend can connect to the database.

   Example:

   ```bash
   docker run -d \
     --name slotfinder-db \
     -e POSTGRES_USER=slotfinder \
     -e POSTGRES_PASSWORD=slotfinder \
     -e POSTGRES_DB=slotfinder \
     -p 5432:5432 \
     slotfinder-postgres
   ```

   Replace the values with those from your `back/.env` if they are different.


3. **Stop and remove the container**

   ```bash
   docker stop slotfinder-db
   docker rm slotfinder-db
   ```

