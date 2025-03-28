FROM golang:1.23-alpine 

WORKDIR /app
RUN apk add --no-cache gcc musl-dev
# Primero copia los archivos de módulo
COPY go.mod go.sum ./
RUN go mod download

# Luego copia el resto del código
COPY . .

# Compila la aplicación
RUN go build -o backend .

EXPOSE 8080

CMD ["./backend"]