# Collective Website — Simple Flask Container
# No special knowledge needed: just `docker build -t collective-website .`

FROM python:3.12-slim

WORKDIR /app

# Install dependencies first for better layer caching
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

EXPOSE 5001

CMD ["python3", "app.py"]
