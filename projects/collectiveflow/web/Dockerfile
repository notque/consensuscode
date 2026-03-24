# CollectiveFlow Web Interface — Simple Flask Container
# No special knowledge needed: just `docker build -t collectiveflow-web .`

FROM python:3.12-slim

WORKDIR /app

# Install dependencies first for better layer caching
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create data directory
RUN mkdir -p /app/data/proposals

EXPOSE 5000

CMD ["python3", "app.py"]
