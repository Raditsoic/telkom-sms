# Variables
$ContainerName = "telkom-sms-psql-1"
$DbUser = "your_postgres_username"
$DbName = "your_database_name"
$ExportPath = ".\"
$OutputFile = "$ExportPath\items.csv"

# SQL Query
$query = @"
COPY (
    SELECT
        i.id AS item_id,
        i.name AS item_name,
        c.name AS category_name,
        i.quantity,
        i.shelf
    FROM items i
    LEFT JOIN categories c ON i.category_id = c.id
) TO STDOUT WITH CSV HEADER;
"@

# Execute the Query in Docker
docker exec -i $ContainerName psql -U $DbUser -d $DbName -c $query | Out-File -FilePath $OutputFile -Encoding utf8

# Check if the Export Succeeded
if ($?) {
    Write-Host "Export complete. File saved to $OutputFile"
} else {
    Write-Host "Export failed."
}
