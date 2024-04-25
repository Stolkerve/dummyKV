# dummyKV

##### Es una base de datos clave-valor, inspirada en redis.

# Requisitos
Tener instalado la version 1.22.1 de [Golang](https://go.dev/)

# Comandos
| Comando | Argumentos | Descripcion                                                                             |
|---------|------------|-----------------------------------------------------------------------------------------|
| PING    |            | Sirve para ver el estado del servidor y sus latencias. El servidor Responde con un PONG |
| ECHO    | 1          | El servidor responde con el texto enviado                                               |
| GET     | 1          | El servidor responde el valor de clave enviada, si existe                               |
| SET     | 2          | El servidor almacena la clave y su valor asociado                                       |

# Servidor
El servidor escucha por defecto en la direccion: 0.0.0.0:8000.

#### Como ejecutarlo
1) compilar el programa con el comando
`go build -o dummyKV cmd/server/main.go`
2) Ver las opciones con el comando `dummyKV --help`
2) Ejecutarlo con el comando
`dummyKV`

# Cliente
El cliente ofrece interaccion mediente argumentos que se le pase al ejecutable o un modo interactivo
#### Como ejecutarlo
1) compilar el programa con el comando
`go build -o REPL cmd/REPL/main.go`
2) Ver las opciones con el comando `REPL --help`
2) Ejecutarlo en modo interactivo con el comando
`REPL`