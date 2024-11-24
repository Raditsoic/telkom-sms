#!/bin/bash

CONTAINER_NAME="telkom-sms-psql-1"
DB_USER="your_postgres_username"
DB_NAME="your_database_name"
EXPORT_PATH="./"
OUTPUT_FILE="$EXPORT_PATH/items.csv"

QUERY="
COPY (
    SELECT
        i.id AS item_id,
        c.name AS category_name,
        i.name AS item_name,
        i.quantity,
        i.shelf
    FROM items i
    LEFT JOIN categories c ON i.category_id = c.id
) TO STDOUT WITH CSV HEADER;
"

docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -c "$QUERY" > $OUTPUT_FILE

if [[ $? -eq 0 ]]; then
    echo "Export complete. File saved to $OUTPUT_FILE"
else
    echo "Export failed."
fi
