# Étape 1 : Builder l'application
FROM golang:1.23 as builder

# Définir le répertoire de travail dans le conteneur
WORKDIR /app

# Copier les fichiers nécessaires pour le build
COPY go.mod go.sum ./
RUN go mod download

# Copier le code source et le fichier .env
COPY . ./
COPY .env /app/.env

# Construire l'exécutable
RUN go build -o vote-api .

# Étape 2 : Image finale
FROM debian:bookworm-slim

# Installer les dépendances nécessaires
RUN apt-get update && apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Définir le répertoire de travail
WORKDIR /app

# Copier l'exécutable depuis l'étape de build
COPY --from=builder /app/vote-api .

# Copier le fichier .env dans l'image finale
COPY --from=builder /app/.env .

# Exposer le port utilisé par l'application
EXPOSE 8080

# Commande par défaut pour démarrer l'application
CMD ["./vote-api"]
