# Rapport de Déploiement avec Docker et Docker Compose

## Introduction
Dans ce projet, j'ai configuré et déployé une application multi-services composée de :
- Une API (`vote-api`),
- Une interface utilisateur (`web-client`),
- Une documentation statique (`docs`),
- Une base de données PostgreSQL (`db`).

J'ai utilisé Docker pour créer des images pour chaque service et Docker Compose pour orchestrer leur déploiement et faciliter leur exécution simultanée.

---

## Configuration avec Docker Compose

Voici le fichier `docker-compose.yml` utilisé pour orchestrer les services :

```yaml
version: '3.8'

services:
  api:
    build: ./vote-api
    ports:
      - "8080:8080"
    environment:
      - PG_URL=postgres://vote_user:password123@db:5432/vote_db?sslmode=disable
      - JSON_LOG=true
    depends_on:
      db:
        condition: service_healthy

  web-client: 
    build: 
      context: ./web-client
    ports:
      - "3000:3000"
    environment:
      - VOTE_API_BASE_URL=http://api:8080
    depends_on:
      - api

  docs:
    build:
      context: ./docs
    ports:
      - "4000:80"

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=vote_user
      - POSTGRES_PASSWORD=password123
      - POSTGRES_DB=vote_db
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U vote_user -d vote_db"]
      interval: 5s
      timeout: 5s
      retries: 5
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### Explications

1. **Service `api`** :
   - Construit à partir du Dockerfile dans `vote-api`.
   - Exposé sur le port `8080`.
   - Connecté à PostgreSQL via la variable d'environnement `PG_URL`.
   - Dépend de la base de données (`db`) et attend qu'elle soit prête avant de démarrer grâce à `depends_on`.

2. **Service `web-client`** :
   - Construit à partir du Dockerfile dans `web-client`.
   - Exposé sur le port `3000`.
   - Communique avec l'API via la variable `VOTE_API_BASE_URL`.

3. **Service `docs`** :
   - Construit à partir du Dockerfile dans `docs`.
   - Exposé sur le port `4000` pour servir la documentation statique.

4. **Service `db`** :
   - Utilise l'image officielle `postgres:15-alpine`.
   - Configuré avec des variables d'environnement pour l'utilisateur, le mot de passe et le nom de la base de données.
   - Utilise un volume nommé `postgres_data` pour persister les données.

---

## Dockerfile des Services

### Dockerfile pour `web-client`

```dockerfile
# Étape 1 : Build de l'application
FROM node:18 as builder
WORKDIR /app

# Copier les fichiers de configuration
COPY package.json yarn.lock ./
RUN yarn install

# Copier le reste du code
COPY . ./
RUN yarn build

# Étape 2 : Image pour l'exécution
FROM node:18-slim
WORKDIR /app

# Copier les fichiers nécessaires pour l'exécution
COPY --from=builder /app/package.json /app/yarn.lock ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

RUN yarn install --production

# Exposer le port
EXPOSE 3000

# Commande de démarrage
CMD ["yarn", "start"]
```

### Dockerfile pour `api`

```dockerfile
FROM golang:1.23-rc-alpine
WORKDIR /app
RUN apk add --no-cache postgresql-client
COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080

CMD ["./main"]
```

### Dockerfile pour `docs`

```dockerfile
# Étape 1 : Build de la documentation
FROM node:18 as builder
WORKDIR /app

# Copier les fichiers de configuration
COPY package.json yarn.lock ./

# Installer les dépendances nécessaires
RUN yarn install

# Copier le reste du code
COPY . ./

# Construire la documentation
RUN yarn build

# Étape 2 : Image pour servir la documentation
FROM nginx:alpine
WORKDIR /usr/share/nginx/html

# Copier la documentation générée dans l'image finale
COPY --from=builder /app/build .

# Exposer le port 80
EXPOSE 80

# Commande de démarrage
CMD ["nginx", "-g", "daemon off;"]
```

---

## Instructions pour exécuter le projet

1. Assurez-vous d'avoir Docker et Docker Compose installés sur votre machine.
2. Clonez le projet :
   ```bash
   git clone <repository_url>
   cd <repository_folder>
   ```
3. Lancez tous les services avec Docker Compose :
   ```bash
   docker-compose up --build
   ```
   Cette commande construit les images et démarre tous les services.

4. Une fois démarré, accédez aux services via les ports exposés :
   - **API** : [http://localhost:8080](http://localhost:8080)
   - **Interface utilisateur** : [http://localhost:3000](http://localhost:3000)
   - **Documentation** : [http://localhost:4000](http://localhost:4000)

5. Pour arrêter les services :
   ```bash
   docker-compose down
   ```

---

## Conclusion
Ce projet est maintenant entièrement conteneurisé et peut être exécuté facilement avec Docker Compose. Chaque service est isolé, mais ils interagissent grâce à la configuration centralisée. Les fichiers Docker et Docker Compose simplifient la gestion des dépendances et assurent la portabilité du projet.
```