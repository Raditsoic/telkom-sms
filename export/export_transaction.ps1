$ContainerName = "telkom-sms-psql-1"
$DbUser = "your_postgres_username"
$DbName = "your_database_name"
$ExportPath = ".\"
$OutputFile = "$ExportPath\transactions.csv"

$query = @"
COPY (
    SELECT 
        'LoanTransaction' AS transaction_type,
        lt.id,
        lt.uuid,
        lt.employee_name,
        lt.employee_department,
        lt.employee_position,
        c.name AS category_name,
        i.name AS item_name,
        lt.quantity,
        lt.status,
        lt.notes,
        lt.time,
        lt.item_id,
        lt.loan_time,
        lt.return_time,
        lt.completed_time,
        lt.returned_time,
        NULL::TEXT AS image
    FROM loan_transactions lt
    LEFT JOIN items i ON lt.item_id = i.id
    LEFT JOIN categories c ON i.category_id = c.id

    UNION ALL

    SELECT 
        'InquiryTransaction' AS transaction_type,
        it.id,
        it.uuid,
        it.employee_name,
        it.employee_department,
        it.employee_position,
        c.name AS category_name,
        i.name AS item_name,
        it.quantity,
        it.status,
        it.notes,
        it.time,
        it.item_id,
        NULL,
        NULL,
        it.completed_time,
        NULL,
        NULL::TEXT AS image
    FROM inquiry_transactions it
    LEFT JOIN items i ON it.item_id = i.id
    LEFT JOIN categories c ON i.category_id = c.id

    UNION ALL

    SELECT 
        'InsertionTransaction' AS transaction_type,
        int.id,
        int.uuid,
        int.employee_name,
        int.employee_department,
        int.employee_position,
        c.name AS category_name,
        i.name AS item_name,
        int.item_request_quantity AS quantity,
        int.status,
        int.notes,
        int.time,
        int.item_id,
        NULL,
        NULL,
        int.completed_time,
        NULL,
        ENCODE(int.image, 'base64') AS image
    FROM insertion_transactions int
    LEFT JOIN items i ON int.item_id = i.id
    LEFT JOIN categories c ON i.category_id = c.id
) TO STDOUT WITH CSV HEADER;
"@

docker exec -i $ContainerName psql -U $DbUser -d $DbName -c $query | Out-File -FilePath $OutputFile -Encoding utf8

if ($?) {
    Write-Host "Export complete. File saved to $OutputFile"
} else {
    Write-Host "Export failed."
}
