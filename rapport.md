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

<<<<<<< Updated upstream
Voici le fichier `docker-compose.yml` utilisé pour orchestrer les services :
=======
Je travaille avec la méthode agile, ce qui implique une organisation stricte des branches pour faciliter le développement collaboratif :
- **Branche `main`** : Contient la version stable de l'application, prête à être déployée en production.
- **Branche `develop`** : Utilisée pour regrouper toutes les fonctionnalités en cours de développement avant leur validation finale.
- **Branches `feature/*`** : Créées à partir de `develop` pour chaque nouvelle fonctionnalité ou tâche. Une fois le travail terminé, elles sont fusionnées dans `develop`.
- **Hotfixes** : Les corrections urgentes (`hotfix`) sont créées directement depuis `main` pour résoudre des bugs critiques en production.

Ce workflow permet une gestion fluide et cohérente des versions tout en maintenant la stabilité de la branche principale.

---

## Rollbacks

Pour mon projet, j'ai choisi une approche basée sur les **tags Docker** pour effectuer des rollbacks. Cette méthode présente plusieurs avantages dans le cadre d'un projet conteneurisé et déployé via Docker Compose :
- Les images Docker taguées permettent de revenir rapidement à une version stable et validée.
- Les tags facilitent la gestion des versions en production, surtout en cas de problème critique.

### Méthode Utilisée : Rollback avec un Tag Docker

1. **Création d'un Tag lors d'une Version Stable :**
   À chaque version validée, je crée un tag Docker correspondant, par exemple `v1.0.0`. Cela garantit que je peux toujours revenir à cette version en cas de problème avec les déploiements futurs.

   ```bash
   docker tag vote-api:latest vote-api:v1.0.0
   docker push vote-api:v1.0.0
   ```

2. **Rollback en Cas de Problème :**
   Si une nouvelle version introduit un bug critique, je peux effectuer un rollback en utilisant directement l'image Docker taguée.

   Exemple avec Docker Compose :
   ```bash
   docker-compose down
   docker-compose pull vote-api:v1.0.0
   docker-compose up -d vote-api
   ```

   Cette méthode remet instantanément en ligne la version stable précédemment taguée sans avoir à reconstruire l'image ou à modifier le code source.

3. **Automatisation dans le Pipeline CI/CD :**
   Dans le cadre d'une intégration continue, j'aurais pu ajouter un job dédié au déploiement d'une version taguée en production. Voici un exemple de configuration possible :

   ```yaml
   jobs:
     rollback:
       runs-on: ubuntu-latest
       steps:
         - name: Deploy Tagged Version
           run: |
             docker pull vote-api:v1.0.0
             docker run -d -p 8080:8080 vote-api:v1.0.0
   ```

---

### Pourquoi ce Choix ?

J'ai opté pour cette méthode pour les raisons suivantes :
1. **Fiabilité :** Les images Docker taguées encapsulent tout ce qui est nécessaire pour exécuter une version spécifique de l'application, réduisant les risques liés aux dépendances ou aux configurations.
2. **Simplicité :** Cette méthode est simple à mettre en œuvre, surtout dans un environnement où Docker est déjà utilisé pour le déploiement.
3. **Rapidité :** Revenir à une version stable avec Docker Compose ne nécessite que quelques commandes et est presque instantané.

---

## Intégration Continue avec GitHub Actions

### Fichier `test.yml`
>>>>>>> Stashed changes

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

<<<<<<< Updated upstream
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
=======
## Résultats

### Tests Unitaires
- Les tests unitaires valident que chaque fonctionnalité individuelle fonctionne correctement.

### Tests E2E
- Les tests E2E valident les interactions entre les différents composants pour garantir une expérience utilisateur fluide.
>>>>>>> Stashed changes

---

## Conclusion
<<<<<<< Updated upstream
Ce projet est maintenant entièrement conteneurisé et peut être exécuté facilement avec Docker Compose. Chaque service est isolé, mais ils interagissent grâce à la configuration centralisée. Les fichiers Docker et Docker Compose simplifient la gestion des dépendances et assurent la portabilité du projet.
```
=======

En adoptant une approche centrée sur Docker pour les rollbacks, j'ai priorisé la simplicité, la rapidité et la fiabilité. Associée à une pipeline CI robuste et à une gestion rigoureuse des branches, cette solution garantit un workflow agile, une stabilité en production, et une facilité de maintenance en cas d'incident.
>>>>>>> Stashed changes
