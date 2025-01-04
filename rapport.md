# Rapport de Déploiement avec Docker et Docker Compose

## Introduction
Dans ce projet, j'ai configuré et déployé une application composée de plusieurs services (`vote-api`, `web-client`, une base de données PostgreSQL et la documentation basée sur Docusaurus) en utilisant Docker et Docker Compose. J'ai créé des Dockerfiles spécifiques pour chaque service et utilisé un fichier `docker-compose.yml` pour orchestrer leur interaction.

---

## Configuration avec Docker Compose

Le fichier `docker-compose.yml` a été indispensable pour faire fonctionner l'application, car il m'a permis de gérer les dépendances entre les services et de les exécuter simultanément de manière cohérente.

### Structure du `docker-compose.yml`

```yaml
version: '3.3'

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
      - NEXT_PUBLIC_API_URL=http://localhost:8080
    depends_on:
      - api

  docs:
    build: ./docs
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
   - J'ai configuré l'API `vote-api` pour qu'elle soit construite avec le Dockerfile situé dans le répertoire `vote-api`.
   - Le service est exposé sur le port `8080` et utilise une variable d'environnement `PG_URL` pour se connecter à la base de données.

2. **Service `web-client`** :
   - J'ai configuré le client web pour qu'il soit construit avec le Dockerfile du répertoire `web-client`.
   - Le service est exposé sur le port `3000` et la variable `NEXT_PUBLIC_API_URL` est définie pour pointer vers l'API.

3. **Service `docs`** :
   - J'ai intégré la documentation Docusaurus en utilisant le Dockerfile du répertoire `docs`.
   - Le service est exposé sur le port `4000`.

4. **Service `db`** :
   - J'ai utilisé l'image officielle `postgres:15-alpine` pour la base de données.
   - Les variables d'environnement (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`) permettent de configurer l'utilisateur, le mot de passe et le nom de la base de données.
   - J'ai ajouté un **healthcheck** pour m'assurer que PostgreSQL est prêt avant que les autres services ne s'y connectent.

5. **Volumes** :
   - Le volume `postgres_data` est utilisé pour persister les données de la base de données.

---

## Dockerfile pour `vote-api`

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

### Explications

1. **Base de l'image** : 
   - J'ai utilisé l'image `golang:1.23-rc-alpine` pour réduire la taille finale de l'image.
2. **Installation des dépendances** :
   - Les fichiers `go.mod` et `go.sum` sont copiés, puis les dépendances sont téléchargées avec `go mod download`.
3. **Construction de l'application** :
   - J'ai compilé l'application en un exécutable nommé `main`.

---

## Dockerfile pour `web-client`

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

### Explications

1. **Construction en deux étapes** :
   - La première étape compile le projet avec `yarn build`.
   - La deuxième étape utilise une image plus légère (`node:18-slim`) pour exécuter le projet.

---

## Dockerfile pour `docs`

```dockerfile
# Étape 1 : Build de la documentation
FROM node:18 as builder
WORKDIR /app

# Copier les fichiers de configuration
COPY package.json yarn.lock ./

RUN yarn install

# Copier le reste du code
COPY . ./

RUN yarn build

# Étape 2 : Image pour servir la documentation
FROM nginx:alpine
WORKDIR /usr/share/nginx/html

COPY --from=builder /app/build .

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### Explications

1. **Étape 1 : Construction** :
   - J'ai utilisé `node:18` pour installer les dépendances et construire la documentation avec `yarn build`.
2. **Étape 2 : Servir la documentation** :
   - J'ai utilisé l'image `nginx:alpine` pour servir les fichiers statiques générés.

---

## Pourquoi Docker Compose était indispensable
Docker Compose était indispensable pour :
- Gérer les dépendances entre les services.
- Synchroniser le démarrage des conteneurs avec `depends_on` et les **healthchecks**.
- Simplifier l'exécution locale des services avec une seule commande : `docker-compose up`.

---

## Conclusion
J'ai configuré une application multi-services en utilisant Docker et Docker Compose. Le projet est maintenant prêt à être utilisé dans différents environnements, notamment en production une fois déployé.
```