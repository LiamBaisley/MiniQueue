Write-Host "Welcome to the MiniQ set up prompt" -ForegroundColor Green
Write-Host "| ------------------------------------------- |" -ForegroundColor Green
Write-Host "| ------------------------------------------- |" -ForegroundColor Green
$port= Read-Host -Prompt "Please provide a port nuumber on which the queue will be run"
docker build -t miniq ..
docker run -p ${port}:8080 -d miniq
Write-Host "| ------------------------------------------- |" -ForegroundColor Green
Write-Host "| ------------------------------------------- |" -ForegroundColor Green
Write-Host "Application now available on port: $port" -ForegroundColor Green
