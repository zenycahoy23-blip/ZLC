        stage('Deploy') {
            steps {
                echo 'Deploying fresh services'
                sh '''
                    set -e
                    
                    echo "Step 1: Stopping old containers..."
                    docker-compose down --remove-orphans || true
                    
                    echo "Step 2: Removing lingering containers by name..."
                    docker rm -f postgres-uno backend-uno frontend-uno prometheus-uno grafana-uno 2>/dev/null || true
                    
                    echo "Step 3: Waiting for cleanup..."
                    sleep 3
                    
                    echo "Step 4: Building fresh images (no cache)..."
                    docker-compose build --no-cache --force-rm
                    
                    echo "Step 5: Creating network..."
                    docker network create uno-reverse 2>/dev/null || true
                    
                    echo "Step 6: Starting containers (force recreate)..."
                    docker-compose up -d --force-recreate --remove-orphans
                    
                    echo "Step 7: Waiting for containers to stabilize..."
                    sleep 8
                    
                    echo "Step 8: Verifying deployment..."
                    docker-compose ps
                '''
            }
        }
