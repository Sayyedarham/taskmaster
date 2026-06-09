# Fix dirty database migration
# Run this in PowerShell from the taskmaster-starter directory

Write-Host "Stopping migrate container..."
docker-compose stop migrate

Write-Host "Resetting migration state..."
docker-compose exec -T postgres psql -U taskmaster -d taskmaster -c "UPDATE schema_migrations SET version = 1, dirty = false;"

Write-Host "Re-running migrations..."
docker-compose up migrate

Write-Host "Done!"
