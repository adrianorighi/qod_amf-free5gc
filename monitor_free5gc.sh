#!/bin/bash

# =============================================================
# CONFIGURAÇÕES
# =============================================================

# Arquivo de log (CSV)
LOG_FILE="resource_monitoring_free5gc.csv"

# Intervalo de monitoramento (em segundos)
MONITOR_INTERVAL=5

# Filtros para containers:
# Procurar por containers que contenham 'free5gc' OU 'ueransim' no nome
# Para monitorar TODOS os containers, substitua a linha abaixo por CONTAINER_IDS=$(docker ps -q)
CONTAINER_IDS=$(docker ps -q)

# Formato de saída do docker stats (usando ; como separador)
# NOTA: {{.MemUsage}} retorna a string bruta (ex: 100MiB / 2.00GiB), que é ideal
FORMAT_STRING='{{.Name}};{{.CPUPerc}};{{.MemUsage}}'

# =============================================================
# INICIALIZAÇÃO
# =============================================================

# 1. Verifica se os containers foram encontrados
if [ -z "$CONTAINER_IDS" ]; then
    echo "ERRO: Nenhum container com 'free5gc', 'ueransim' ou 'upf' foi encontrado em execução."
    echo "Verifique se os serviços estão rodando e se os nomes de filtro estão corretos."
    exit 1
fi

# 2. Cria o arquivo de log e adiciona o cabeçalho CSV
if [ ! -f "$LOG_FILE" ]; then
    echo "Criando arquivo de log: $LOG_FILE"
    echo "Timestamp;ContainerName;CPU_Percent;Memory_Usage_and_Limit" > "$LOG_FILE"
fi

echo "Iniciando monitoramento dos containers: $CONTAINER_IDS"
echo "Pressione [Ctrl+C] para parar o monitoramento."
echo "--------------------------------------------------------"

# =============================================================
# LOOP DE MONITORAMENTO
# =============================================================

while true; do
    # 1. Captura o Timestamp
    TIMESTAMP=$(date +%Y-%m-%d\ %H:%M:%S)
    
    # 2. Captura as estatísticas (sem stream)
    STATS=$(docker stats --no-stream --format "$FORMAT_STRING" $CONTAINER_IDS 2>/dev/null)
    
    # Verifica se há dados para evitar erros e loops vazios
    if [ -n "$STATS" ]; then
        # 3. Processa cada linha e adiciona ao arquivo de log
        echo "$STATS" | while IFS= read -r LINE; do
            # Adiciona o Timestamp e remove qualquer caractere de retorno de carro (\r)
            echo "${TIMESTAMP};${LINE//$'\r'/}" >> "$LOG_FILE"
        done
    else
        echo "[$(date +%H:%M:%S)] Aviso: Nenhum dado de estatística recebido. Os containers pararam?"
    fi

    # 4. Aguarda o intervalo definido
    sleep "$MONITOR_INTERVAL"
done
