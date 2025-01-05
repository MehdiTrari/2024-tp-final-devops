# Rapport de Déploiement avec Docker, Docker Compose et Intégration Continue

## Introduction
Dans ce projet, j'ai configuré et déployé une application multi-services composée de :
- Une API (`vote-api`),
- Une interface utilisateur (`web-client`),
- Une documentation statique (`docs`),
- Une base de données PostgreSQL (`db`).

J'ai utilisé Docker pour conteneuriser chaque service, Docker Compose pour orchestrer leur déploiement, et une pipeline d'intégration continue (CI) pour automatiser les tests et la validation du code.

---

## Méthode de Développement Agile et Gestion des Branches

Je travaille avec la méthode agile, ce qui implique une organisation stricte des branches pour faciliter le développement collaboratif :
- **Branche `main`** : Contient la version stable de l'application, prête à être déployée en production.
- **Branche `develop`** : Utilisée pour regrouper toutes les fonctionnalités en cours de développement avant leur validation finale.
- **Branches `feature/*`** : Créées à partir de `develop` pour chaque nouvelle fonctionnalité ou tâche. Une fois le travail terminé, elles sont fusionnées dans `develop`.
- **Hotfixes** : Les corrections urgentes (`hotfix`) sont créées directement depuis `main` pour résoudre des bugs critiques.

### Validation des Pull Requests
J'ai configuré GitHub pour :
- Exiger une validation automatique par les tests CI avant la fusion des PR.
- Empêcher la fusion si des tests échouent ou si la branche n'a pas été approuvée par un relecteur.

Cette configuration assure une gestion fluide et cohérente des versions tout en maintenant la stabilité de la branche principale.

---

## Intégration Continue avec GitHub Actions

### Fichier `test.yml`

```yaml
name: CI Pipeline

on:
  push:
    branches:
      - main
      - develop
      - feature/*
  pull_request:
    branches:
      - main
      - develop
      - feature/*

jobs:
  build:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_USER: user
          POSTGRES_PASSWORD: password
          POSTGRES_DB: testdb
        ports:
          - 5432:5432
        options: >-
          --health-cmd="pg_isready -U user -d testdb"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=5

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21

      - name: Wait for PostgreSQL to be ready
        run: |
          for i in {1..10}; do
            nc -z localhost 5432 && echo "Postgres is ready" && exit 0
            echo "Waiting for Postgres..."
            sleep 3
          done
          echo "Postgres failed to start" && exit 1

      - name: Set environment variables
        env:
          PG_URL: ${{ secrets.PG_URL }}
        run: echo "PG_URL=${{ secrets.PG_URL }}" >> $GITHUB_ENV

      - name: Install dependencies (vote-api)
        working-directory: vote-api
        run: go mod tidy

      - name: Run Go tests
        working-directory: vote-api
        run: go test ./... -v

  webclient-tests:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'yarn'

      - name: Install dependencies
        working-directory: web-client
        run: yarn install --frozen-lockfile

      - name: Install Playwright Browsers
        working-directory: web-client
        run: npx playwright install --with-deps

      - name: Run WebClient Unit Tests
        working-directory: web-client
        run: yarn test

  build-artifacts:
    runs-on: ubuntu-latest
    needs:
      - webclient-tests
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Install dependencies and build
        working-directory: web-client
        run: |
          yarn install --frozen-lockfile
          yarn build

      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: web-client-build
          path: web-client/.next
```

---

### Explications de la Pipeline CI

#### 1. **Déclencheurs**
La pipeline est déclenchée automatiquement sur les événements suivants :
- **Push** : Lorsqu'un commit est poussé sur les branches `main`, `develop`, ou toute branche `feature/*`.
- **Pull Request** : Lorsqu'une pull request est ouverte ou mise à jour sur les mêmes branches.

#### 2. **Job `build`**
- **Objectif** : Tester l'API `vote-api` en exécutant les tests unitaires écrits en Go.
- **Configuration** :
  - Utilisation d'un conteneur PostgreSQL simulé comme base de données.
  - Installation des dépendances Go.
  - Exécution des tests avec `go test`.

#### 3. **Job `webclient-tests`**
- **Objectif** : Tester le client web (`web-client`) en exécutant les tests unitaires et les tests E2E.
- **Configuration** :
  - Installation de Node.js et des dépendances via Yarn.
  - Installation des navigateurs nécessaires pour Playwright.
  - Exécution des tests unitaires avec `yarn test`.

#### 4. **Job `build-artifacts`**
- **Objectif** : Générer les artefacts de build pour le déploiement.
- **Configuration** :
  - Compilation du client web avec `yarn build`.
  - Téléchargement des artefacts pour utilisation ultérieure.

---

## Rollback
Pour ce projet, j'ai choisi de gérer les rollbacks avec des **tags Git**. Cette méthode me permet :
1. D'identifier rapidement une version stable.
2. De déployer une version précédente en cas de problème, simplement en basculant vers le commit associé au tag.
3. D'assurer un suivi clair des versions, ce qui facilite la gestion de la production.

---

## Résultats et Conclusion

Cette configuration garantit une validation continue du code et une gestion efficace des branches et des versions. La robustesse de la pipeline CI et la gestion agile des branches offrent une base solide pour le développement collaboratif et le déploiement fiable.