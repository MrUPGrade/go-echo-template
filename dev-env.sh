HOST_IP=$(hostname -I | cut -d \  -f 1)

# DB
export ET_DB_HOST=${HOST_IP}
export ET_DB_PORT="21001"
export ET_DB_USER="user"
export ET_DB_PASS="pass"
export ET_DB_NAME="db"
export ET_DB_DIALECT="postgres"

# Prometheus
export ET_PROM_PGW_PORT=21002

env | grep ET_ > .dev.env
