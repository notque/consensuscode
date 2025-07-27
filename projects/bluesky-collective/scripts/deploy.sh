#!/bin/bash

# Collective Deployment Script
# Autonomous deployment with collective safeguards
# No gatekeepers, maximum transparency

set -euo pipefail

# Collective configuration
COLLECTIVE_NAME="consensus-code-collective"
PROJECT_NAME="bluesky-collective"
REGISTRY="ghcr.io"
NAMESPACE="collective"

# Colors for collective feedback
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Collective logging function
log() {
    echo -e "${BLUE}[COLLECTIVE]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[COLLECTIVE WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[COLLECTIVE ERROR]${NC} $1"
    exit 1
}

success() {
    echo -e "${GREEN}[COLLECTIVE SUCCESS]${NC} $1"
}

# Function to check collective consensus
check_collective_consensus() {
    local environment=$1
    
    log "Checking collective consensus for $environment deployment..."
    
    if [[ "$environment" == "production" ]]; then
        # Production requires consensus file or special flag
        if [[ ! -f ".collective-consensus" ]] && [[ "${COLLECTIVE_CONSENSUS:-}" != "true" ]]; then
            error "Production deployment requires collective consensus. Create '.collective-consensus' file or set COLLECTIVE_CONSENSUS=true"
        fi
        
        # Check for minimum reviewers on latest commit
        if command -v gh &> /dev/null; then
            local commit_sha=$(git rev-parse HEAD)
            local pr_number=$(gh pr list --state merged --search "$commit_sha" --json number --jq '.[0].number // empty')
            
            if [[ -n "$pr_number" ]]; then
                local reviews=$(gh pr view "$pr_number" --json reviews --jq '.reviews | length')
                if [[ "$reviews" -lt 2 ]]; then
                    warn "Production deployment with fewer than 2 reviews. Collective consensus override required."
                    if [[ "${COLLECTIVE_OVERRIDE:-}" != "true" ]]; then
                        error "Set COLLECTIVE_OVERRIDE=true to proceed without full consensus"
                    fi
                fi
            fi
        fi
    fi
}

# Function to build collective container
build_container() {
    local version=$1
    local image_tag="${REGISTRY}/${COLLECTIVE_NAME}/${PROJECT_NAME}:${version}"
    
    log "Building collective container: $image_tag"
    
    # Build with collective metadata
    docker build \
        --build-arg VERSION="$version" \
        --build-arg BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        --build-arg COLLECTIVE="$COLLECTIVE_NAME" \
        --build-arg GOVERNANCE="horizontal" \
        --tag "$image_tag" \
        .
    
    # Tag as latest for staging
    if [[ "$version" == "staging" ]]; then
        docker tag "$image_tag" "${REGISTRY}/${COLLECTIVE_NAME}/${PROJECT_NAME}:latest"
    fi
    
    success "Container built successfully"
    echo "$image_tag"
}

# Function to push container to registry
push_container() {
    local image_tag=$1
    
    log "Pushing container to collective registry..."
    
    # Login to registry if not already authenticated
    if ! docker info | grep -q "Username:"; then
        log "Authenticating with container registry..."
        echo "$GITHUB_TOKEN" | docker login "$REGISTRY" -u "$GITHUB_ACTOR" --password-stdin
    fi
    
    docker push "$image_tag"
    
    # Push latest tag for staging
    if [[ "$image_tag" == *":staging" ]]; then
        docker push "${REGISTRY}/${COLLECTIVE_NAME}/${PROJECT_NAME}:latest"
    fi
    
    success "Container pushed to registry"
}

# Function to deploy to Kubernetes
deploy_kubernetes() {
    local environment=$1
    local image_tag=$2
    
    log "Deploying to collective Kubernetes cluster ($environment)..."
    
    # Check if kubectl is configured
    if ! kubectl cluster-info &> /dev/null; then
        error "kubectl not configured. Please configure access to collective cluster."
    fi
    
    # Create namespace if it doesn't exist
    kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
    
    # Update image in deployment
    local deployment_file="infrastructure/kubernetes/bluesky-collective.yaml"
    
    # Create a temporary file with updated image
    local temp_file=$(mktemp)
    cp "$deployment_file" "$temp_file"
    
    # Replace image tag
    sed -i.bak "s|image: ghcr.io/.*|image: $image_tag|g" "$temp_file"
    
    # Apply deployment
    kubectl apply -f "$temp_file" -n "$NAMESPACE"
    
    # Wait for rollout to complete
    kubectl rollout status deployment/bluesky-collective -n "$NAMESPACE" --timeout=300s
    
    # Cleanup
    rm "$temp_file" "$temp_file.bak"
    
    success "Deployment completed successfully"
}

# Function to deploy to AWS ECS
deploy_ecs() {
    local environment=$1
    local image_tag=$2
    
    log "Deploying to collective ECS cluster ($environment)..."
    
    # Check if AWS CLI is configured
    if ! aws sts get-caller-identity &> /dev/null; then
        error "AWS CLI not configured. Please configure access to collective AWS account."
    fi
    
    local cluster_name="${COLLECTIVE_NAME}-cluster"
    local service_name="${COLLECTIVE_NAME}-service"
    local task_family="${COLLECTIVE_NAME}-task"
    
    # Update task definition with new image
    local task_definition=$(aws ecs describe-task-definition --task-definition "$task_family" --query 'taskDefinition')
    
    # Create new task definition with updated image
    local new_task_definition=$(echo "$task_definition" | jq --arg image "$image_tag" '
        .containerDefinitions[0].image = $image |
        del(.taskDefinitionArn, .revision, .status, .requiresAttributes, .placementConstraints, .compatibilities, .registeredAt, .registeredBy)
    ')
    
    # Register new task definition
    local new_task_arn=$(echo "$new_task_definition" | aws ecs register-task-definition --cli-input-json file:///dev/stdin --query 'taskDefinition.taskDefinitionArn' --output text)
    
    # Update service with new task definition
    aws ecs update-service \
        --cluster "$cluster_name" \
        --service "$service_name" \
        --task-definition "$new_task_arn" \
        --force-new-deployment
    
    # Wait for deployment to complete
    aws ecs wait services-stable \
        --cluster "$cluster_name" \
        --services "$service_name"
    
    success "ECS deployment completed successfully"
}

# Function to run health checks
health_check() {
    local environment=$1
    local max_attempts=30
    local attempt=0
    
    log "Running collective health checks..."
    
    # Determine health check URL based on environment
    local health_url
    case "$environment" in
        "staging")
            health_url="https://staging.bluesky.consensuscode.dev/health"
            ;;
        "production")
            health_url="https://bluesky.consensuscode.dev/health"
            ;;
        *)
            warn "Unknown environment: $environment. Skipping health check."
            return 0
            ;;
    esac
    
    while [[ $attempt -lt $max_attempts ]]; do
        if curl -sf "$health_url" > /dev/null 2>&1; then
            success "Health check passed"
            return 0
        fi
        
        ((attempt++))
        log "Health check attempt $attempt/$max_attempts failed, retrying in 10 seconds..."
        sleep 10
    done
    
    error "Health check failed after $max_attempts attempts"
}

# Function to notify collective
notify_collective() {
    local environment=$1
    local status=$2
    local image_tag=$3
    
    log "Notifying collective of deployment status..."
    
    # Create notification message
    local message="🚀 Collective Deployment Update
Environment: $environment
Status: $status
Image: $image_tag
Time: $(date -u +%Y-%m-%dT%H:%M:%SZ)
Deployer: ${GITHUB_ACTOR:-$(whoami)}
Commit: $(git rev-parse --short HEAD)"

    # Send to Slack/Discord/Matrix if webhook configured
    if [[ -n "${COLLECTIVE_WEBHOOK_URL:-}" ]]; then
        curl -X POST "$COLLECTIVE_WEBHOOK_URL" \
            -H "Content-Type: application/json" \
            -d "{\"text\": \"$message\"}" \
            2>/dev/null || warn "Failed to send webhook notification"
    fi
    
    # Log to deployment log
    echo "$message" >> "deployments.log"
    
    success "Collective notification sent"
}

# Main deployment function
main() {
    local command=${1:-"help"}
    local environment=${2:-"staging"}
    local version=${3:-$(git rev-parse --short HEAD)}
    
    log "Starting collective deployment process..."
    log "Command: $command, Environment: $environment, Version: $version"
    
    case "$command" in
        "build")
            image_tag=$(build_container "$version")
            echo "Built image: $image_tag"
            ;;
        "push")
            image_tag="${REGISTRY}/${COLLECTIVE_NAME}/${PROJECT_NAME}:${version}"
            push_container "$image_tag"
            ;;
        "deploy")
            check_collective_consensus "$environment"
            
            image_tag="${REGISTRY}/${COLLECTIVE_NAME}/${PROJECT_NAME}:${version}"
            
            # Determine deployment target
            if kubectl cluster-info &> /dev/null; then
                deploy_kubernetes "$environment" "$image_tag"
            elif aws sts get-caller-identity &> /dev/null; then
                deploy_ecs "$environment" "$image_tag"
            else
                error "No deployment target configured (Kubernetes or AWS ECS)"
            fi
            
            health_check "$environment"
            notify_collective "$environment" "SUCCESS" "$image_tag"
            ;;
        "full")
            check_collective_consensus "$environment"
            
            # Full deployment pipeline
            image_tag=$(build_container "$version")
            push_container "$image_tag"
            
            if kubectl cluster-info &> /dev/null; then
                deploy_kubernetes "$environment" "$image_tag"
            elif aws sts get-caller-identity &> /dev/null; then
                deploy_ecs "$environment" "$image_tag"
            else
                error "No deployment target configured"
            fi
            
            health_check "$environment"
            notify_collective "$environment" "SUCCESS" "$image_tag"
            ;;
        "health")
            health_check "$environment"
            ;;
        "help"|*)
            echo "Collective Deployment Script"
            echo "Usage: $0 <command> [environment] [version]"
            echo ""
            echo "Commands:"
            echo "  build    - Build container image"
            echo "  push     - Push container to registry"
            echo "  deploy   - Deploy to cluster"
            echo "  full     - Build, push, and deploy"
            echo "  health   - Run health checks"
            echo "  help     - Show this help"
            echo ""
            echo "Environments: staging, production"
            echo "Version: defaults to git commit SHA"
            echo ""
            echo "Collective Environment Variables:"
            echo "  COLLECTIVE_CONSENSUS - Set to 'true' for production override"
            echo "  COLLECTIVE_OVERRIDE  - Set to 'true' to bypass review checks"
            echo "  COLLECTIVE_WEBHOOK_URL - Webhook for collective notifications"
            echo "  GITHUB_TOKEN         - Token for container registry"
            echo "  GITHUB_ACTOR         - GitHub username for attribution"
            ;;
    esac
}

# Run main function with all arguments
main "$@"