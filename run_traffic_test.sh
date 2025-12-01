#!/bin/bash

# =============================================================
# VARIÁVEIS DE CONFIGURAÇÃO
# =============================================================

UPF_SEARCH_TERM="upf"
UERANSIM_SEARCH_TERM="ueransim"
IPERF_SERVER_PORT="5201" # Porta padrão do iperf3

# Função para encontrar o nome do container em execução
find_container() {
    local search_term=$1
    local container_name=$(docker ps --filter name="$search_term" --format "{{.Names}}" | head -n 1)

    if [ -z "$container_name" ]; then
        echo "ERRO: Nenhum container em execução com o termo '$search_term' encontrado."
        exit 1
    fi
    echo "$container_name"
}

# Encontrar containers
echo "--- 1. Encontrando Containers ---"
UPF_CONTAINER=$(find_container "$UPF_SEARCH_TERM")
UERANSIM_CONTAINER=$(find_container "$UERANSIM_SEARCH_TERM")

echo "UPF Container encontrado: $UPF_CONTAINER"
echo "UERANSIM Container encontrado: $UERANSIM_CONTAINER"
echo "---------------------------------"

# =============================================================
# 2. INICIAR SERVIDOR IPERF3 (UPF)
# =============================================================
echo "--- 2. Iniciando iperf3 Server no UPF ($UPF_CONTAINER) ---"
#docker exec -d "$UPF_CONTAINER" iperf3 -s -B 0.0.0.0 -D

#if [ $? -ne 0 ]; then
#    echo "AVISO: Falha ao iniciar iperf3 no UPF. Verifique se o iperf3 está instalado no container UPF."
#fi
#echo "Servidor iperf3 iniciado em segundo plano na porta $IPERF_SERVER_PORT."
echo "---------------------------------"

# =============================================================
# 3. INICIAR UES (UERANSIM) EM PARALELO
# =============================================================
echo "--- 3. Iniciando UEs em Paralelo no UERANSIM ($UERANSIM_CONTAINER) ---"
docker exec -d "$UERANSIM_CONTAINER" sh -c \
  './nr-ue -c config/uecfg.yaml &'

if [ $? -ne 0 ]; then
    echo "ERRO: Falha ao iniciar simuladores de UE. Verifique se os arquivos de configuração existem e o executável nr-ue está acessível no container."
    exit 1
fi
echo "Ambos os UEs foram iniciados e estão tentando se registrar na rede 5G."
echo "Aguardando 10 segundos para o estabelecimento da sessão PDU..."
sleep 10
echo "---------------------------------"


# =============================================================
# 4. EXECUTAR SCRIPTS DE TRÁFEGO
# =============================================================
echo "--- 4. Executando Testes de Tráfego (test_ueransim.sh uesimtun0) ---"

# Executa o primeiro script
if [ -f "test_ueransim.sh" ]; then
    echo "Executando test_ueransim.sh uesimtun0..."
    ./test_ueransim.sh uesimtun0
else
    echo "AVISO: Arquivo test_ueransim.sh não encontrado no diretório atual."
fi

# Executa o segundo script
#if [ -f "test_ueransim.sh" ]; then
#    echo "Executando test_ueransim.sh uesimtun1..."
#    ./test_ueransim.sh uesimtun1
#else
#    echo "AVISO: Arquivo test_ueransim.sh não encontrado no diretório atual."
#fi

echo "---------------------------------"

echo "--- 5. Fim do Teste ---"
