<h1>cPanel Finder</h1>

<p>A script to identify cPanel on websites</p>

<hr>

<h2>Instalação</h2>

    git clone https://github.com/k45t0/cPanel-Finder.git

    cd cPanel-Finder

    go build cpanelfinder.go

    cp cpanelfinder /usr/local/bin

<hr>
<h2>Modo de Uso</h2>

    cpaneldiscovery -l domains.txt -t 20 -o cpnalvalid.txt

<hr>
<h2>Help</h2>     

    -d string
      	Domínio único para verificar (sem protocolo, exemplo: example.com)
    -l string
      	Arquivo contendo a lista de URLs/domínios
    -o string
      	Nome do arquivo de saída para domínios válidos (default "cpnalvalid.txt")
    -p int
      	Porta para verificar o cPanel (padrão: 2083) (default 2083)
    -t int
      	Número de threads/concorrência (default 10)
