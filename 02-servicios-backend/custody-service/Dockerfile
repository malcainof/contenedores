# --------------------------------------------------------------------------- #
# Parte 1: build del binario
# --------------------------------------------------------------------------- #
ARG GO_VER=20

# La compilación la realizaremos usando una imagen de Golang estándar
FROM golang:1.${GO_VER} as builder

# Nos posicionamos en un directorio de trabajo
WORKDIR /workspace

# Copiamos los archivos go.mod y go.sum para poder importar las librerías
COPY go.mod go.sum ./

# Ejecutamos el comando que permite importar las librerías
RUN go mod download

# Copiamos todos los otros archivos al directorio de trabajo
COPY . .

# Ejecutamos la compilación del microservicio. Esto generará un binario nativo y optimizado, llamado "server", en el directorio de trabajo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o server server.go

# --------------------------------------------------------------------------- #
# Parte 2: entorno de ejecución
# --------------------------------------------------------------------------- #
# Se exporta siempre el último "FROM", los otros contenedores se consideran de paso
# La exportación estará basada en Distroless, una distribución de Linux que no tiene intérprete de comandos (shell), lo que reduce los posibles vectores de ataque
FROM gcr.io/distroless/static:nonroot as output

# Nos posicionamos en un directorio de trabajo
WORKDIR /go/src/app

# Copiamos desde el paso anterior (Parte 1: build del binario) el binario nativo al actual directorio de trabajo.
# La instrucción "--from=builder" indica que el origen del archivo (/workspace/server) proviene de la etapa marcada con la instrucción "as builder"
COPY --from=builder /workspace/server ./server

# Nos cambiamos a un usuario sin privilegios. En este caso, la imagen de base Distroless tiene un usuario llamado "nonroot".
# Es mala práctica crear imágenes donde se requiera el uso del usuario ROOT
USER nonroot

# Definimos el comando de arranque de nuestro microservicio
CMD [ "/go/src/app/server" ]

# --------------------------------------------------------------------------- #