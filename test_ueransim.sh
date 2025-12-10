#!/bin/bash

# =================CONFIGURAÇÕES=================
# Nome do container do UE (conforme seu docker-compose)
UE_CONTAINER="ueransim" 

# Nome do container do UPF
UPF_CONTAINER="upf"

# ----------------------------------------------------------------------
#VERIFICAÇÃO DE PARÂMETRO
# ----------------------------------------------------------------------
if [ -z "$1" ]; then
    echo "---------------------------------------------------------------"
    echo "ERRO: O nome da interface de rede do túnel deve ser fornecido como parâmetro."
    echo "Uso: ./test-ueransim-0.sh <NOME_DA_INTERFACE_DO_TÚNEL>"
    echo "---------------------------------------------------------------"
    exit 1
fi

TUNNEL_INTERFACE=$1
echo "Interface de Túnel Definida: $TUNNEL_INTERFACE"
IPERF3_PORT=$2
# ----------------------------------------------------------------------

# IP ou Hostname do Servidor iperf3 de destino
UPF_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$UPF_CONTAINER")
if [ -z "$UPF_IP" ]; then
    echo "ERRO: Não foi possível obter o endereço IP do container UPF."
    exit 1
fi

echo "UPF IP: '$UPF_IP'."

TARGET_SERVER="201.54.13.159" 
echo "TARGET SERVER: '$TARGET_SERVER'."

# Duração do teste em segundos
DURATION="60"
BANDWIDTH="250M"

LOG_FILE="traffic-${TUNNEL_INTERFACE}.json"
PING_LOG_FILE="ping-${TUNNEL_INTERFACE}.log"
# ===============================================

echo "--- Iniciando Teste de Tráfego e Latência 5G via UERANSIM ---"

# 1. Verificar se o container está rodando
if [ ! "$(docker ps -q -f name=$UE_CONTAINER)" ]; then
    echo "ERRO: O container '$UE_CONTAINER' não foi encontrado ou não está rodando."
    exit 1
fi

echo "[1/5] Container UE encontrado: $UE_CONTAINER"

# 2. Obter o IP da interface de túnel
UE_IP=$(docker exec $UE_CONTAINER ip -4 addr show $TUNNEL_INTERFACE | grep -oP '(?<=inet\s)\d+(\.\d+){3}')

if [ -z "$UE_IP" ]; then
    echo "ERRO: Não foi possível obter o IP da interface $TUNNEL_INTERFACE."
    echo "Verifique se o UE completou o registro na rede 5G com sucesso."
    exit 1
fi

echo "[2/5] IP do Túnel 5G detectado: $UE_IP (Interface: $TUNNEL_INTERFACE)"

# 3. Execução do PING em background
echo "[3/5] Iniciando PING contínuo em paralelo para $TARGET_SERVER (saída em $PING_LOG_FILE)..."
docker exec $UE_CONTAINER ping -I $TUNNEL_INTERFACE -i 0.2 -W 1 -c $((DURATION * 5)) "$TARGET_SERVER" > "$PING_LOG_FILE" &
PING_PID=$! 

echo "---------------------------------------------------------------"

# 4. Execução do iperf3 em foreground
echo "[4/5] Executando iperf3 (Cliente: $UE_IP -> Servidor: $TARGET_SERVER)..."
docker exec $UE_CONTAINER iperf3 \
    -c "$TARGET_SERVER" \
    -B "$UE_IP" \
    -t "$DURATION" \
    -b "$BANDWIDTH" \
    -p "$IPERF3_PORT" \
    -R -u -J \
    > "$LOG_FILE" 

# 5. Esperar o processo ping
echo "[5/5] Aguardando o término do processo PING..."
wait $PING_PID

echo "---------------------------------------------------------------"
echo "Teste de Throughput e Latência finalizado."

# Análise de Métricas
echo "Analizando métricas de Throughput (iperf3) no arquivo $LOG_FILE..."
THROUGHPUT=$(jq '.end.sum_received.bits_per_second' "$LOG_FILE")
JITTER=$(jq '.end.sum.jitter_ms' "$LOG_FILE")

echo "--- Resultados Consolidados ---"
echo "Throughput: $(echo "scale=2; $THROUGHPUT / 1000000" | bc) Mbps"
echo "Jitter: $JITTER ms"

echo "--- Resultados de Latência (Ping) no arquivo $PING_LOG_FILE ---"
PING_STATS=$(grep "rtt min/avg/max/mdev" "$PING_LOG_FILE")
echo "Latência: $PING_STATS"
